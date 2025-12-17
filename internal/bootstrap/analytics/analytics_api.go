package analytics

import (
	"github.com/sashapremium/events/internal/api/analytics_service_api"
	analyticsService "github.com/sashapremium/events/internal/services/analyticsService"
)

func InitAnalyticsServiceAPI(svc *analyticsService.Service) *analytics_service_api.AnalyticsServiceAPI {
	return analytics_service_api.NewAnalyticsServiceAPI(svc)
}
