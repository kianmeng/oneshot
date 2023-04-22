package messages

import (
	"time"

	"github.com/pion/webrtc/v3"
)

type Message interface {
	Type() string
}

type VersionInfo struct {
	APIVersion string
	Version    string
}

type Handshake struct {
	ID          string
	VersionInfo VersionInfo
	Error       string
}

func (h *Handshake) Type() string {
	return "Handshake"
}

type BasicAuth struct {
	UsernameHash []byte
	PasswordHash []byte
}

type SessionURLRequest struct {
	URL      string
	Required bool
}

// sent from the oneshot server to the signalling server after VersionInfo has been exchanged
type ServerArrivalRequest struct {
	BasicAuth          *BasicAuth
	URL                *SessionURLRequest
	Redirect           string
	RedirectOnly       bool
	IsUsingPortMapping bool
	TTL                time.Duration
	Hostname           string
	Cmd                string
}

func (a *ServerArrivalRequest) Type() string {
	return "ServerArrivalRequest"
}

// sent from the signalling server to the oneshot server when it first connects in response to an ArrivalRequest
type ServerArrivalResponse struct {
	AssignedURL string
	Error       string
}

func (a *ServerArrivalResponse) Type() string {
	return "ServerArrivalResponse"
}

type AnswerOfferRequest struct {
	SessionID string
	Answer    string
}

func (a *AnswerOfferRequest) Type() string {
	return "AnswerOfferRequest"
}

type AnswerOfferResponse struct {
	SessionID string
	Error     string
}

func (a *AnswerOfferResponse) Type() string {
	return "AnswerOfferResponse"
}

// sent from the signalling server to the oneshot server when a new session has been request by a client
type GetOfferRequest struct {
	SessionID     string
	Configuration *webrtc.Configuration `json:",omitempty"`
}

func (g *GetOfferRequest) Type() string {
	return "GetOfferRequest"
}

// sent from the oneshot server to the signalling server when it has crafted an offer for the client requesting a session
type GetOfferResponse struct {
	SessionID string
	Offer     string
}

func (g *GetOfferResponse) Type() string {
	return "GetOfferResponse"
}

// sent from the signalling server to the oneshot server when a client has answered the offer
type GotAnswerRequest struct {
	SessionID string
	Answer    string
}

func (g *GotAnswerRequest) Type() string {
	return "GotAnswerRequest"
}

// sent from the oneshot server to the signalling server when it has accepted the answer and started the session
type GotAnswerResponse struct {
	SessionID string
	Error     string
}

func (g *GotAnswerResponse) Type() string {
	return "GotAnswerResponse"
}

// sent from the oneshot server to the signalling server when a session has ended
type FinishedSessionRequest struct {
	SessionID string
	Error     string
}

func (f *FinishedSessionRequest) Type() string {
	return "FinishedSessionRequest"
}

// sent from the signalling server to the oneshot server when it has received the FinishedSessionRequest
type FinishedSessionResponse struct {
	SessionID string
	Error     string
}

func (f *FinishedSessionResponse) Type() string {
	return "FinishedSessionResponse"
}

type Ping struct{}

func (p *Ping) Type() string {
	return "Ping"
}

type UpdatePingRateRequest struct {
	Period time.Duration
}

func (u *UpdatePingRateRequest) Type() string {
	return "UpdatePingRateRequest"
}
