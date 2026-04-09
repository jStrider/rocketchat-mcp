package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sancare/rocketchat-mcp/internal/rocketchat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListChannelsHandler(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rocketchat.ChannelListResponse{
			Channels: []rocketchat.Channel{
				{ID: "ch1", Name: "general", MembersCount: 42, MsgCount: 100, Topic: "General chat"},
				{ID: "ch2", Name: "random", MembersCount: 10, MsgCount: 50},
			},
			Pagination: rocketchat.Pagination{Count: 2, Total: 2},
		})
	})

	handler := makeListChannelsHandler(client)
	result, err := handler(context.Background(), callToolRequest("list_channels", map[string]any{}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "#general")
	assert.Contains(t, text, "42 members")
	assert.Contains(t, text, "#random")
}

func TestListJoinedChannelsHandler(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/channels.list.joined", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rocketchat.ChannelListResponse{
			Channels:   []rocketchat.Channel{{ID: "ch1", Name: "dev"}},
			Pagination: rocketchat.Pagination{Count: 1, Total: 1},
		})
	})

	handler := makeListJoinedChannelsHandler(client)
	result, err := handler(context.Background(), callToolRequest("list_joined_channels", map[string]any{}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "#dev")
}

func TestGetChannelInfoHandler(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"channel": rocketchat.Channel{
				ID: "ch1", Name: "general", Type: "c",
				MembersCount: 42, MsgCount: 1000,
				Topic: "General discussion", Description: "Main channel",
			},
		})
	})

	handler := makeGetChannelInfoHandler(client)
	result, err := handler(context.Background(), callToolRequest("get_channel_info", map[string]any{"channel": "general"}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "#general")
	assert.Contains(t, text, "public channel")
	assert.Contains(t, text, "42")
}

func TestGetChannelInfoHandler_MissingChannel(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})

	handler := makeGetChannelInfoHandler(client)
	result, err := handler(context.Background(), callToolRequest("get_channel_info", map[string]any{}))
	require.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestGetChannelMessagesHandler(t *testing.T) {
	reqCount := 0
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		reqCount++
		w.Header().Set("Content-Type", "application/json")
		if reqCount == 1 {
			// ResolveRoomID
			json.NewEncoder(w).Encode(map[string]any{
				"channel": rocketchat.Channel{ID: "ch1", Name: "general"},
			})
			return
		}
		// channels.history
		json.NewEncoder(w).Encode(rocketchat.MessageListResponse{
			Messages: []rocketchat.Message{
				{ID: "m1", Text: "hello world", User: rocketchat.MessageUser{Username: "alice"}, Timestamp: "2026-04-09T10:00:00Z"},
			},
		})
	})

	handler := makeGetChannelMessagesHandler(client)
	result, err := handler(context.Background(), callToolRequest("get_channel_messages", map[string]any{"channel": "general"}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "@alice")
	assert.Contains(t, text, "hello world")
}
