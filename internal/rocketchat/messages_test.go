package rocketchat

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchMessages(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/chat.search", r.URL.Path)
		assert.Equal(t, "room1", r.URL.Query().Get("roomId"))
		assert.Equal(t, "deploy", r.URL.Query().Get("searchText"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessageListResponse{
			Messages: []Message{
				{ID: "m1", Text: "deploy v2.0", User: MessageUser{Username: "ops"}},
			},
			Pagination: Pagination{Count: 1, Total: 1},
		})
	})

	resp, err := client.SearchMessages(context.Background(), "room1", "deploy", ListOptions{})
	require.NoError(t, err)
	assert.Len(t, resp.Messages, 1)
	assert.Contains(t, resp.Messages[0].Text, "deploy")
}

func TestPostMessage(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/chat.postMessage", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, "general", body["channel"])
		assert.Equal(t, "hello world", body["text"])

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"message": Message{ID: "m1", Text: "hello world", RoomID: "ch1"},
		})
	})

	msg, err := client.PostMessage(context.Background(), "general", "hello world")
	require.NoError(t, err)
	assert.Equal(t, "m1", msg.ID)
	assert.Equal(t, "hello world", msg.Text)
}

func TestCreateDM(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/im.create", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"room": Room{ID: "dm-room-id"},
		})
	})

	roomID, err := client.CreateDM(context.Background(), "alice")
	require.NoError(t, err)
	assert.Equal(t, "dm-room-id", roomID)
}

func TestCreateDM_EmptyRoomID(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"room": Room{ID: ""},
		})
	})

	_, err := client.CreateDM(context.Background(), "nobody")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "empty room ID")
}

func TestGetMessage(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/chat.getMessage", r.URL.Path)
		assert.Equal(t, "msg123", r.URL.Query().Get("msgId"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"message": Message{ID: "msg123", RoomID: "ch1", Text: "test msg"},
		})
	})

	msg, err := client.GetMessage(context.Background(), "msg123")
	require.NoError(t, err)
	assert.Equal(t, "msg123", msg.ID)
	assert.Equal(t, "ch1", msg.RoomID)
}

func TestAddReaction(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/chat.react", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true}`))
	})

	err := client.AddReaction(context.Background(), "msg1", ":thumbsup:")
	require.NoError(t, err)
}
