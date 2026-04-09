package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sancare/rocketchat-mcp/internal/rocketchat"
)

func registerMessageTools(srv *server.MCPServer, client *rocketchat.Client) {
	srv.AddTool(searchMessagesTool, makeSearchMessagesHandler(client))
	srv.AddTool(sendMessageTool, makeSendMessageHandler(client))
	srv.AddTool(sendDMTool, makeSendDMHandler(client))
}

var searchMessagesTool = mcp.NewTool("search_messages",
	mcp.WithDescription("Search messages in a channel by text content."),
	mcp.WithString("channel", mcp.Required(), mcp.Description("Channel name or ID")),
	mcp.WithString("search_text", mcp.Required(), mcp.Description("Text to search for")),
	mcp.WithNumber("count", mcp.Description("Number of results (default 50, max 200)")),
)

func makeSearchMessagesHandler(client *rocketchat.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channel := stringArg(req, "channel")
		searchText := stringArg(req, "search_text")
		if channel == "" || searchText == "" {
			return mcp.NewToolResultError("channel and search_text are required"), nil
		}
		roomID, err := client.ResolveRoomID(ctx, channel)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		opts := listOptsFromRequest(req)
		resp, err := client.SearchMessages(ctx, roomID, searchText, opts)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "Search results for %q in #%s (%d found):\n", searchText, channel, len(resp.Messages))
		for _, m := range resp.Messages {
			fmt.Fprintf(&sb, "- [%s] @%s: %s\n", m.Timestamp, m.User.Username, truncate(m.Text, 500))
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

var sendMessageTool = mcp.NewTool("send_message",
	mcp.WithDescription("Send a message to a channel or group."),
	mcp.WithString("channel", mcp.Required(), mcp.Description("Channel name (without #) or ID")),
	mcp.WithString("text", mcp.Required(), mcp.Description("Message text")),
)

func makeSendMessageHandler(client *rocketchat.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channel := stringArg(req, "channel")
		text := stringArg(req, "text")
		if channel == "" || text == "" {
			return mcp.NewToolResultError("channel and text are required"), nil
		}
		msg, err := client.PostMessage(ctx, channel, text)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Message sent to #%s (id: %s)", channel, msg.ID)), nil
	}
}

var sendDMTool = mcp.NewTool("send_dm",
	mcp.WithDescription("Send a direct message to a user."),
	mcp.WithString("username", mcp.Required(), mcp.Description("Recipient username (without @)")),
	mcp.WithString("text", mcp.Required(), mcp.Description("Message text")),
)

func makeSendDMHandler(client *rocketchat.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		username := stringArg(req, "username")
		text := stringArg(req, "text")
		if username == "" || text == "" {
			return mcp.NewToolResultError("username and text are required"), nil
		}
		msg, err := client.SendDM(ctx, username, text)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("DM sent to @%s (id: %s)", username, msg.ID)), nil
	}
}
