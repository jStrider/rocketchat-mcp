package middleware

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadOnlyGuard_BlocksWriteTools(t *testing.T) {
	called := false
	inner := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		called = true
		return mcp.NewToolResultText("ok"), nil
	}

	handler := ReadOnlyGuard(true)(inner)

	for _, tool := range []string{"send_message", "send_dm", "reply_to_thread", "add_reaction", "create_channel", "archive_channel", "invite_to_channel"} {
		called = false
		result, err := handler(context.Background(), makeTestRequest(tool))
		require.NoError(t, err, "tool: %s", tool)
		assert.True(t, result.IsError, "tool %s should be blocked", tool)
		assert.False(t, called, "inner handler should not be called for %s", tool)
	}
}

func TestReadOnlyGuard_AllowsReadTools(t *testing.T) {
	called := false
	inner := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		called = true
		return mcp.NewToolResultText("ok"), nil
	}

	handler := ReadOnlyGuard(true)(inner)

	for _, tool := range []string{"list_channels", "get_channel_info", "search_messages", "list_users", "get_me"} {
		called = false
		result, err := handler(context.Background(), makeTestRequest(tool))
		require.NoError(t, err, "tool: %s", tool)
		assert.False(t, result.IsError, "tool %s should be allowed", tool)
		assert.True(t, called, "inner handler should be called for %s", tool)
	}
}

func TestReadOnlyGuard_DisabledPassesThrough(t *testing.T) {
	called := false
	inner := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		called = true
		return mcp.NewToolResultText("ok"), nil
	}

	handler := ReadOnlyGuard(false)(inner)

	result, err := handler(context.Background(), makeTestRequest("send_message"))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	assert.True(t, called)
}

func TestIsWriteTool(t *testing.T) {
	assert.True(t, IsWriteTool("send_message"))
	assert.True(t, IsWriteTool("send_dm"))
	assert.False(t, IsWriteTool("list_channels"))
	assert.False(t, IsWriteTool("get_me"))
}
