package tools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sancare/rocketchat-mcp/internal/rocketchat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newToolTestClient(t *testing.T, handler http.HandlerFunc) *rocketchat.Client {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	return rocketchat.NewClient(ts.URL, "test-uid", "test-token", 10*1024*1024)
}

func callToolRequest(name string, args map[string]any) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      name,
			Arguments: args,
		},
	}
}

func TestGetThreadMessagesHandler(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "thread1", r.URL.Query().Get("tmid"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rocketchat.MessageListResponse{
			Messages: []rocketchat.Message{
				{ID: "m1", Text: "reply 1", User: rocketchat.MessageUser{Username: "alice"}, Timestamp: "2026-04-09T10:00:00Z"},
				{ID: "m2", Text: "reply 2", User: rocketchat.MessageUser{Username: "bob"}, Timestamp: "2026-04-09T10:01:00Z"},
			},
		})
	})

	handler := makeGetThreadMessagesHandler(client)
	result, err := handler(context.Background(), callToolRequest("get_thread_messages", map[string]any{"thread_id": "thread1"}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "2 messages")
	assert.Contains(t, text, "@alice")
}

func TestGetThreadMessagesHandler_MissingID(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})

	handler := makeGetThreadMessagesHandler(client)
	result, err := handler(context.Background(), callToolRequest("get_thread_messages", map[string]any{}))
	require.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestAddReactionHandler(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/chat.react", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	})

	handler := makeAddReactionHandler(client)
	result, err := handler(context.Background(), callToolRequest("add_reaction", map[string]any{
		"message_id": "msg1",
		"emoji":      ":thumbsup:",
	}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, ":thumbsup:")
}

func TestAddReactionHandler_MissingArgs(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})

	handler := makeAddReactionHandler(client)
	result, err := handler(context.Background(), callToolRequest("add_reaction", map[string]any{}))
	require.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestListRoomFilesHandler(t *testing.T) {
	reqCount := 0
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		reqCount++
		w.Header().Set("Content-Type", "application/json")
		if reqCount == 1 {
			// ResolveRoomID calls channels.info
			json.NewEncoder(w).Encode(map[string]any{
				"channel": rocketchat.Channel{ID: "ch1", Name: "general"},
			})
			return
		}
		// channels.files
		json.NewEncoder(w).Encode(rocketchat.FileListResponse{
			Files: []rocketchat.RoomFile{
				{ID: "f1", Name: "report.pdf", Type: "application/pdf", Size: 1024, User: rocketchat.MessageUser{Username: "alice"}},
			},
			Pagination: rocketchat.Pagination{Total: 1},
		})
	})

	handler := makeListRoomFilesHandler(client)
	result, err := handler(context.Background(), callToolRequest("list_room_files", map[string]any{"channel": "general"}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "report.pdf")
	assert.Contains(t, text, "@alice")
}
