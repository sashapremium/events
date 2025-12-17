package events_service_api

import (
	"context"

	"github.com/sashapremium/events/internal/pb/events_api"
)

func (a *EventsServiceAPI) LikePost(ctx context.Context, req *events_api.LikePostRequest) (*events_api.LikePostResponse, error) {
	if err := a.svc.LikePost(ctx, req.Id, req.UserHash); err != nil {
		return nil, err
	}
	return &events_api.LikePostResponse{}, nil
}

func (a *EventsServiceAPI) RepostPost(ctx context.Context, req *events_api.RepostPostRequest) (*events_api.RepostPostResponse, error) {
	if err := a.svc.RepostPost(ctx, req.Id, req.UserHash); err != nil {
		return nil, err
	}
	return &events_api.RepostPostResponse{}, nil
}

func (a *EventsServiceAPI) AddComment(ctx context.Context, req *events_api.AddCommentRequest) (*events_api.AddCommentResponse, error) {
	comment, err := a.svc.AddComment(ctx, req.Id, req.UserHash, req.Text)
	if err != nil {
		return nil, err
	}

	return &events_api.AddCommentResponse{
		Comment: comment,
	}, nil
}
