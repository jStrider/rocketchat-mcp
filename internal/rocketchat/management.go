package rocketchat

import (
	"context"
	"encoding/json"
	"fmt"
)

// CreateChannel creates a new public channel.
func (c *Client) CreateChannel(ctx context.Context, name string, members []string, readOnly bool) (*Channel, error) {
	body := map[string]any{
		"name":     name,
		"readOnly": readOnly,
	}
	if len(members) > 0 {
		body["members"] = members
	}
	raw, err := c.post(ctx, "/channels.create", body)
	if err != nil {
		return nil, fmt.Errorf("create channel %q: %w", name, err)
	}
	var resp struct {
		Channel Channel `json:"channel"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse create channel response: %w", err)
	}
	return &resp.Channel, nil
}

// ArchiveChannel archives a channel by ID.
func (c *Client) ArchiveChannel(ctx context.Context, roomID string) error {
	_, err := c.post(ctx, "/channels.archive", map[string]string{
		"roomId": roomID,
	})
	if err != nil {
		return fmt.Errorf("archive channel %q: %w", roomID, err)
	}
	return nil
}

// InviteToChannel invites a user to a channel.
func (c *Client) InviteToChannel(ctx context.Context, roomID, userID string) error {
	_, err := c.post(ctx, "/channels.invite", map[string]string{
		"roomId": roomID,
		"userId": userID,
	})
	if err != nil {
		return fmt.Errorf("invite user %q to channel %q: %w", userID, roomID, err)
	}
	return nil
}

// ListGroups returns private groups the authenticated user is a member of.
func (c *Client) ListGroups(ctx context.Context, opts ListOptions) (*ChannelListResponse, error) {
	params := pagingParams(opts, 100)
	raw, err := c.get(ctx, "/groups.list", params)
	if err != nil {
		return nil, fmt.Errorf("list groups: %w", err)
	}
	var resp struct {
		Groups []Channel `json:"groups"`
		Pagination
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse groups list response: %w", err)
	}
	return &ChannelListResponse{
		Channels:   resp.Groups,
		Pagination: resp.Pagination,
	}, nil
}
