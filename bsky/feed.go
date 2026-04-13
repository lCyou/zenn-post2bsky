package bsky

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Post struct {
	Text      string
	CreatedAt string
}

func (c *Client) GetOwnRecentPosts(limit int) ([]Post, error) {
	url := fmt.Sprintf(
		"https://api.bsky.app/xrpc/app.bsky.feed.getAuthorFeed?actor=%s&limit=%d&filter=posts_no_replies",
		c.DID, limit,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AccessJWT)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("feed error: status %d", resp.StatusCode)
	}

	var result struct {
		Feed []struct {
			Post struct {
				Record struct {
					Text      string `json:"text"`
					CreatedAt string `json:"createdAt"`
				} `json:"record"`
			} `json:"post"`
		} `json:"feed"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	posts := make([]Post, 0, len(result.Feed))
	for _, item := range result.Feed {
		posts = append(posts, Post{
			Text:      item.Post.Record.Text,
			CreatedAt: item.Post.Record.CreatedAt,
		})
	}
	return posts, nil
}
