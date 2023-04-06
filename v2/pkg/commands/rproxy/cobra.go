package rproxy

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/raphaelreyna/oneshot/v2/pkg/commands"
	"github.com/raphaelreyna/oneshot/v2/pkg/configuration"
	"github.com/raphaelreyna/oneshot/v2/pkg/events"
	"github.com/raphaelreyna/oneshot/v2/pkg/output"
	"github.com/spf13/cobra"
)

func New(config *configuration.Root) *Cmd {
	return &Cmd{
		config: config,
	}
}

type Cmd struct {
	cobraCommand *cobra.Command
	host         string
	config       *configuration.Root
}

func (c *Cmd) Cobra() *cobra.Command {
	if c.cobraCommand != nil {
		return c.cobraCommand
	}

	c.cobraCommand = &cobra.Command{
		Use:     "reverse-proxy host",
		Aliases: []string{"rproxy"},
		Short:   "Reverse proxy all requests to the specified host",
		Long:    `Reverse proxy all requests to the specified host. The host may be a URL or a host:port combination.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			config := c.config.Subcommands.RProxy
			config.MergeFlags()
			if err := config.Validate(); err != nil {
				return errors.New("invalid configuration")
			}
			if err := config.Hydrate(); err != nil {
				return errors.New("failed to hydrate configuration")
			}
			return nil
		},
		RunE: c.setHandlerFunc,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("proxy host required")
			}
			if 1 < len(args) {
				return errors.New("too many arguments, only 1 host may be used")
			}
			return nil
		},
	}

	c.cobraCommand.SetUsageTemplate(usageTemplate)

	c.config.Subcommands.RProxy.SetFlags(c.cobraCommand, c.cobraCommand.LocalFlags())

	return c.cobraCommand
}

func (c *Cmd) setHandlerFunc(cmd *cobra.Command, args []string) error {
	var (
		ctx = cmd.Context()

		host      = args[0]
		config    = c.config.Subcommands.RProxy
		spoofHost = config.SpoofHost
	)

	output.IncludeBody(ctx)

	hostURL, err := url.Parse(host)
	if err != nil {
		return err
	}

	var spoofedHostURL *url.URL
	if config.MatchHost {
		spoofedHostURL = hostURL
	}
	if spoofHost != "" {
		spoofHost = strings.TrimPrefix(spoofHost, "http://")
		spoofHost = strings.TrimPrefix(spoofHost, "https://")
		if idx := strings.Index(spoofHost, "/"); -1 < idx {
			spoofHost = spoofHost[:idx]
		}
		spoofedHostURL = &url.URL{
			Host: spoofHost,
		}
	}

	c.host = host

	if spoofedHostURL != nil {
		if len(config.RequestHeader) == 0 {
			config.RequestHeader = make(map[string][]string)
		}
		config.RequestHeader["Host"] = []string{spoofedHostURL.Host}
	}

	rp := httputil.NewSingleHostReverseProxy(hostURL)
	rp.ModifyResponse = func(resp *http.Response) error {
		ctx := c.cobraCommand.Context()
		originalHeader := resp.Header.Clone()

		events.Raise(ctx, &events.HTTPResponse{
			StatusCode: resp.StatusCode,
			Header:     originalHeader,
		})

		if 0 < len(config.ResponseHeader) {
			for k, v := range config.ResponseHeader {
				resp.Header[k] = v
			}
		}
		if config.StatusCode != 0 {
			resp.StatusCode = config.StatusCode
		}
		return nil
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx := c.cobraCommand.Context()
		events.Raise(ctx, output.NewHTTPRequest(r))

		if 0 < len(config.RequestHeader) {
			for k, v := range config.RequestHeader {
				r.Header[k] = v
			}
		}

		if spoofedHostURL != nil {
			r.Host = spoofedHostURL.Host
			r.URL = spoofedHostURL
		}

		if config.Method != "" {
			r.Method = strings.ToUpper(config.Method)
		}

		var jsonOutput bool
		if format, _ := output.GetFormatAndOpts(ctx); format == "json" {
			jsonOutput = true
		}

		if config.Tee || jsonOutput {
			bw, getBufByte := output.NewBufferedWriter(ctx, w)

			ww := bw.(http.ResponseWriter)
			rp.ServeHTTP(ww, r)

			events.Raise(ctx, &events.File{
				Content: getBufByte,
			})
		} else {
			rp.ServeHTTP(w, r)
		}

		events.Success(ctx)
	}

	commands.SetHTTPHandlerFunc(ctx, handler)
	return nil
}
