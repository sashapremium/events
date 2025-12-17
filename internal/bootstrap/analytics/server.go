package analytics

import (
	"github.com/sashapremium/events/internal/pb/analytics_api"
	"google.golang.org/grpc"
)

func InitGRPCServer(api analytics_api.AnalyticsServiceServer) *grpc.Server {
	s := grpc.NewServer()
	analytics_api.RegisterAnalyticsServiceServer(s, api)
	return s
}
