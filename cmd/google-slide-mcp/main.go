// mcp-stdio starts the MCP server over stdio.
package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{Name: "google-slides-mcp", Version: "v0.0.0"}, nil)

	type args struct {
		Text string `json:"text" jsonschema:"text to echo back"`
	}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "echo",
		Description: "Echo input text",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args args) (*mcp.CallToolResult, any, error) {
		_ = ctx
		_ = req
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: args.Text}},
		}, nil, nil
	})

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Printf("Server failed: %v", err)
	}
}
