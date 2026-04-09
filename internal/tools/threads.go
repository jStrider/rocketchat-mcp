package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sancare/rocketchat-mcp/internal/rocketchat"
)

func registerThreadTools(srv *server.MCPServer, client *rocketchat.Client) {
	srv.AddTool(getThreadMessagesTool, makeGetThreadMessagesHandler(client))
	srv.AddTool(replyToThreadTool, makeReplyToThreadHandler(client))
	srv.AddTool(addReactionTool, makeAddReactionHandler(client))
	srv.AddTool(listRoomFilesTool, makeListRoomFilesHandler(client))
}

var getThreadMessagesTool = mcp.NewTool("get_thread_messages",
	mcp.WithDescription("Get messages from a thread. The thread_id is the ID of the parent message."),
	mcp.WithString("thread_id", mcp.Required(), mcp.Description("ID of the parent message")),
	mcp.WithNumber("count", mcp.Description("Number of messages (default 50, max 200)")),
)

func makeGetThreadMessagesHandler(client *rocketchat.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		threadID := stringArg(req, "thread_id")
		if threadID == "" {
			return mcp.NewToolResultError("thread_id is required"), nil
		}
		opts := listOptsFromRequest(req)
		resp, err := client.GetThreadMessages(ctx, threadID, opts)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "Thread %s (%d messages):\n", threadID, len(resp.Messages))
		for _, m := range resp.Messages {
			fmt.Fprintf(&sb, "- [%s] @%s: %s\n", m.Timestamp, m.User.Username, truncate(m.Text, 500))
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

var replyToThreadTool = mcp.NewTool("reply_to_thread",
	mcp.WithDescription("Reply to a thread. Provide the thread_id (parent message ID) and optionally the channel."),
	mcp.WithString("thread_id", mcp.Required(), mcp.Description("ID of the parent message")),
	mcp.WithString("text", mcp.Required(), mcp.Description("Reply text")),
	mcp.WithString("channel", mcp.Description("Channel name or room ID (auto-detected if omitted)")),
)

func makeReplyToThreadHandler(client *rocketchat.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		threadID := stringArg(req, "thread_id")
		text := stringArg(req, "text")
		if threadID == "" || text == "" {
			return mcp.NewToolResultError("thread_id and text are required"), nil
		}

		roomID := stringArg(req, "channel")
		if roomID == "" {
			parent, err := client.GetMessage(ctx, threadID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("could not resolve thread room: %v", err)), nil
			}
			roomID = parent.RoomID
		} else {
			resolved, err := client.ResolveRoomID(ctx, roomID)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			roomID = resolved
		}

		msg, err := client.ReplyToThread(ctx, roomID, threadID, text)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Reply sent to thread %s (id: %s)", threadID, msg.ID)), nil
	}
}

var addReactionTool = mcp.NewTool("add_reaction",
	mcp.WithDescription("Add an emoji reaction to a message. Use colons around the emoji name, e.g. :thumbsup: or :white_check_mark:"),
	mcp.WithString("message_id", mcp.Required(), mcp.Description("ID of the message")),
	mcp.WithString("emoji", mcp.Required(), mcp.Description("Emoji name with colons, e.g. :thumbsup:")),
)

func makeAddReactionHandler(client *rocketchat.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		msgID := stringArg(req, "message_id")
		emoji := stringArg(req, "emoji")
		if msgID == "" || emoji == "" {
			return mcp.NewToolResultError("message_id and emoji are required"), nil
		}
		if err := client.AddReaction(ctx, msgID, emoji); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Reaction %s added to message %s", emoji, msgID)), nil
	}
}

var listRoomFilesTool = mcp.NewTool("list_room_files",
	mcp.WithDescription("List files uploaded to a channel or group."),
	mcp.WithString("channel", mcp.Required(), mcp.Description("Channel name or ID")),
	mcp.WithNumber("count", mcp.Description("Number of files (default 50, max 200)")),
	mcp.WithNumber("offset", mcp.Description("Pagination offset")),
)

func makeListRoomFilesHandler(client *rocketchat.Client) server.ToolHandlerFunc {
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
		resp, err := client.ListRoomFiles(ctx, roomID, opts)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		var sb strings.Builder
		fmt.Fprintf(&sb, "Files in #%s (%d/%d):\n", channel, len(resp.Files), resp.Total)
		for _, f := range resp.Files {
			fmt.Fprintf(&sb, "- %s (%s, %d bytes) by @%s\n", f.Name, f.Type, f.Size, f.User.Username)
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}
