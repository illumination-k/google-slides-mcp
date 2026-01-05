package main

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/illumination-k/google-slides-mcp/internal/domain"
)

type authChecker interface {
	Check(ctx context.Context) error
}

type slidesMetadataGetter interface {
	GetPresentationMetadata(ctx context.Context, presentationID string) (domain.Presentation, error)
}

type echoArgs struct {
	Text string `json:"text" jsonschema:"text to echo back"`
}

type slidesPingArgs struct {
	PresentationID string `json:"presentation_id" jsonschema:"Google Slides presentation ID to fetch"`
}

func newServer(slides slidesMetadataGetter, auth authChecker) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "google-slides-mcp", Version: "v0.0.0"}, nil)
	server.AddReceivingMiddleware(authCheckOnInitializeMiddleware(auth))

	mcp.AddTool(server, &mcp.Tool{
		Name:        "echo",
		Description: "Echo input text",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args echoArgs) (*mcp.CallToolResult, any, error) {
		_ = ctx
		_ = req
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: args.Text}},
		}, nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "slides_get_metadata",
		Description: "Get Google Slides presentation metadata (id, title) using Application Default Credentials (e.g., GOOGLE_APPLICATION_CREDENTIALS)",
	}, slidesGetMetadataHandler(slides))

	return server
}

func authCheckOnInitializeMiddleware(auth authChecker) mcp.Middleware {
	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
			if method == "initialize" {
				if err := auth.Check(ctx); err != nil {
					return nil, err
				}
			}
			return next(ctx, method, req)
		}
	}
}

func slidesGetMetadataHandler(slides slidesMetadataGetter) func(context.Context, *mcp.CallToolRequest, slidesPingArgs) (*mcp.CallToolResult, any, error) {
	return func(ctx context.Context, req *mcp.CallToolRequest, args slidesPingArgs) (*mcp.CallToolResult, any, error) {
		_ = req

		presentation, err := slides.GetPresentationMetadata(ctx, args.PresentationID)
		if err != nil {
			return nil, nil, err
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "ok: " + presentation.ID + " " + presentation.Title}},
		}, nil, nil
	}
}
