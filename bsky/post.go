package bsky

import "time"

func (c *Client) CreatePost(text string) error {
	body := map[string]any{
		"repo":       c.DID,
		"collection": "app.bsky.feed.post",
		"record": map[string]any{
			"$type":     "app.bsky.feed.post",
			"text":      text,
			"createdAt": time.Now().UTC().Format(time.RFC3339),
		},
	}
	return c.doXRPC("POST", "com.atproto.repo.createRecord", body, nil)
}
