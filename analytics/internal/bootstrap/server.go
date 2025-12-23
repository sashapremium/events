package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/sashapremium/events/analytics/internal/pb/analytics_api"
)

const (
	grpcListenAddr     = ":50052"
	grpcDialAddr       = "localhost:50052"
	httpAddr           = ":8082"
	defaultSwaggerPath = "/swagger/analytics.swagger.json"
	swaggerEnv         = "swaggerPath"
)

type Consumer interface {
	Consume(ctx context.Context)
}

func AppRun(apiSrv analytics_api.AnalyticsServiceServer, consumer Consumer) {
	go func() {
		slog.Info("analytics consumer starting")
		consumer.Consume(context.Background())
	}()

	go func() {
		if err := runGRPCServer(apiSrv); err != nil {
			panic(fmt.Errorf("failed to run gRPC server: %v", err))
		}
	}()

	if err := runGatewayServer(); err != nil {
		panic(fmt.Errorf("failed to run gateway server: %v", err))
	}
}

func runGRPCServer(apiSrv analytics_api.AnalyticsServiceServer) error {
	lis, err := net.Listen("tcp", grpcListenAddr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	analytics_api.RegisterAnalyticsServiceServer(s, apiSrv)

	slog.Info("analytics gRPC server listening", "addr", grpcListenAddr)
	return s.Serve(lis)
}

func runGatewayServer() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	swaggerPath := os.Getenv(swaggerEnv)
	if swaggerPath == "" {
		swaggerPath = defaultSwaggerPath
	}

	if _, err := os.Stat(swaggerPath); err != nil {
		return fmt.Errorf("swagger file not found: %s, err=%v", swaggerPath, err)
	}

	r := chi.NewRouter()

	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, swaggerPath)
	})

	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := analytics_api.RegisterAnalyticsServiceHandlerFromEndpoint(ctx, mux, grpcDialAddr, opts); err != nil {
		return err
	}

	r.Mount("/", mux)

	slog.Info("analytics gRPC-Gateway server listening", "addr", httpAddr, "swagger", swaggerPath)
	return http.ListenAndServe(httpAddr, r)
}
