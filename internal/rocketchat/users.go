package rocketchat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// GetMe returns the currently authenticated user.
func (c *Client) GetMe(ctx context.Context) (*User, error) {
	raw, err := c.get(ctx, "/me", nil)
	if err != nil {
		return nil, fmt.Errorf("get current user: %w", err)
	}
	var user User
	if err := json.Unmarshal(raw, &user); err != nil {
		return nil, fmt.Errorf("parse current user response: %w", err)
	}
	return &user, nil
}

// ListUsers returns a paginated list of users, optionally filtered by query.
func (c *Client) ListUsers(ctx context.Context, opts ListOptions) (*UserListResponse, error) {
	params := pagingParams(opts, 100)
	if opts.Query != "" {
		params.Set("query", opts.Query)
	}
	raw, err := c.get(ctx, "/users.list", params)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	var resp UserListResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse users list response: %w", err)
	}
	return &resp, nil
}

// GetUserInfo returns detailed information about a user by username.
func (c *Client) GetUserInfo(ctx context.Context, username string) (*User, error) {
	params := url.Values{"username": {username}}
	raw, err := c.get(ctx, "/users.info", params)
	if err != nil {
		return nil, fmt.Errorf("get user info for %q: %w", username, err)
	}
	var resp struct {
		User User `json:"user"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse user info response: %w", err)
	}
	return &resp.User, nil
}
