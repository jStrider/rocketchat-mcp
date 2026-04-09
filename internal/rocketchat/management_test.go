package rocketchat

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateChannel(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/channels.create", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "new-channel", body["name"])

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"channel": Channel{ID: "ch-new", Name: "new-channel"},
		})
	})

	ch, err := client.CreateChannel(context.Background(), "new-channel", nil, false)
	require.NoError(t, err)
	assert.Equal(t, "ch-new", ch.ID)
	assert.Equal(t, "new-channel", ch.Name)
}

func TestCreateChannel_WithMembers(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		members, ok := body["members"].([]any)
		assert.True(t, ok)
		assert.Len(t, members, 2)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"channel": Channel{ID: "ch1", Name: "team-channel"},
		})
	})

	ch, err := client.CreateChannel(context.Background(), "team-channel", []string{"alice", "bob"}, false)
	require.NoError(t, err)
	assert.Equal(t, "team-channel", ch.Name)
}

func TestArchiveChannel(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/channels.archive", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	})

	err := client.ArchiveChannel(context.Background(), "ch1")
	require.NoError(t, err)
}

func TestInviteToChannel(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/channels.invite", r.URL.Path)

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "ch1", body["roomId"])
		assert.Equal(t, "u1", body["userId"])

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	})

	err := client.InviteToChannel(context.Background(), "ch1", "u1")
	require.NoError(t, err)
}

func TestListGroups(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/groups.list", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"groups": []Channel{
				{ID: "g1", Name: "secret-team", Type: "p", MembersCount: 5},
			},
			"count":  1,
			"offset": 0,
			"total":  1,
		})
	})

	resp, err := client.ListGroups(context.Background(), ListOptions{})
	require.NoError(t, err)
	assert.Len(t, resp.Channels, 1)
	assert.Equal(t, "secret-team", resp.Channels[0].Name)
	assert.Equal(t, 1, resp.Total)
}
