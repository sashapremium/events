package events_service_api

import (
	"context"

	eventmodel "github.com/sashapremium/events/events/internal/models"
	"github.com/sashapremium/events/events/internal/pb/events_api"
	pbmodels "github.com/sashapremium/events/events/internal/pb/models"
)

type EventsService interface {
	GetPost(ctx context.Context, id uint64) (*eventmodel.PostInfo, error)
	ViewPost(ctx context.Context, id uint64, userHash string) error
	LikePost(ctx context.Context, id uint64, userHash string) error
	RepostPost(ctx context.Context, id uint64, userHash string) error
	AddComment(ctx context.Context, id uint64, userHash, text string) (*pbmodels.CommentModel, error)
}

type EventsServiceAPI struct {
	events_api.UnimplementedEventsServiceServer
	svc EventsService
}

func NewEventsServiceAPI(svc EventsService) *EventsServiceAPI {
	return &EventsServiceAPI{svc: svc}
}
