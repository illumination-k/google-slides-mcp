# google-slides-mcp

A minimal MCP server (stdio transport) that can:

- Validate Google auth on startup (ADC token acquisition)
- Fetch Google Slides presentation metadata via an MCP tool

## Prerequisites

- Go (pinned via `mise` in this repo)
- A Google Cloud project with **Google Slides API** enabled
- Application Default Credentials (ADC)
  - Recommended: set `GOOGLE_APPLICATION_CREDENTIALS` to a service account JSON
  - The service account must have access to the target presentation (share the slide with the service account email)

## Install toolchain (recommended)

```sh
mise install
```

## Run the MCP server (stdio)

```sh
GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json \
  go run ./cmd/google-slide-mcp
```

The server performs an auth preflight during the MCP `initialize` handshake. If ADC is not configured correctly, it will fail fast.

## Get presentation metadata (Slides API)

This repo exposes an MCP tool named `slides_get_metadata`.

- Tool name: `slides_get_metadata`
- Args:
  - `presentation_id` (string): Google Slides presentation ID

Example MCP tool call arguments:

```json
{
  "presentation_id": "YOUR_PRESENTATION_ID"
}
```

If authentication and permissions are correct, the tool returns something like:

```text
ok: <presentationId> <title>
```

## Dev commands

- Format: `mise run fmt`
- Lint: `mise run lint`
