# Mocklet MCP Server

Mocklet MCP Server is a [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server that allows AI assistants (such as Claude Desktop, Cursor, Windsurf, Zed, etc.) to natively interact with the Mocklet (Harmockery) mocking platform.

With this server, an AI assistant can seamlessly upload HAR files, create templates, and spawn disposable (ephemeral) mock servers out-of-the-box for autonomous frontend prototyping and E2E testing.

## Why use this?

- **Autonomous Test Generation:** The AI can independently spin up a mock server from a HAR file, write Cypress/Playwright tests to verify UI, and then tear down the mock.
- **Frontend Prototyping:** The AI agent can create a realistic backend from a Mocklet template in milliseconds and immediately start developing the dashboard or UI.
- **API Debugging:** Get mock usage statistics (hits/misses) and debug missing routes directly within the chat interface with your AI.

---

## 🛠 Installation

### Quick install (recommended)
The [`scripts/install.sh`](scripts/install.sh) script downloads the latest release binary for your OS/architecture and registers it with your client in one step:

```bash
curl -fsSL https://gitlab.com/keystr0ke/mocklet-mcp/-/raw/main/scripts/install.sh | bash -s -- claude --token "your_service_token_here"
curl -fsSL https://gitlab.com/keystr0ke/mocklet-mcp/-/raw/main/scripts/install.sh | bash -s -- codex  --token "your_service_token_here"
curl -fsSL https://gitlab.com/keystr0ke/mocklet-mcp/-/raw/main/scripts/install.sh | bash -s -- agy    --token "your_service_token_here"
```

Supports Linux and macOS (amd64/arm64). It picks up `MOCKLET_API_URL`/`MOCKLET_SERVICE_TOKEN` from the environment too, and falls back to printing the manual config snippet if the target client's CLI isn't found on `PATH`. See the client-specific sections below if you'd rather configure things by hand, or need Windows/Cursor/Zed.

### Requirements (build from source)
- [Go](https://go.dev/) version 1.22 or newer.

### Build from source
Clone the repository and build the binary:

```bash
git clone https://gitlab.com/keystr0ke/mocklet-mcp.git
cd mocklet-mcp
go mod tidy
go build -o mocklet-mcp .
```

This will create an executable `mocklet-mcp` binary in the current directory. Make sure to note the absolute path to this file (e.g., `/home/user/coding/mocklet-mcp/mocklet-mcp`), as you will need it to configure your clients.

---

## ⚙️ Configuration (Environment Variables)

The server requires the following environment variables to function properly:

| Variable | Description | Example |
| --- | --- | --- |
| `MOCKLET_API_URL` | The base URL of your Mocklet API. If omitted, defaults to `http://localhost:8080`. | `https://api.mocklet.dev` |
| `MOCKLET_SERVICE_TOKEN` | Service token (Bearer Token) for Mocklet API authentication. | `mckt_123456789...` |

---

## 🚀 Client Setup

### Claude Desktop
Open your Claude Desktop configuration file (usually located at `~/Library/Application Support/Claude/claude_desktop_config.json` on macOS or `%APPDATA%\Claude\claude_desktop_config.json` on Windows) and add the `mcpServers` section:

```json
{
  "mcpServers": {
    "mocklet": {
      "command": "/absolute/path/to/mocklet-mcp",
      "args": [],
      "env": {
        "MOCKLET_API_URL": "http://localhost:8080",
        "MOCKLET_SERVICE_TOKEN": "your_service_token_here"
      }
    }
  }
}
```
*Restart Claude Desktop after modifying the file.*

### Claude Code (CLI)
```bash
claude mcp add mocklet -s user -e MOCKLET_API_URL="http://localhost:8080" -e MOCKLET_SERVICE_TOKEN="your_service_token_here" -- /absolute/path/to/mocklet-mcp
```
`scripts/install.sh claude` does this for you.

### Cursor
1. Go to **Settings > Features > MCP**.
2. Click **+ Add New MCP Server**.
3. Type: `command`
4. Name: `mocklet`
5. Command: `MOCKLET_API_URL="http://localhost:8080" MOCKLET_SERVICE_TOKEN="your_token" /absolute/path/to/mocklet-mcp`

### Google Antigravity (AGY)
Add the server to your MCP configuration file at `~/.gemini/config/mcp_config.json` (older installs may still use the legacy path `~/.gemini/antigravity-cli/mcp_config.json`):

```json
{
  "mcpServers": {
    "mocklet": {
      "command": "/absolute/path/to/mocklet-mcp",
      "args": [],
      "env": {
        "MOCKLET_API_URL": "http://localhost:8080",
        "MOCKLET_SERVICE_TOKEN": "your_service_token_here"
      }
    }
  }
}
```
*Restart AGY after saving the file.* `scripts/install.sh agy` does this for you.

### Codex CLI
Codex's config is TOML, not JSON. Either run:

```bash
codex mcp add mocklet --env MOCKLET_API_URL="http://localhost:8080" --env MOCKLET_SERVICE_TOKEN="your_service_token_here" -- /absolute/path/to/mocklet-mcp
```

or edit `~/.codex/config.toml` directly:

```toml
[mcp_servers.mocklet]
command = "/absolute/path/to/mocklet-mcp"
env = { MOCKLET_API_URL = "http://localhost:8080", MOCKLET_SERVICE_TOKEN = "your_service_token_here" }
```

### Cline / VS Code
If you are using the Cline (formerly Claude Dev) extension in VS Code, edit its MCP settings file (e.g. `~/Library/Application Support/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json` on macOS):

```json
{
  "mcpServers": {
    "mocklet": {
      "command": "/absolute/path/to/mocklet-mcp",
      "args": [],
      "env": {
        "MOCKLET_API_URL": "http://localhost:8080",
        "MOCKLET_SERVICE_TOKEN": "your_service_token_here"
      }
    }
  }
}
```

### Zed IDE
For local usage in [Zed](https://zed.dev/), open your settings file (`~/.config/zed/settings.json`) and add the `context_servers` block:

```json
{
  "context_servers": {
    "mocklet": {
      "command": "/absolute/path/to/mocklet-mcp",
      "args": [],
      "env": {
        "MOCKLET_API_URL": "http://localhost:8080",
        "MOCKLET_SERVICE_TOKEN": "your_service_token_here"
      }
    }
  }
}
```

---

## 🧰 Available Tools

Once connected, the AI assistant will have access to the following operations:

- `mocklet_validate_har` — Validates a HAR file before deployment.
- `mocklet_create_mock` — Creates a disposable mock server from a HAR file.
- `mocklet_list_mocks` — Retrieves a list of active mocks.
- `mocklet_get_mock_stats` — Retrieves statistics (hits/misses) for a specific mock.
- `mocklet_delete_mock` — Stops and deletes a mock server.
- `mocklet_create_template` — Uploads a HAR file to create a reusable template.
- `mocklet_list_templates` — Searches for existing templates.
- `mocklet_spawn_mock` — Quickly launches an ephemeral mock based on an existing template.
- `mocklet_upload_template_revision` — Updates the logic of an existing template using a new HAR file.

## 💬 Built-in Prompts

The server also provides pre-configured prompts accessible in supported clients (like Claude) to automate popular tasks:
- **Spawn Dependency Mock**
- **Debug Mock Usage**
- **Generate Frontend with Mock Data**
- **Integration Testing Setup**
- **HAR Validation & Cleanup**
