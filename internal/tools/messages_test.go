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

func TestSearchMessagesHandler(t *testing.T) {
	reqCount := 0
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		reqCount++
		w.Header().Set("Content-Type", "application/json")
		if reqCount == 1 {
			json.NewEncoder(w).Encode(map[string]any{
				"channel": rocketchat.Channel{ID: "ch1", Name: "general"},
			})
			return
		}
		json.NewEncoder(w).Encode(rocketchat.MessageListResponse{
			Messages: []rocketchat.Message{
				{ID: "m1", Text: "deploy v2.0 successful", User: rocketchat.MessageUser{Username: "ops"}, Timestamp: "2026-04-09T10:00:00Z"},
			},
		})
	})

	handler := makeSearchMessagesHandler(client)
	result, err := handler(context.Background(), callToolRequest("search_messages", map[string]any{
		"channel":     "general",
		"search_text": "deploy",
	}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "deploy")
	assert.Contains(t, text, "@ops")
}

func TestSearchMessagesHandler_MissingArgs(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})

	handler := makeSearchMessagesHandler(client)
	result, err := handler(context.Background(), callToolRequest("search_messages", map[string]any{}))
	require.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestSendMessageHandler(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"message": rocketchat.Message{ID: "m-new", Text: "hello"},
		})
	})

	handler := makeSendMessageHandler(client)
	result, err := handler(context.Background(), callToolRequest("send_message", map[string]any{
		"channel": "general",
		"text":    "hello",
	}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "Message sent")
	assert.Contains(t, text, "#general")
}

func TestSendMessageHandler_MissingArgs(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})

	handler := makeSendMessageHandler(client)
	result, err := handler(context.Background(), callToolRequest("send_message", map[string]any{}))
	require.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestSendDMHandler(t *testing.T) {
	reqCount := 0
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		reqCount++
		w.Header().Set("Content-Type", "application/json")
		if reqCount == 1 {
			// im.create
			json.NewEncoder(w).Encode(map[string]any{
				"room": rocketchat.Room{ID: "dm-room"},
			})
			return
		}
		// chat.sendMessage
		json.NewEncoder(w).Encode(map[string]any{
			"message": rocketchat.Message{ID: "dm-msg", Text: "hey"},
		})
	})

	handler := makeSendDMHandler(client)
	result, err := handler(context.Background(), callToolRequest("send_dm", map[string]any{
		"username": "alice",
		"text":     "hey",
	}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "DM sent")
	assert.Contains(t, text, "@alice")
}

func TestSendDMHandler_MissingArgs(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})

	handler := makeSendDMHandler(client)
	result, err := handler(context.Background(), callToolRequest("send_dm", map[string]any{}))
	require.NoError(t, err)
	assert.True(t, result.IsError)
}
