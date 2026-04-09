package rocketchat

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMe(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/me", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(User{
			ID:       "uid1",
			Username: "admin",
			Name:     "Admin User",
			Status:   "online",
			Active:   true,
		})
	})

	user, err := client.GetMe(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "admin", user.Username)
	assert.Equal(t, "uid1", user.ID)
	assert.True(t, user.Active)
}

func TestListUsers(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/users.list", r.URL.Path)
		assert.Equal(t, "10", r.URL.Query().Get("count"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(UserListResponse{
			Users: []User{
				{ID: "u1", Username: "alice"},
				{ID: "u2", Username: "bob"},
			},
			Pagination: Pagination{Count: 2, Total: 2, Offset: 0},
		})
	})

	resp, err := client.ListUsers(context.Background(), ListOptions{Count: 10})
	require.NoError(t, err)
	assert.Len(t, resp.Users, 2)
	assert.Equal(t, "alice", resp.Users[0].Username)
	assert.Equal(t, 2, resp.Total)
}

func TestGetUserInfo(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/users.info", r.URL.Path)
		assert.Equal(t, "alice", r.URL.Query().Get("username"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"user": User{ID: "u1", Username: "alice", Name: "Alice Smith", Status: "online"},
		})
	})

	user, err := client.GetUserInfo(context.Background(), "alice")
	require.NoError(t, err)
	assert.Equal(t, "alice", user.Username)
	assert.Equal(t, "Alice Smith", user.Name)
}
