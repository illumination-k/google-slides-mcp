package main

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/illumination-k/google-slides-mcp/internal/domain"
)

type mockSlidesMetadataGetter struct {
	err          error
	presentation domain.Presentation
}

func (m mockSlidesMetadataGetter) GetPresentationMetadata(ctx context.Context, presentationID string) (domain.Presentation, error) {
	_ = ctx
	_ = presentationID
	if m.err != nil {
		return domain.Presentation{}, m.err
	}
	return m.presentation, nil
}

func TestSlidesGetMetadataHandler_OK(t *testing.T) {
	h := slidesGetMetadataHandler(mockSlidesMetadataGetter{presentation: domain.Presentation{ID: "pid", Title: "hello"}})

	res, _, err := h(context.Background(), &mcp.CallToolRequest{}, slidesPingArgs{PresentationID: "pid"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res == nil {
		t.Fatalf("expected result")
	}
	if len(res.Content) != 1 {
		t.Fatalf("expected 1 content item, got %d", len(res.Content))
	}

	text, ok := res.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent")
	}
	if text.Text != "ok: pid hello" {
		t.Fatalf("unexpected text: %q", text.Text)
	}
}

type mockAuthChecker struct {
	err   error
	calls int
}

func (m *mockAuthChecker) Check(ctx context.Context) error {
	_ = ctx
	m.calls++
	return m.err
}

func TestAuthCheckOnInitializeMiddleware_ChecksOnlyOnInitialize(t *testing.T) {
	auth := &mockAuthChecker{}
	mw := authCheckOnInitializeMiddleware(auth)

	nextCalls := 0
	h := mw(func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
		_ = ctx
		_ = method
		_ = req
		nextCalls++
		return nil, nil
	})

	if _, err := h(context.Background(), "tools/list", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if auth.calls != 0 {
		t.Fatalf("expected 0 auth checks, got %d", auth.calls)
	}
	if nextCalls != 1 {
		t.Fatalf("expected next to be called once, got %d", nextCalls)
	}

	if _, err := h(context.Background(), "initialize", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if auth.calls != 1 {
		t.Fatalf("expected 1 auth check, got %d", auth.calls)
	}
	if nextCalls != 2 {
		t.Fatalf("expected next to be called twice, got %d", nextCalls)
	}
}

func TestAuthCheckOnInitializeMiddleware_BlocksOnAuthError(t *testing.T) {
	auth := &mockAuthChecker{err: context.Canceled}
	mw := authCheckOnInitializeMiddleware(auth)

	nextCalls := 0
	h := mw(func(ctx context.Context, method string, req mcp.Request) (mcp.Result, error) {
		_ = ctx
		_ = method
		_ = req
		nextCalls++
		return nil, nil
	})

	if _, err := h(context.Background(), "initialize", nil); err == nil {
		t.Fatalf("expected error")
	}
	if auth.calls != 1 {
		t.Fatalf("expected 1 auth check, got %d", auth.calls)
	}
	if nextCalls != 0 {
		t.Fatalf("expected next not to be called, got %d", nextCalls)
	}
}
