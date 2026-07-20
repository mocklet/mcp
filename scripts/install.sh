#!/usr/bin/env bash
# Downloads the latest mocklet-mcp release for the current OS/architecture
# and registers it as an MCP server with Claude Code, Codex, or Antigravity.
#
# Usage:
#   ./install.sh <claude|codex|agy> [--api-url URL] [--token TOKEN]
#
# Environment variables (used when the matching flag is omitted):
#   MOCKLET_API_URL       Mocklet API base URL (default: http://localhost:8080)
#   MOCKLET_SERVICE_TOKEN Mocklet service token
#   MOCKLET_MCP_INSTALL_DIR  Where to place the binary (default: ~/.local/share/mocklet-mcp)

set -euo pipefail

PROJECT_API="https://gitlab.com/api/v4/projects/84644902"

usage() {
  cat <<'EOF'
Usage: install.sh <claude|codex|agy> [--api-url URL] [--token TOKEN]

Downloads the latest mocklet-mcp release for your OS/architecture and
registers it as an MCP server with the chosen client:
  claude  Claude Code (uses `claude mcp add`)
  codex   OpenAI Codex CLI (uses `codex mcp add`)
  agy     Google Antigravity (writes ~/.gemini/config/mcp_config.json)

Options:
  --api-url URL   Mocklet API base URL (default: $MOCKLET_API_URL or http://localhost:8080)
  --token TOKEN   Mocklet service token (default: $MOCKLET_SERVICE_TOKEN)
EOF
}

CLIENT="${1:-}"
case "$CLIENT" in
  claude|codex|agy) shift ;;
  -h|--help|"") usage; exit 0 ;;
  *) echo "Unknown client: $CLIENT" >&2; usage; exit 1 ;;
esac

API_URL="${MOCKLET_API_URL:-http://localhost:8080}"
TOKEN="${MOCKLET_SERVICE_TOKEN:-}"

while [ $# -gt 0 ]; do
  case "$1" in
    --api-url) API_URL="$2"; shift 2 ;;
    --token) TOKEN="$2"; shift 2 ;;
    *) echo "Unknown option: $1" >&2; usage; exit 1 ;;
  esac
done

detect_platform() {
  local os arch
  case "$(uname -s)" in
    Linux) os=linux ;;
    Darwin) os=darwin ;;
    *)
      echo "Unsupported OS: $(uname -s). Download a binary manually from:" >&2
      echo "  https://gitlab.com/keystr0ke/mocklet-mcp/-/releases" >&2
      exit 1
      ;;
  esac
  case "$(uname -m)" in
    x86_64|amd64) arch=amd64 ;;
    arm64|aarch64) arch=arm64 ;;
    *)
      echo "Unsupported architecture: $(uname -m)." >&2
      exit 1
      ;;
  esac
  echo "${os}-${arch}"
}

fetch_latest_asset_url() {
  local platform="$1" release_json
  release_json="$(curl -fsSL "$PROJECT_API/releases/permalink/latest")"

  if command -v jq >/dev/null 2>&1; then
    printf '%s' "$release_json" | jq -r --arg suffix "mocklet-mcp-${platform}" \
      '.assets.links[] | select(.direct_asset_url | endswith($suffix)) | .direct_asset_url' | head -1
  elif command -v python3 >/dev/null 2>&1; then
    printf '%s' "$release_json" | python3 -c '
import json, sys
data = json.load(sys.stdin)
suffix = "mocklet-mcp-" + sys.argv[1]
for link in data["assets"]["links"]:
    if link["direct_asset_url"].endswith(suffix):
        print(link["direct_asset_url"])
        break
' "$platform"
  else
    printf '%s' "$release_json" | grep -oE '"direct_asset_url":"[^"]*mocklet-mcp-'"${platform}"'"' \
      | head -1 | sed -E 's/^"direct_asset_url":"//; s/"$//'
  fi
}

install_binary() {
  local platform asset_url install_dir bin_path
  platform="$(detect_platform)"
  echo "Detected platform: ${platform}" >&2

  asset_url="$(fetch_latest_asset_url "$platform")"
  if [ -z "$asset_url" ]; then
    echo "Could not find a release asset for platform '${platform}'." >&2
    echo "Check https://gitlab.com/keystr0ke/mocklet-mcp/-/releases" >&2
    exit 1
  fi

  install_dir="${MOCKLET_MCP_INSTALL_DIR:-$HOME/.local/share/mocklet-mcp}"
  bin_path="${install_dir}/mocklet-mcp"
  mkdir -p "$install_dir"

  echo "Downloading ${asset_url}" >&2
  curl -fsSL "$asset_url" -o "$bin_path"
  chmod +x "$bin_path"

  if [ "$(uname -s)" = "Darwin" ] && command -v xattr >/dev/null 2>&1; then
    xattr -d com.apple.quarantine "$bin_path" 2>/dev/null || true
  fi

  echo "Installed mocklet-mcp to ${bin_path}" >&2
  printf '%s' "$bin_path"
}

configure_claude() {
  local bin_path="$1"
  if ! command -v claude >/dev/null 2>&1; then
    cat >&2 <<EOF
'claude' CLI not found on PATH. Add the server manually, e.g.:
  claude mcp add mocklet -s user -e MOCKLET_API_URL="${API_URL}" -e MOCKLET_SERVICE_TOKEN="${TOKEN}" -- "${bin_path}"
EOF
    exit 1
  fi

  local args=(mcp add mocklet -s user -e "MOCKLET_API_URL=${API_URL}")
  if [ -n "$TOKEN" ]; then
    args+=(-e "MOCKLET_SERVICE_TOKEN=${TOKEN}")
  fi
  args+=(-- "$bin_path")

  claude "${args[@]}"
  echo "Registered 'mocklet' with Claude Code (user scope)." >&2
}

configure_codex() {
  local bin_path="$1"
  if ! command -v codex >/dev/null 2>&1; then
    cat >&2 <<EOF
'codex' CLI not found on PATH. Add the server manually, e.g.:
  codex mcp add mocklet --env MOCKLET_API_URL="${API_URL}" --env MOCKLET_SERVICE_TOKEN="${TOKEN}" -- "${bin_path}"

Or edit ~/.codex/config.toml directly:
  [mcp_servers.mocklet]
  command = "${bin_path}"
  env = { MOCKLET_API_URL = "${API_URL}", MOCKLET_SERVICE_TOKEN = "${TOKEN}" }
EOF
    exit 1
  fi

  local args=(mcp add mocklet --env "MOCKLET_API_URL=${API_URL}")
  if [ -n "$TOKEN" ]; then
    args+=(--env "MOCKLET_SERVICE_TOKEN=${TOKEN}")
  fi
  args+=(-- "$bin_path")

  codex "${args[@]}"
  echo "Registered 'mocklet' with Codex." >&2
}

configure_agy() {
  local bin_path="$1"
  local current="$HOME/.gemini/config/mcp_config.json"
  local legacy="$HOME/.gemini/antigravity-cli/mcp_config.json"
  local config="$current"
  if [ -f "$legacy" ] && [ ! -f "$current" ]; then
    config="$legacy"
  fi

  if ! command -v python3 >/dev/null 2>&1; then
    cat >&2 <<EOF
python3 not found, cannot safely edit JSON. Add this to ${config} by hand:
  {
    "mcpServers": {
      "mocklet": {
        "command": "${bin_path}",
        "args": [],
        "env": { "MOCKLET_API_URL": "${API_URL}", "MOCKLET_SERVICE_TOKEN": "${TOKEN}" }
      }
    }
  }
EOF
    exit 1
  fi

  mkdir -p "$(dirname "$config")"
  python3 -c '
import json, sys
path, bin_path, api_url, token = sys.argv[1:5]
try:
    with open(path) as f:
        data = json.load(f)
except (FileNotFoundError, json.JSONDecodeError):
    data = {}
data.setdefault("mcpServers", {})
env = {"MOCKLET_API_URL": api_url}
if token:
    env["MOCKLET_SERVICE_TOKEN"] = token
data["mcpServers"]["mocklet"] = {"command": bin_path, "args": [], "env": env}
with open(path, "w") as f:
    json.dump(data, f, indent=2)
    f.write("\n")
' "$config" "$bin_path" "$API_URL" "$TOKEN"

  echo "Registered 'mocklet' in ${config}. Restart Antigravity to pick it up." >&2
}

if [ -z "$TOKEN" ]; then
  echo "Warning: no service token set (--token or \$MOCKLET_SERVICE_TOKEN). You can add it later." >&2
fi

BIN_PATH="$(install_binary)"

case "$CLIENT" in
  claude) configure_claude "$BIN_PATH" ;;
  codex) configure_codex "$BIN_PATH" ;;
  agy) configure_agy "$BIN_PATH" ;;
esac
