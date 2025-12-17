package events_service_api

import (
	"context"
	"time"

	eventmodel "github.com/sashapremium/events/internal/models"
	"github.com/sashapremium/events/internal/pb/events_api"
	pbmodels "github.com/sashapremium/events/internal/pb/models"
)

func (a *EventsServiceAPI) GetPost(ctx context.Context, req *events_api.GetPostRequest) (*events_api.GetPostResponse, error) {
	post, err := a.svc.GetPost(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	_ = a.svc.ViewPost(ctx, req.Id, "")

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
