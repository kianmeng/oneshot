package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	discoveryserver "github.com/forestnode-io/oneshot/v2/pkg/commands/discovery-server"
	"github.com/rs/zerolog"
)

func NegotiateOfferRequest(ctx context.Context, url, username, password string, client *http.Client) (*discoveryserver.ClientOfferRequestResponse, error) {
	log := zerolog.Ctx(ctx)

	// perform the first request which saves our spot in the queue.
	// we're going to use the same pathways as browser clients to we
	// set the accept header to text/http and the user agent to oneshot.
	// the server will respond differently based on the user agent, it wont send html.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Accept", "text/html")
	req.Header.Set("User-Agent", "oneshot")
	if username != "" || password != "" {
		req.SetBasicAuth(username, password)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get token request response: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	log.Debug().
		Int("status", resp.StatusCode).
		Interface("headers", resp.Header).
		Msg("got token request response")

	cookies := resp.Cookies()
	sessionToken := ""
	for _, cookie := range cookies {
		if cookie.Name == "session_token" {
			sessionToken = cookie.Value
			break
		}
	}
	if sessionToken == "" {
		return nil, fmt.Errorf("failed to get session token")
	}

	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create offer request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Session-Token", sessionToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request offer response: %w", err)
	}
	defer resp.Body.Close()

	var corr discoveryserver.ClientOfferRequestResponse
	if err := json.NewDecoder(resp.Body).Decode(&corr); err != nil {
		return nil, fmt.Errorf("failed to decode offer request response: %w", err)
	}

	log.Debug().
		Int("status", resp.StatusCode).
		Interface("headers", resp.Header).
		Interface("response_object", corr).
		Msg("got offer request response")

	return &corr, nil
}
