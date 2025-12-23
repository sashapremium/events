package events_service_api

import (
	"context"

	"github.com/sashapremium/events/events/internal/pb/events_api"
)

func (a *EventsServiceAPI) AddComment(ctx context.Context, req *events_api.AddCommentRequest) (*events_api.AddCommentResponse, error) {
	comment, err := a.svc.AddComment(ctx, req.Id, req.UserHash, req.Text)
	if err != nil {
		return nil, err
	}

	return &events_api.AddCommentResponse{
		Comment: comment,
	}, nil
}
