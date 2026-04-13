package bsky

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


type Client struct {
	PDS         string
	HTTPClient  *http.Client
	AccessJWT   string
	DID         string
	Handle      string
	appPassword string
}

func (c *Client) SetCredentials(handle, appPassword string) {
	c.Handle = handle
	c.appPassword = appPassword
}

// Authenticate tries the cached session first, falls back to createSession.
func (c *Client) Authenticate() error {
	cache, err := LoadSession()
	if err == nil && cache != nil {
		if rerr := c.RefreshSession(cache.RefreshJWT); rerr == nil {
			return nil
		}
	}
	if c.Handle == "" || c.appPassword == "" {
		return fmt.Errorf("no credentials available")
	}
	return c.CreateSession(c.Handle, c.appPassword)
}

func NewClient(pds string) *Client {
	return &Client{
		PDS:        pds,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) doXRPC(method, endpoint string, body any, result any) error {
	url := c.PDS + "/xrpc/" + endpoint

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.AccessJWT != "" {
		req.Header.Set("Authorization", "Bearer "+c.AccessJWT)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("XRPC error %d: %s", resp.StatusCode, string(respBytes))
	}

	if result != nil {
		if err := json.Unmarshal(respBytes, result); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}
	}

	return nil
}
