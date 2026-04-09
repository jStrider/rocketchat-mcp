package rocketchat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// ListChannels returns a paginated list of public channels.
func (c *Client) ListChannels(ctx context.Context, opts ListOptions) (*ChannelListResponse, error) {
	params := pagingParams(opts, 100)
	if opts.Query != "" {
		params.Set("query", opts.Query)
	}
	raw, err := c.get(ctx, "/channels.list", params)
	if err != nil {
		return nil, fmt.Errorf("list channels: %w", err)
	}
	var resp ChannelListResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse channels list response: %w", err)
	}
	return &resp, nil
}

// ListJoinedChannels returns channels the authenticated user has joined.
func (c *Client) ListJoinedChannels(ctx context.Context, opts ListOptions) (*ChannelListResponse, error) {
	params := pagingParams(opts, 100)
	raw, err := c.get(ctx, "/channels.list.joined", params)
	if err != nil {
		return nil, fmt.Errorf("list joined channels: %w", err)
	}
	var resp ChannelListResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse joined channels response: %w", err)
	}
	return &resp, nil
}

// GetChannelInfo returns info about a channel by name or ID.
func (c *Client) GetChannelInfo(ctx context.Context, channel string) (*Channel, error) {
	// Try by name first, then by ID.
	params := url.Values{"roomName": {channel}}
	raw, err := c.get(ctx, "/channels.info", params)
	if err != nil {
		params = url.Values{"roomId": {channel}}
		raw, err = c.get(ctx, "/channels.info", params)
		if err != nil {
			return nil, fmt.Errorf("get channel info for %q: %w", channel, err)
		}
	}
	var resp struct {
		Channel Channel `json:"channel"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse channel info response: %w", err)
	}
	return &resp.Channel, nil
}

// GetChannelHistory returns messages from a channel, newest first.
func (c *Client) GetChannelHistory(ctx context.Context, roomID string, opts ListOptions) (*MessageListResponse, error) {
	params := pagingParams(opts, 50)
	params.Set("roomId", roomID)
	raw, err := c.get(ctx, "/channels.history", params)
	if err != nil {
		return nil, fmt.Errorf("get channel history for %q: %w", roomID, err)
	}
	var resp MessageListResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse channel history response: %w", err)
	}
	return &resp, nil
}

// ResolveRoomID resolves a channel name or ID to a room ID.
// If the input looks like an existing room ID, it tries to validate it.
func (c *Client) ResolveRoomID(ctx context.Context, channel string) (string, error) {
	ch, err := c.GetChannelInfo(ctx, channel)
	if err != nil {
		return "", fmt.Errorf("resolve room %q: %w", channel, err)
	}
	return ch.ID, nil
}
