package analytics

import (
	"github.com/sashapremium/events/analytics/internal/api/analytics_service_api"
	analyticsService "github.com/sashapremium/events/analytics/internal/services/analyticsService"
)

func InitAnalyticsServiceAPI(svc *analyticsService.Service) *analytics_service_api.AnalyticsServiceAPI {
	return analytics_service_api.NewAnalyticsServiceAPI(svc)
}
