// mcp-stdio starts the MCP server over stdio.
package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/illumination-k/google-slides-mcp/internal/domain"
	"github.com/illumination-k/google-slides-mcp/internal/infra/googleauth"
	"github.com/illumination-k/google-slides-mcp/internal/infra/slidesapi"
)

func main() {
	ctx := context.Background()

	repo, err := slidesapi.NewRepository(ctx)
	if err != nil {
		log.Printf("Failed to initialize Google Slides repository: %v", err)
		return
	}

	slidesSvc := domain.NewSlidesService(repo)
	auth := googleauth.ADCChecker{Scopes: []string{"https://www.googleapis.com/auth/presentations.readonly"}}
	server := newServer(slidesSvc, auth)

	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Printf("Server failed: %v", err)
	}
}
