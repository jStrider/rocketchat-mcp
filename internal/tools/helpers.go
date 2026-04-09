package tools

import (
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sancare/rocketchat-mcp/internal/rocketchat"
)

// stringArg extracts a string argument from the request.
func stringArg(req mcp.CallToolRequest, key string) string {
	args, _ := req.Params.Arguments.(map[string]any)
	if args == nil {
		return ""
	}
	v, _ := args[key].(string)
	return v
}

// boolArg extracts a boolean argument from the request.
func boolArg(req mcp.CallToolRequest, key string) bool {
	args, _ := req.Params.Arguments.(map[string]any)
	if args == nil {
		return false
	}
	v, _ := args[key].(bool)
	return v
}

// intArg extracts an integer argument with a default value.
func intArg(req mcp.CallToolRequest, key string, def int) int {
	args, _ := req.Params.Arguments.(map[string]any)
	if args == nil {
		return def
	}
	v, ok := args[key].(float64)
	if !ok {
		return def
	}
	return int(v)
}

// listOptsFromRequest extracts ListOptions from a tool request.
func listOptsFromRequest(req mcp.CallToolRequest) rocketchat.ListOptions {
	return rocketchat.ListOptions{
		Count:  intArg(req, "count", 0),
		Offset: intArg(req, "offset", 0),
		Query:  stringArg(req, "query"),
	}
}

// truncate shortens a string to maxLen and appends "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// formatMessageList formats a list of messages for LLM consumption.
func formatMessageList(channel string, resp *rocketchat.MessageListResponse) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Messages in #%s (%d messages):\n", channel, len(resp.Messages))
	for _, m := range resp.Messages {
		text := truncate(m.Text, 500)
		fmt.Fprintf(&sb, "- [%s] @%s: %s", m.Timestamp, m.User.Username, text)
		if m.ThreadCount > 0 {
			fmt.Fprintf(&sb, " [%d replies]", m.ThreadCount)
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
