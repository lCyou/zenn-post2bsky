package bsky

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type SessionCache struct {
	AccessJWT  string `json:"access_jwt"`
	RefreshJWT string `json:"refresh_jwt"`
	DID        string `json:"did"`
	Handle     string `json:"handle"`
}

func (c *Client) CreateSession(handle, appPassword string) error {
	body := map[string]string{
		"identifier": handle,
		"password":   appPassword,
	}
	var result struct {
		AccessJwt  string `json:"accessJwt"`
		RefreshJwt string `json:"refreshJwt"`
		Did        string `json:"did"`
		Handle     string `json:"handle"`
	}
	if err := c.doXRPC("POST", "com.atproto.server.createSession", body, &result); err != nil {
		return err
	}
	c.AccessJWT = result.AccessJwt
	c.DID = result.Did
	c.Handle = result.Handle
	return c.SaveSession(SessionCache{
		AccessJWT:  result.AccessJwt,
		RefreshJWT: result.RefreshJwt,
		DID:        result.Did,
		Handle:     result.Handle,
	})
}

func (c *Client) RefreshSession(refreshJWT string) error {
	// Temporarily use the refresh token as the auth header
	orig := c.AccessJWT
	c.AccessJWT = refreshJWT
	defer func() {
		if c.AccessJWT == refreshJWT {
			c.AccessJWT = orig
		}
	}()

	var result struct {
		AccessJwt  string `json:"accessJwt"`
		RefreshJwt string `json:"refreshJwt"`
		Did        string `json:"did"`
		Handle     string `json:"handle"`
	}
	if err := c.doXRPC("POST", "com.atproto.server.refreshSession", nil, &result); err != nil {
		return err
	}
	c.AccessJWT = result.AccessJwt
	c.DID = result.Did
	c.Handle = result.Handle
	return c.SaveSession(SessionCache{
		AccessJWT:  result.AccessJwt,
		RefreshJWT: result.RefreshJwt,
		DID:        result.Did,
		Handle:     result.Handle,
	})
}

func (c *Client) SaveSession(s SessionCache) error {
	path, err := sessionCachePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0600)
}

func LoadSession() (*SessionCache, error) {
	path, err := sessionCachePath()
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var s SessionCache
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func sessionCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "post2bsky", "session.json"), nil
}
