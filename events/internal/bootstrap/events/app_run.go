package events

import (
	"fmt"
	"net"

	"github.com/sashapremium/events/events/config"
	"github.com/sashapremium/events/events/internal/pb/events_api"
	"google.golang.org/grpc"
)

func AppRun(cfg *config.Config, api events_api.EventsServiceServer) error {
	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	events_api.RegisterEventsServiceServer(s, api)

	return s.Serve(lis)
}
