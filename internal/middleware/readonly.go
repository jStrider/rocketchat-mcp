package middleware

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// writeTools lists all tools that perform write operations.
var writeTools = map[string]bool{
	"send_message":      true,
	"send_dm":           true,
	"reply_to_thread":   true,
	"add_reaction":      true,
	"create_channel":    true,
	"archive_channel":   true,
	"invite_to_channel": true,
}

// ReadOnlyGuard returns a middleware that blocks write operations when enabled.
func ReadOnlyGuard(readOnly bool) server.ToolHandlerMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		if !readOnly {
			return next
		}
		return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if writeTools[req.Params.Name] {
				return mcp.NewToolResultError("server is in read-only mode: write operations are disabled"), nil
			}
			return next(ctx, req)
		}
	}
}

// IsWriteTool returns true if the tool name is a write operation.
func IsWriteTool(name string) bool {
	return writeTools[name]
}
