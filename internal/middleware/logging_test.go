package middleware

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeTestRequest(toolName string) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: toolName,
		},
	}
}

func TestLogging_Success(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	inner := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultText("ok"), nil
	}

	handler := Logging(logger)(inner)
	result, err := handler(context.Background(), makeTestRequest("list_channels"))

	require.NoError(t, err)
	assert.False(t, result.IsError)
	assert.Contains(t, buf.String(), "tool_call_start")
	assert.Contains(t, buf.String(), "tool_call_done")
	assert.Contains(t, buf.String(), "list_channels")
}

func TestLogging_Error(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	inner := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return nil, fmt.Errorf("connection refused")
	}

	handler := Logging(logger)(inner)
	_, err := handler(context.Background(), makeTestRequest("send_message"))

	require.Error(t, err)
	assert.Contains(t, buf.String(), "tool_call_error")
	assert.Contains(t, buf.String(), "connection refused")
}

func TestLogging_ToolError(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	inner := func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return mcp.NewToolResultError("channel not found"), nil
	}

	handler := Logging(logger)(inner)
	result, err := handler(context.Background(), makeTestRequest("get_channel_info"))

	require.NoError(t, err)
	assert.True(t, result.IsError)
	assert.Contains(t, buf.String(), `"is_error":true`)
}
