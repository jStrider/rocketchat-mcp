package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sancare/rocketchat-mcp/internal/rocketchat"
)

func registerUserTools(srv *server.MCPServer, client *rocketchat.Client) {
	srv.AddTool(listUsersTool, makeListUsersHandler(client))
	srv.AddTool(getUserInfoTool, makeGetUserInfoHandler(client))
	srv.AddTool(getMeTool, makeGetMeHandler(client))
}

var listUsersTool = mcp.NewTool("list_users",
	mcp.WithDescription("List users. Optionally filter by name or username."),
	mcp.WithNumber("count", mcp.Description("Items per page (default 100, max 200)")),
	mcp.WithNumber("offset", mcp.Description("Pagination offset")),
	mcp.WithString("query", mcp.Description("Filter by name or username")),
)

func makeListUsersHandler(client *rocketchat.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		opts := listOptsFromRequest(req)
		resp, err := client.ListUsers(ctx, opts)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "Users (%d/%d, offset %d):\n", resp.Count, resp.Total, resp.Offset)
		for _, u := range resp.Users {
			name := u.Name
			if name == "" {
				name = u.Username
			}
			fmt.Fprintf(&sb, "- @%s (%s) — %s\n", u.Username, name, u.Status)
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

var getUserInfoTool = mcp.NewTool("get_user_info",
	mcp.WithDescription("Get detailed information about a user by username."),
	mcp.WithString("username", mcp.Required(), mcp.Description("Username")),
)

func makeGetUserInfoHandler(client *rocketchat.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		username := stringArg(req, "username")
		if username == "" {
			return mcp.NewToolResultError("username is required"), nil
		}
		user, err := client.GetUserInfo(ctx, username)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "User: @%s\n", user.Username)
		fmt.Fprintf(&sb, "  ID: %s\n", user.ID)
		fmt.Fprintf(&sb, "  Name: %s\n", user.Name)
		fmt.Fprintf(&sb, "  Status: %s\n", user.Status)
		fmt.Fprintf(&sb, "  Active: %v\n", user.Active)
		if len(user.Roles) > 0 {
			fmt.Fprintf(&sb, "  Roles: %s\n", strings.Join(user.Roles, ", "))
		}
		return mcp.NewToolResultText(sb.String()), nil
	}
}

var getMeTool = mcp.NewTool("get_me",
	mcp.WithDescription("Get the currently authenticated user's profile."),
)

func makeGetMeHandler(client *rocketchat.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		user, err := client.GetMe(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "Authenticated as: @%s\n", user.Username)
		fmt.Fprintf(&sb, "  Name: %s\n", user.Name)
		fmt.Fprintf(&sb, "  ID: %s\n", user.ID)
		fmt.Fprintf(&sb, "  Status: %s\n", user.Status)
		return mcp.NewToolResultText(sb.String()), nil
	}
}
