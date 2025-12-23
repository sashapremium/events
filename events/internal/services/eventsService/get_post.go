package eventsService

import (
	"context"
	"errors"

	eventmodel "github.com/sashapremium/events/events/internal/models"
)

func (s *Service) GetPost(ctx context.Context, id uint64) (*eventmodel.PostInfo, error) {
	if s.storage == nil {
		return nil, errors.New("storage is nil")
	}
	return s.storage.GetPost(ctx, id)
}
