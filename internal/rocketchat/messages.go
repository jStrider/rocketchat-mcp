package rocketchat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// SearchMessages searches for messages containing the given text in a channel.
func (c *Client) SearchMessages(ctx context.Context, roomID, searchText string, opts ListOptions) (*MessageListResponse, error) {
	params := pagingParams(opts, 50)
	params.Set("roomId", roomID)
	params.Set("searchText", searchText)
	raw, err := c.get(ctx, "/chat.search", params)
	if err != nil {
		return nil, fmt.Errorf("search messages in %q: %w", roomID, err)
	}
	var resp MessageListResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse search response: %w", err)
	}
	return &resp, nil
}

// PostMessage sends a message to a channel (by name or ID).
func (c *Client) PostMessage(ctx context.Context, channel, text string) (*Message, error) {
	raw, err := c.post(ctx, "/chat.postMessage", map[string]string{
		"channel": channel,
		"text":    text,
	})
	if err != nil {
		return nil, fmt.Errorf("post message to %q: %w", channel, err)
	}
	var resp struct {
		Message Message `json:"message"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse post message response: %w", err)
	}
	return &resp.Message, nil
}

// CreateDM creates a direct message room with the given username and returns the room ID.
func (c *Client) CreateDM(ctx context.Context, username string) (string, error) {
	raw, err := c.post(ctx, "/im.create", map[string]string{
		"username": username,
	})
	if err != nil {
		return "", fmt.Errorf("create DM with %q: %w", username, err)
	}
	var resp struct {
		Room Room `json:"room"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return "", fmt.Errorf("parse DM create response: %w", err)
	}
	if resp.Room.ID == "" {
		return "", fmt.Errorf("failed to create DM room with %q: empty room ID", username)
	}
	return resp.Room.ID, nil
}

// SendDM creates a DM room with the user and sends a message.
func (c *Client) SendDM(ctx context.Context, username, text string) (*Message, error) {
	roomID, err := c.CreateDM(ctx, username)
	if err != nil {
		return nil, err
	}
	raw, err := c.post(ctx, "/chat.sendMessage", map[string]any{
		"message": map[string]string{
			"rid": roomID,
			"msg": text,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("send DM to %q: %w", username, err)
	}
	var resp struct {
		Message Message `json:"message"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse send DM response: %w", err)
	}
	return &resp.Message, nil
}

// GetThreadMessages returns messages from a thread.
func (c *Client) GetThreadMessages(ctx context.Context, threadID string, opts ListOptions) (*MessageListResponse, error) {
	params := pagingParams(opts, 50)
	params.Set("tmid", threadID)
	raw, err := c.get(ctx, "/chat.getThreadMessages", params)
	if err != nil {
		return nil, fmt.Errorf("get thread messages for %q: %w", threadID, err)
	}
	var resp MessageListResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse thread messages response: %w", err)
	}
	return &resp, nil
}

// ReplyToThread sends a reply in a thread.
func (c *Client) ReplyToThread(ctx context.Context, roomID, threadID, text string) (*Message, error) {
	raw, err := c.post(ctx, "/chat.sendMessage", map[string]any{
		"message": map[string]string{
			"rid":  roomID,
			"tmid": threadID,
			"msg":  text,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("reply to thread %q: %w", threadID, err)
	}
	var resp struct {
		Message Message `json:"message"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse thread reply response: %w", err)
	}
	return &resp.Message, nil
}

// GetMessage returns a single message by ID.
func (c *Client) GetMessage(ctx context.Context, msgID string) (*Message, error) {
	params := url.Values{"msgId": {msgID}}
	raw, err := c.get(ctx, "/chat.getMessage", params)
	if err != nil {
		return nil, fmt.Errorf("get message %q: %w", msgID, err)
	}
	var resp struct {
		Message Message `json:"message"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse get message response: %w", err)
	}
	return &resp.Message, nil
}

// AddReaction adds an emoji reaction to a message.
func (c *Client) AddReaction(ctx context.Context, msgID, emoji string) error {
	_, err := c.post(ctx, "/chat.react", map[string]string{
		"messageId": msgID,
		"emoji":     emoji,
	})
	if err != nil {
		return fmt.Errorf("add reaction %q to message %q: %w", emoji, msgID, err)
	}
	return nil
}

// ListRoomFiles returns files uploaded to a room.
func (c *Client) ListRoomFiles(ctx context.Context, roomID string, opts ListOptions) (*FileListResponse, error) {
	params := pagingParams(opts, 50)
	params.Set("roomId", roomID)
	raw, err := c.get(ctx, "/channels.files", params)
	if err != nil {
		return nil, fmt.Errorf("list room files for %q: %w", roomID, err)
	}
	var resp FileListResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse room files response: %w", err)
	}
	return &resp, nil
}
