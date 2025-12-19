package events_service_api

import (
	"context"
	"time"

	eventmodel "github.com/sashapremium/events/internal/models"
	"github.com/sashapremium/events/internal/pb/events_api"
	pbmodels "github.com/sashapremium/events/internal/pb/models"
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

func (a *EventsServiceAPI) GetPost(ctx context.Context, req *events_api.GetPostRequest) (*events_api.GetPostResponse, error) {
	post, err := a.svc.GetPost(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if userHash := userHashFromMetadata(ctx); userHash != "" {
		_ = a.svc.ViewPost(ctx, req.Id, userHash)
	}

	return &events_api.GetPostResponse{
		Post: mapPost(post),
	}, nil
}

func mapPost(p *eventmodel.PostInfo) *pbmodels.PostModel {
	if p == nil {
		return nil
	}

	return &pbmodels.PostModel{
		Id:          p.ID,
		Title:       p.Title,
		AuthorId:    p.AuthorID,
		Category:    p.Category,
		Content:     p.Content,
		PublishedAt: p.PublishedAt.Format(time.RFC3339),
	}
}
