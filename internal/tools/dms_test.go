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

func TestListDMsHandler(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/im.list", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rocketchat.DMListResponse{
			IMs: []rocketchat.DMRoom{
				{
					ID:        "dm1",
					Usernames: []string{"jrenaud", "alice"},
					LastMessage: &rocketchat.Message{Text: "salut!"},
				},
				{
					ID:        "dm2",
					Usernames: []string{"jrenaud", "bob"},
				},
			},
			Pagination: rocketchat.Pagination{Count: 2, Total: 2},
		})
	})

	handler := makeListDMsHandler(client)
	result, err := handler(context.Background(), callToolRequest("list_dms", map[string]any{}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "alice")
	assert.Contains(t, text, "salut!")
	assert.Contains(t, text, "bob")
}

func TestGetDMMessagesHandler(t *testing.T) {
	reqCount := 0
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		reqCount++
		w.Header().Set("Content-Type", "application/json")
		if reqCount == 1 {
			// im.create
			json.NewEncoder(w).Encode(map[string]any{
				"room": rocketchat.Room{ID: "dm-room-1"},
			})
			return
		}
		// im.history
		json.NewEncoder(w).Encode(rocketchat.MessageListResponse{
			Messages: []rocketchat.Message{
				{ID: "m1", Text: "hello!", User: rocketchat.MessageUser{Username: "alice"}, Timestamp: "2026-04-09T10:00:00Z"},
				{ID: "m2", Text: "how are you?", User: rocketchat.MessageUser{Username: "jrenaud"}, Timestamp: "2026-04-09T10:01:00Z"},
			},
		})
	})

	handler := makeGetDMMessagesHandler(client)
	result, err := handler(context.Background(), callToolRequest("get_dm_messages", map[string]any{"username": "alice"}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "DM with @alice")
	assert.Contains(t, text, "hello!")
	assert.Contains(t, text, "@jrenaud")
}

func TestGetDMMessagesHandler_MissingUsername(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})

	handler := makeGetDMMessagesHandler(client)
	result, err := handler(context.Background(), callToolRequest("get_dm_messages", map[string]any{}))
	require.NoError(t, err)
	assert.True(t, result.IsError)
}
