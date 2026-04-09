package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sancare/rocketchat-mcp/internal/rocketchat"
)

func registerDMTools(srv *server.MCPServer, client *rocketchat.Client) {
	srv.AddTool(listDMsTool, makeListDMsHandler(client))
	srv.AddTool(getDMMessagesTool, makeGetDMMessagesHandler(client))
}

var listDMsTool = mcp.NewTool("list_dms",
	mcp.WithDescription("List direct message conversations for the current user."),
	mcp.WithNumber("count", mcp.Description("Items per page (default 50, max 200)")),
	mcp.WithNumber("offset", mcp.Description("Pagination offset")),
)

func makeListDMsHandler(client *rocketchat.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		opts := listOptsFromRequest(req)
		resp, err := client.ListDMs(ctx, opts)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "Direct messages (%d/%d, offset %d):\n", len(resp.IMs), resp.Total, resp.Offset)
		for _, dm := range resp.IMs {
			users := strings.Join(dm.Usernames, ", ")
			lastMsg := ""
			if dm.LastMessage != nil {
				lastMsg = truncate(dm.LastMessage.Text, 80)
			}
			if lastMsg != "" {
				fmt.Fprintf(&sb, "- %s (id: %s) — last: %s\n", users, dm.ID, lastMsg)
			} else {
				fmt.Fprintf(&sb, "- %s (id: %s)\n", users, dm.ID)
			}
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

var getDMMessagesTool = mcp.NewTool("get_dm_messages",
	mcp.WithDescription("Read messages from a direct message conversation. Use the username of the other person, or a DM room ID."),
	mcp.WithString("username", mcp.Required(), mcp.Description("Username of the other person (creates/opens the DM room)")),
	mcp.WithNumber("count", mcp.Description("Number of messages (default 50, max 200)")),
	mcp.WithNumber("offset", mcp.Description("Pagination offset")),
)

func makeGetDMMessagesHandler(client *rocketchat.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		username := stringArg(req, "username")
		if username == "" {
			return mcp.NewToolResultError("username is required"), nil
		}

		// im.create returns the existing DM room if it already exists.
		roomID, err := client.CreateDM(ctx, username)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("could not open DM with %q: %v", username, err)), nil
		}

		opts := listOptsFromRequest(req)
		resp, err := client.GetDMHistory(ctx, roomID, opts)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "DM with @%s (%d messages):\n", username, len(resp.Messages))
		for _, m := range resp.Messages {
			fmt.Fprintf(&sb, "- [%s] @%s: %s\n", m.Timestamp, m.User.Username, truncate(m.Text, 500))
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}
