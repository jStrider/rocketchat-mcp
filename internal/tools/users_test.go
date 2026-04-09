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

func TestListUsersHandler(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rocketchat.UserListResponse{
			Users: []rocketchat.User{
				{ID: "u1", Username: "alice", Name: "Alice Smith", Status: "online"},
				{ID: "u2", Username: "bob", Name: "Bob Jones", Status: "away"},
			},
			Pagination: rocketchat.Pagination{Count: 2, Total: 2},
		})
	})

	handler := makeListUsersHandler(client)
	result, err := handler(context.Background(), callToolRequest("list_users", map[string]any{}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "@alice")
	assert.Contains(t, text, "@bob")
	assert.Contains(t, text, "online")
}

func TestGetUserInfoHandler(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"user": rocketchat.User{
				ID: "u1", Username: "alice", Name: "Alice Smith",
				Status: "online", Active: true, Roles: []string{"admin", "user"},
			},
		})
	})

	handler := makeGetUserInfoHandler(client)
	result, err := handler(context.Background(), callToolRequest("get_user_info", map[string]any{"username": "alice"}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "@alice")
	assert.Contains(t, text, "Alice Smith")
	assert.Contains(t, text, "admin")
}

func TestGetUserInfoHandler_MissingUsername(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})

	handler := makeGetUserInfoHandler(client)
	result, err := handler(context.Background(), callToolRequest("get_user_info", map[string]any{}))
	require.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestGetMeHandler(t *testing.T) {
	client := newToolTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/me", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rocketchat.User{
			ID: "uid1", Username: "admin", Name: "Admin User", Status: "online",
		})
	})

	handler := makeGetMeHandler(client)
	result, err := handler(context.Background(), callToolRequest("get_me", map[string]any{}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "@admin")
	assert.Contains(t, text, "Admin User")
}
