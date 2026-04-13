package ui

import "github.com/kyou/post2bsky/bsky"

// PostDoneMsg is sent when CreatePost succeeds.
type PostDoneMsg struct{}

// PostErrMsg is sent when CreatePost fails.
type PostErrMsg struct{ Err error }

// FeedLoadedMsg is sent when GetOwnRecentPosts returns.
type FeedLoadedMsg struct{ Posts []bsky.Post }

// FeedErrMsg is sent when GetOwnRecentPosts fails.
type FeedErrMsg struct{ Err error }

// AuthDoneMsg is sent when authentication completes.
type AuthDoneMsg struct{}

// AuthErrMsg is sent when authentication fails.
type AuthErrMsg struct{ Err error }
