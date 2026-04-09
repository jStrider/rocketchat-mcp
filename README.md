# rocketchat-mcp

MCP server for Rocket.Chat in Go. Exposes channels, messages, threads, DMs, reactions, users, and files as MCP tools.

## Quick start

```bash
make build

# stdio transport (Claude Code)
ROCKETCHAT_URL=https://chat.example.com \
ROCKETCHAT_USER_ID=your-user-id \
ROCKETCHAT_AUTH_TOKEN=your-token \
  ./dist/rocketchat-mcp -transport stdio

# HTTP transport (Docker)
docker build -t rocketchat-mcp .
docker run -d -p 8080:8080 \
  -e ROCKETCHAT_URL=https://chat.example.com \
  -e ROCKETCHAT_AUTH_TOKEN=your-token \
  -e ROCKETCHAT_USER_ID=your-user-id \
  rocketchat-mcp
```

Get token and User ID: Rocket.Chat > Profile > Personal Access Tokens.

## Claude Code configuration

### stdio (recommended)

```json
{
  "mcpServers": {
    "rocketchat": {
      "command": "/path/to/rocketchat-mcp",
      "args": ["-transport", "stdio"],
      "env": {
        "ROCKETCHAT_URL": "https://chat.example.com",
        "ROCKETCHAT_USER_ID": "your-user-id",
        "ROCKETCHAT_AUTH_TOKEN": "your-token"
      }
    }
  }
}
```

### HTTP

```json
{
  "mcpServers": {
    "rocketchat": {
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

## Environment variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `ROCKETCHAT_URL` | yes | - | Rocket.Chat server URL (must be HTTPS) |
| `ROCKETCHAT_USER_ID` | yes | - | User ID for API auth |
| `ROCKETCHAT_AUTH_TOKEN` | yes | - | Personal Access Token |
| `ROCKETCHAT_READ_ONLY` | no | `false` | Block all write operations |
| `ROCKETCHAT_ALLOW_HTTP` | no | `false` | Allow non-HTTPS URLs (dev only) |
| `ROCKETCHAT_LOG_LEVEL` | no | `info` | Log level: debug, info, warn, error |
| `ROCKETCHAT_MCP_ADDR` | no | `:8080` | HTTP listen address |

## Tools

| Group | Tools |
|-------|-------|
| Channels | `list_channels` `list_joined_channels` `get_channel_info` `get_channel_messages` |
| Messages | `search_messages` `send_message` `send_dm` |
| Users | `list_users` `get_user_info` `get_me` |
| Threads | `get_thread_messages` `reply_to_thread` |
| Reactions | `add_reaction` |
| Files | `list_room_files` |
| Management | `create_channel` `archive_channel` `invite_to_channel` `list_groups` |

## Development

```bash
make help     # show all targets
make build    # build binary
make test     # run tests with coverage
make lint     # run golangci-lint
make check    # fmt + vet + lint + test
make docker   # build Docker image
make secrets  # scan for leaked secrets
```

## License

MIT
