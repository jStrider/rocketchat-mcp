package tools

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/sancare/rocketchat-mcp/internal/config"
	"github.com/sancare/rocketchat-mcp/internal/rocketchat"
)

// RegisterAll registers all MCP tools on the server.
func RegisterAll(srv *server.MCPServer, client *rocketchat.Client, _ *config.Config) {
	registerChannelTools(srv, client)
	registerMessageTools(srv, client)
	registerUserTools(srv, client)
	registerThreadTools(srv, client)
	registerManagementTools(srv, client)
}
