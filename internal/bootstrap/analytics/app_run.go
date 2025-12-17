package analytics

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sashapremium/events/config"
	"github.com/sashapremium/events/internal/pb/analytics_api"
)

type Consumer interface {
	Consume(ctx context.Context)
}

func AppRun(cfg *config.Config, api analytics_api.AnalyticsServiceServer, consumer Consumer) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Kafka consumer
	go consumer.Consume(ctx)

	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := InitGRPCServer(api)

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
		_ = lis.Close()
	}()

	log.Printf("analytics grpc server listening on %s", addr)
	return grpcServer.Serve(lis)
}
