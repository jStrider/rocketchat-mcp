package middleware

import (
	"context"
	"log/slog"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Logging returns a middleware that logs every tool invocation with slog.
func Logging(logger *slog.Logger) server.ToolHandlerMiddleware {
	return func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
		return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			start := time.Now()
			logger.InfoContext(ctx, "tool_call_start",
				"tool", req.Params.Name,
			)

			result, err := next(ctx, req)
			duration := time.Since(start)

			if err != nil {
				logger.ErrorContext(ctx, "tool_call_error",
					"tool", req.Params.Name,
					"duration_ms", duration.Milliseconds(),
					"error", err.Error(),
				)
				return result, err
			}

			isError := result != nil && result.IsError
			logger.InfoContext(ctx, "tool_call_done",
				"tool", req.Params.Name,
				"duration_ms", duration.Milliseconds(),
				"is_error", isError,
			)
			return result, nil
		}
	}
}
