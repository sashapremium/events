package events

import (
	api "github.com/sashapremium/events/internal/api/events_service_api"
	"github.com/sashapremium/events/internal/services/eventsService"
)

func InitEventsServiceAPI(svc *eventsService.Service) *api.EventsServiceAPI {
	return api.NewEventsServiceAPI(svc)
}
