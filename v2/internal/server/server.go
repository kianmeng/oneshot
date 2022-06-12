package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/raphaelreyna/oneshot/v2/internal/summary"
)

var ErrTimeout = errors.New("timeout")

type connKey struct{}

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) (*summary.Request, error)
	ServeExpiredHTTP(http.ResponseWriter, *http.Request)
}

type _wr struct {
	w         http.ResponseWriter
	r         *http.Request
	done      func()
	startTime time.Time
}

type Server struct {
	http.Server

	requestsQueue chan _wr
	requestsWG    sync.WaitGroup

	summary     *summary.Summary
	succChan    chan struct{}
	timeoutChan chan struct{}
	stopWorkers func()

	serveHTTP        HandlerFunc
	serveExpiredHTTP http.HandlerFunc

	bufferRequestBody bool
	timeout           time.Duration
}

func NewServer(handler Handler) *Server {
	s := Server{
		requestsQueue:    make(chan _wr),
		serveHTTP:        handler.ServeHTTP,
		serveExpiredHTTP: handler.ServeExpiredHTTP,
		succChan:         make(chan struct{}),
		summary:          summary.NewSummary(time.Now()),
	}

	return &s
}

func (s *Server) AddMiddleware(m Middleware) {
	if m == nil {
		return
	}

	s.serveHTTP = m(s.serveHTTP)
	var hf = m(newHandlerFunc(s.serveExpiredHTTP, false))
	s.serveExpiredHTTP = func(w http.ResponseWriter, r *http.Request) { hf(w, r) }
}

func (s *Server) Summary() *summary.Summary {
	return s.summary
}

func (s *Server) Serve(ctx context.Context, l net.Listener) error {
	r := mux.NewRouter()
	r.HandleFunc("/", s.rootHandler)

	s.Handler = r
	s.BaseContext = func(l net.Listener) context.Context {
		return ctx
	}
	s.ConnContext = func(ctx context.Context, c net.Conn) context.Context {
		return context.WithValue(ctx, connKey{}, c)
	}

	ctx, s.stopWorkers = context.WithCancel(ctx)

	go s.preTransferWorker(ctx)

	if 0 < s.timeout {
		s.timeoutChan = make(chan struct{}, 1)
		l = withTimeout(s.timeout, s.timeoutChan, l)
	}
	var errs = make(chan error)
	go func() {
		errs <- s.Server.Serve(l)
	}()

	go func() {
		<-s.succChan
		s.requestsWG.Wait()
		errs <- nil
	}()

	go func() {
		<-s.timeoutChan
		fmt.Printf("timeout after %v; exiting ...\n", s.timeout)
		errs <- nil
	}()

	err := <-errs
	if err != nil {
		return err
	}

	s.summary.End()
	s.Shutdown(ctx)
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	close(s.requestsQueue)
	s.stopWorkers()

	return s.Server.Close()
}

func (s *Server) BufferRequests() {
	s.bufferRequestBody = true
}

func (s *Server) SetTimeoutDuration(d time.Duration) {
	s.timeout = d
}

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	// go dark if the transfer succeeded
	if s.summary.Succesful() {
		var (
			ctx  = r.Context()
			conn = ctx.Value(connKey{}).(net.Conn)
		)

		conn.Close()
		return
	}

	var (
		done = make(chan struct{})

		wr = _wr{
			w:         w,
			r:         r,
			startTime: time.Now(),
			done: func() {
				done <- struct{}{}
				close(done)
			},
		}
	)

	s.requestsWG.Add(1)
	s.requestsQueue <- wr

	<-done
	s.requestsWG.Done()
}

func (s *Server) preTransferWorker(ctx context.Context) {
	for {
		if s.summary.Succesful() {
			s.succChan <- struct{}{}
			for i := 0; i < runtime.NumCPU(); i++ {
				go s.postTransferWorker(ctx)
			}
			return
		}

		select {
		case <-ctx.Done():
			return
		case wr, ok := <-s.requestsQueue:
			if !ok {
				return
			}

			var (
				w   = summary.NewResponseWriter(wr.w)
				req *summary.Request
				err error
			)

			if s.bufferRequestBody {
				bodyBuf := buffered(wr.r)
				req, err = s.serveHTTP(w, wr.r)
				req.SetBody(bodyBuf)
			} else {
				req, err = s.serveHTTP(w, wr.r)
			}

			if req != nil {
				req.SetTimes(wr.startTime, time.Now())
				req.SetWriteStats(w)

				if err != nil {
					s.summary.AddFailure(req)
				} else {
					s.summary.SucceededWith(req)
				}
			}

			wr.done()
		}
	}
}

func (s *Server) postTransferWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case wr, ok := <-s.requestsQueue:
			if !ok {
				return
			}

			s.serveExpiredHTTP(wr.w, wr.r)

			wr.done()
		}
	}
}

type bufferedBody struct {
	io.ReadCloser
	buf *bufio.Reader
}

func (bb *bufferedBody) Close() error {
	io.ReadAll(bb)
	return bb.ReadCloser.Close()
}

func (bb *bufferedBody) Read(p []byte) (int, error) {
	return bb.buf.Read(p)
}

func buffered(r *http.Request) *bufio.Reader {
	var (
		body = r.Body
		bb   = bufferedBody{
			ReadCloser: body,
			buf:        bufio.NewReader(body),
		}
	)

	(*r).Body = &bb

	return bb.buf
}
