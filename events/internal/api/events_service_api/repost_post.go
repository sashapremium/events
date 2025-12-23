package events_service_api

import (
	"context"

	"github.com/sashapremium/events/events/internal/pb/events_api"
)

func (a *EventsServiceAPI) RepostPost(ctx context.Context, req *events_api.RepostPostRequest) (*events_api.RepostPostResponse, error) {
	if err := a.svc.RepostPost(ctx, req.Id, req.UserHash); err != nil {
		return nil, err
	}
	return &events_api.RepostPostResponse{}, nil
}
