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

func TestCreateChannelHandler(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/channels.create", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"channel": rocketchat.Channel{ID: "ch-new", Name: "test-channel"},
		})
	})

	handler := makeCreateChannelHandler(client)
	result, err := handler(context.Background(), callToolRequest("create_channel", map[string]any{
		"name": "test-channel",
	}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "test-channel")
	assert.Contains(t, text, "created")
}

func TestCreateChannelHandler_WithMembers(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		members, ok := body["members"].([]any)
		assert.True(t, ok)
		assert.Len(t, members, 2)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"channel": rocketchat.Channel{ID: "ch1", Name: "team"},
		})
	})

	handler := makeCreateChannelHandler(client)
	result, err := handler(context.Background(), callToolRequest("create_channel", map[string]any{
		"name":    "team",
		"members": "alice, bob",
	}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
}

func TestCreateChannelHandler_MissingName(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})

	handler := makeCreateChannelHandler(client)
	result, err := handler(context.Background(), callToolRequest("create_channel", map[string]any{}))
	require.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestArchiveChannelHandler(t *testing.T) {
	reqCount := 0
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		reqCount++
		w.Header().Set("Content-Type", "application/json")
		if reqCount == 1 {
			// ResolveRoomID
			json.NewEncoder(w).Encode(map[string]any{
				"channel": rocketchat.Channel{ID: "ch1", Name: "old-channel"},
			})
			return
		}
		// archive
		_, _ = w.Write([]byte(`{"success":true}`))
	})

	handler := makeArchiveChannelHandler(client)
	result, err := handler(context.Background(), callToolRequest("archive_channel", map[string]any{
		"channel": "old-channel",
	}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "archived")
}

func TestInviteToChannelHandler(t *testing.T) {
	reqCount := 0
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		reqCount++
		w.Header().Set("Content-Type", "application/json")
		switch reqCount {
		case 1:
			// ResolveRoomID
			json.NewEncoder(w).Encode(map[string]any{
				"channel": rocketchat.Channel{ID: "ch1", Name: "general"},
			})
		case 2:
			// GetUserInfo
			json.NewEncoder(w).Encode(map[string]any{
				"user": rocketchat.User{ID: "u1", Username: "alice"},
			})
		case 3:
			// InviteToChannel
			_, _ = w.Write([]byte(`{"success":true}`))
		}
	})

	handler := makeInviteToChannelHandler(client)
	result, err := handler(context.Background(), callToolRequest("invite_to_channel", map[string]any{
		"channel":  "general",
		"username": "alice",
	}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "@alice")
	assert.Contains(t, text, "#general")
}

func TestListGroupsHandler(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/groups.list", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"groups": []rocketchat.Channel{
				{ID: "g1", Name: "secret-ops", MembersCount: 3, MsgCount: 100},
			},
			"count":  1,
			"offset": 0,
			"total":  1,
		})
	})

	handler := makeListGroupsHandler(client)
	result, err := handler(context.Background(), callToolRequest("list_groups", map[string]any{}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "secret-ops")
}
