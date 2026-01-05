package slidesapi

import (
	"context"

	"google.golang.org/api/option"
	slides "google.golang.org/api/slides/v1"

	"github.com/illumination-k/google-slides-mcp/internal/domain"
)

type Repository struct {
	svc *slides.Service
}

func NewRepository(ctx context.Context) (*Repository, error) {
	svc, err := slides.NewService(ctx, option.WithScopes(slides.PresentationsReadonlyScope))
	if err != nil {
		return nil, err
	}

	return &Repository{svc: svc}, nil
}

func (r *Repository) GetPresentation(ctx context.Context, presentationID string) (domain.Presentation, error) {
	call := r.svc.Presentations.Get(presentationID).Fields("presentationId,title").Context(ctx)
	p, err := call.Do()
	if err != nil {
		return domain.Presentation{}, err
	}

	return domain.Presentation{ID: p.PresentationId, Title: p.Title}, nil
}
