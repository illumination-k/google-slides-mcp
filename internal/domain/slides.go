package domain

import "context"

type Presentation struct {
	ID    string
	Title string
}

type SlidesRepository interface {
	GetPresentation(ctx context.Context, presentationID string) (Presentation, error)
}

type SlidesService struct {
	repo SlidesRepository
}

func NewSlidesService(repo SlidesRepository) *SlidesService {
	return &SlidesService{repo: repo}
}

func (s *SlidesService) GetPresentationMetadata(ctx context.Context, presentationID string) (Presentation, error) {
	return s.repo.GetPresentation(ctx, presentationID)
}
