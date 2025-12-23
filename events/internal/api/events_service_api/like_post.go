package events_service_api

import (
	"context"

	"github.com/sashapremium/events/events/internal/pb/events_api"
)

func (a *EventsServiceAPI) LikePost(ctx context.Context, req *events_api.LikePostRequest) (*events_api.LikePostResponse, error) {
	if err := a.svc.LikePost(ctx, req.Id, req.UserHash); err != nil {
		return nil, err
	}
	return &events_api.LikePostResponse{}, nil
}
