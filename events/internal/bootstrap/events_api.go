package bootstrap

import (
	api "github.com/sashapremium/events/events/internal/api/events_service_api"
	"github.com/sashapremium/events/events/internal/services/eventsService"
)

func InitEventsServiceAPI(svc *eventsService.Service) *api.EventsServiceAPI {
	return api.NewEventsServiceAPI(svc)
}
