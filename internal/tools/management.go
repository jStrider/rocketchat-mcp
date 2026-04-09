package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sancare/rocketchat-mcp/internal/rocketchat"
)

func registerManagementTools(srv *server.MCPServer, client *rocketchat.Client) {
	srv.AddTool(createChannelTool, makeCreateChannelHandler(client))
	srv.AddTool(archiveChannelTool, makeArchiveChannelHandler(client))
	srv.AddTool(inviteToChannelTool, makeInviteToChannelHandler(client))
	srv.AddTool(listGroupsTool, makeListGroupsHandler(client))
}

var createChannelTool = mcp.NewTool("create_channel",
	mcp.WithDescription("Create a new public channel."),
	mcp.WithString("name", mcp.Required(), mcp.Description("Channel name (lowercase, no spaces)")),
	mcp.WithString("members", mcp.Description("Comma-separated usernames to invite")),
	mcp.WithBoolean("read_only", mcp.Description("Create as read-only channel (default false)")),
)

func makeCreateChannelHandler(client *rocketchat.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := stringArg(req, "name")
		if name == "" {
			return mcp.NewToolResultError("name is required"), nil
		}
		membersStr := stringArg(req, "members")
		var members []string
		if membersStr != "" {
			for _, m := range strings.Split(membersStr, ",") {
				m = strings.TrimSpace(m)
				if m != "" {
					members = append(members, m)
				}
			}
		}
		readOnly := boolArg(req, "read_only")
		ch, err := client.CreateChannel(ctx, name, members, readOnly)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Channel #%s created (id: %s)", ch.Name, ch.ID)), nil
	}
}

var archiveChannelTool = mcp.NewTool("archive_channel",
	mcp.WithDescription("Archive a channel. This is a destructive operation — the channel will be read-only and hidden from channel list."),
	mcp.WithString("channel", mcp.Required(), mcp.Description("Channel name or ID")),
)

func makeArchiveChannelHandler(client *rocketchat.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channel := stringArg(req, "channel")
		if channel == "" {
			return mcp.NewToolResultError("channel is required"), nil
		}
		roomID, err := client.ResolveRoomID(ctx, channel)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if err := client.ArchiveChannel(ctx, roomID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Channel #%s archived", channel)), nil
	}
}

var inviteToChannelTool = mcp.NewTool("invite_to_channel",
	mcp.WithDescription("Invite a user to a channel."),
	mcp.WithString("channel", mcp.Required(), mcp.Description("Channel name or ID")),
	mcp.WithString("username", mcp.Required(), mcp.Description("Username to invite")),
)

func makeInviteToChannelHandler(client *rocketchat.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channel := stringArg(req, "channel")
		username := stringArg(req, "username")
		if channel == "" || username == "" {
			return mcp.NewToolResultError("channel and username are required"), nil
		}
		roomID, err := client.ResolveRoomID(ctx, channel)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		user, err := client.GetUserInfo(ctx, username)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("could not find user %q: %v", username, err)), nil
		}
		if err := client.InviteToChannel(ctx, roomID, user.ID); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("@%s invited to #%s", username, channel)), nil
	}
}

var listGroupsTool = mcp.NewTool("list_groups",
	mcp.WithDescription("List private groups the current user is a member of."),
	mcp.WithNumber("count", mcp.Description("Items per page (default 100, max 200)")),
	mcp.WithNumber("offset", mcp.Description("Pagination offset")),
)

func makeListGroupsHandler(client *rocketchat.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		opts := listOptsFromRequest(req)
		resp, err := client.ListGroups(ctx, opts)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatChannelList(resp)), nil
	}
}
