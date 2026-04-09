package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sancare/rocketchat-mcp/internal/rocketchat"
)

func registerChannelTools(srv *server.MCPServer, client *rocketchat.Client) {
	srv.AddTool(listChannelsTool, makeListChannelsHandler(client))
	srv.AddTool(listJoinedChannelsTool, makeListJoinedChannelsHandler(client))
	srv.AddTool(getChannelInfoTool, makeGetChannelInfoHandler(client))
	srv.AddTool(getChannelMessagesTool, makeGetChannelMessagesHandler(client))
}

var listChannelsTool = mcp.NewTool("list_channels",
	mcp.WithDescription("List public channels. Returns name, topic, member count, and message count."),
	mcp.WithNumber("count", mcp.Description("Items per page (default 100, max 200)")),
	mcp.WithNumber("offset", mcp.Description("Pagination offset")),
	mcp.WithString("query", mcp.Description("Filter channels by name")),
)

func makeListChannelsHandler(client *rocketchat.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		opts := listOptsFromRequest(req)
		resp, err := client.ListChannels(ctx, opts)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatChannelList(resp)), nil
	}
}

var listJoinedChannelsTool = mcp.NewTool("list_joined_channels",
	mcp.WithDescription("List channels the current user has joined."),
	mcp.WithNumber("count", mcp.Description("Items per page (default 100, max 200)")),
	mcp.WithNumber("offset", mcp.Description("Pagination offset")),
)

func makeListJoinedChannelsHandler(client *rocketchat.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		opts := listOptsFromRequest(req)
		resp, err := client.ListJoinedChannels(ctx, opts)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatChannelList(resp)), nil
	}
}

var getChannelInfoTool = mcp.NewTool("get_channel_info",
	mcp.WithDescription("Get detailed information about a channel by name or ID."),
	mcp.WithString("channel", mcp.Required(), mcp.Description("Channel name or ID")),
)

func makeGetChannelInfoHandler(client *rocketchat.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channel := stringArg(req, "channel")
		if channel == "" {
			return mcp.NewToolResultError("channel is required"), nil
		}
		ch, err := client.GetChannelInfo(ctx, channel)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatChannelInfo(ch)), nil
	}
}

var getChannelMessagesTool = mcp.NewTool("get_channel_messages",
	mcp.WithDescription("Get recent messages from a channel (newest first). Use the channel name."),
	mcp.WithString("channel", mcp.Required(), mcp.Description("Channel name or ID")),
	mcp.WithNumber("count", mcp.Description("Number of messages (default 50, max 200)")),
	mcp.WithNumber("offset", mcp.Description("Pagination offset")),
)

func makeGetChannelMessagesHandler(client *rocketchat.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channel := stringArg(req, "channel")
		if channel == "" {
			return mcp.NewToolResultError("channel is required"), nil
		}
		roomID, err := client.ResolveRoomID(ctx, channel)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		opts := listOptsFromRequest(req)
		resp, err := client.GetChannelHistory(ctx, roomID, opts)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(formatMessageList(channel, resp)), nil
	}
}

func formatChannelList(resp *rocketchat.ChannelListResponse) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Channels (%d/%d, offset %d):\n", resp.Count, resp.Total, resp.Offset)
	for _, ch := range resp.Channels {
		topic := ch.Topic
		if len(topic) > 80 {
			topic = topic[:80] + "..."
		}
		if topic != "" {
			fmt.Fprintf(&sb, "- #%s (%d members, %d msgs) — %s\n", ch.DisplayName(), ch.MembersCount, ch.MsgCount, topic)
		} else {
			fmt.Fprintf(&sb, "- #%s (%d members, %d msgs)\n", ch.DisplayName(), ch.MembersCount, ch.MsgCount)
		}
	}
	return sb.String()
}

func formatChannelInfo(ch *rocketchat.Channel) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Channel: #%s\n", ch.DisplayName())
	fmt.Fprintf(&sb, "  ID: %s\n", ch.ID)
	fmt.Fprintf(&sb, "  Type: %s\n", channelTypeName(ch.Type))
	fmt.Fprintf(&sb, "  Members: %d\n", ch.MembersCount)
	fmt.Fprintf(&sb, "  Messages: %d\n", ch.MsgCount)
	if ch.Topic != "" {
		fmt.Fprintf(&sb, "  Topic: %s\n", ch.Topic)
	}
	if ch.Description != "" {
		fmt.Fprintf(&sb, "  Description: %s\n", ch.Description)
	}
	fmt.Fprintf(&sb, "  Read-only: %v\n", ch.ReadOnly)
	fmt.Fprintf(&sb, "  Archived: %v\n", ch.Archived)
	return sb.String()
}

func channelTypeName(t string) string {
	switch t {
	case "c":
		return "public channel"
	case "p":
		return "private group"
	case "d":
		return "direct message"
	default:
		return t
	}
}
