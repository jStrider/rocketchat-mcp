package rocketchat

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *Client) {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	client := NewClient(ts.URL, "test-uid", "test-token", 10*1024*1024)
	return ts, client
}

func TestClient_AuthHeaders(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test-token", r.Header.Get("X-Auth-Token"))
		assert.Equal(t, "test-uid", r.Header.Get("X-User-Id"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})

	_, err := client.get(context.Background(), "/test", nil)
	require.NoError(t, err)
}

func TestClient_ContextCancellation(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.get(ctx, "/test", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestClient_HTTPError401(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"status":"error","message":"You must be logged in to do this."}`))
	})

	_, err := client.get(context.Background(), "/test", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "You must be logged in")
}

func TestClient_HTTPError401_NoBody(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{}`))
	})

	_, err := client.get(context.Background(), "/test", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestClient_HTTPError403(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{}`))
	})

	_, err := client.get(context.Background(), "/test", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestClient_HTTPErrorWithAPIMessage(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"success":false,"error":"Channel_Not_Found"}`))
	})

	_, err := client.get(context.Background(), "/channels.info", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Channel_Not_Found")
}

func TestClient_PostSendsJSON(t *testing.T) {
	_, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})

	_, err := client.post(context.Background(), "/chat.postMessage", map[string]string{"channel": "general", "text": "hello"})
	require.NoError(t, err)
}

func TestPagingParams_Defaults(t *testing.T) {
	p := pagingParams(ListOptions{}, 50)
	assert.Equal(t, "50", p.Get("count"))
	assert.Equal(t, "0", p.Get("offset"))
}

func TestPagingParams_MaxCap(t *testing.T) {
	p := pagingParams(ListOptions{Count: 999}, 50)
	assert.Equal(t, "200", p.Get("count"))
}

func TestPagingParams_NegativeOffset(t *testing.T) {
	p := pagingParams(ListOptions{Offset: -5}, 50)
	assert.Equal(t, "0", p.Get("offset"))
}
