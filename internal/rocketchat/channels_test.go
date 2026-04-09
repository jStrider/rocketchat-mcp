package rocketchat

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListChannels(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/channels.list", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChannelListResponse{
			Channels: []Channel{
				{ID: "ch1", Name: "general", MembersCount: 42, Topic: "General chat"},
			},
			Pagination: Pagination{Count: 1, Total: 1, Offset: 0},
		})
	})

	resp, err := client.ListChannels(context.Background(), ListOptions{Count: 50})
	require.NoError(t, err)
	assert.Len(t, resp.Channels, 1)
	assert.Equal(t, "general", resp.Channels[0].Name)
	assert.Equal(t, 42, resp.Channels[0].MembersCount)
}

func TestListChannels_WithQuery(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "dev", r.URL.Query().Get("query"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChannelListResponse{
			Channels:   []Channel{{ID: "ch2", Name: "dev-ops"}},
			Pagination: Pagination{Count: 1, Total: 1},
		})
	})

	resp, err := client.ListChannels(context.Background(), ListOptions{Query: "dev"})
	require.NoError(t, err)
	assert.Equal(t, "dev-ops", resp.Channels[0].Name)
}

func TestListJoinedChannels(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/channels.list.joined", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChannelListResponse{
			Channels:   []Channel{{ID: "ch1", Name: "general"}},
			Pagination: Pagination{Count: 1, Total: 1},
		})
	})

	resp, err := client.ListJoinedChannels(context.Background(), ListOptions{})
	require.NoError(t, err)
	assert.Len(t, resp.Channels, 1)
}

func TestGetChannelInfo(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/channels.info", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"channel": Channel{ID: "ch1", Name: "general", Topic: "General discussion"},
		})
	})

	ch, err := client.GetChannelInfo(context.Background(), "general")
	require.NoError(t, err)
	assert.Equal(t, "general", ch.Name)
	assert.Equal(t, "General discussion", ch.Topic)
}

func TestGetChannelHistory(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/channels.history", r.URL.Path)
		assert.Equal(t, "ch1", r.URL.Query().Get("roomId"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessageListResponse{
			Messages: []Message{
				{ID: "m1", Text: "hello", User: MessageUser{Username: "alice"}},
				{ID: "m2", Text: "world", User: MessageUser{Username: "bob"}},
			},
			Pagination: Pagination{Count: 2, Total: 10},
		})
	})

	resp, err := client.GetChannelHistory(context.Background(), "ch1", ListOptions{Count: 50})
	require.NoError(t, err)
	assert.Len(t, resp.Messages, 2)
	assert.Equal(t, "hello", resp.Messages[0].Text)
}

func TestResolveRoomID(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"channel": Channel{ID: "resolved-id", Name: "general"},
		})
	})

	id, err := client.ResolveRoomID(context.Background(), "general")
	require.NoError(t, err)
	assert.Equal(t, "resolved-id", id)
}
