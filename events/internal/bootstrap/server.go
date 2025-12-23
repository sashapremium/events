package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	api "github.com/sashapremium/events/events/internal/api/events_service_api"
	"github.com/sashapremium/events/events/internal/pb/events_api"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcAddr = ":50051"
	httpAddr = ":8080"

	defaultSwaggerPath = "/swagger/events.swagger.json"
	swaggerEnv         = "swaggerPath"
)

func AppRun(apiSrv *api.EventsServiceAPI) {
	go func() {
		if err := runGRPCServer(apiSrv); err != nil {
			panic(fmt.Errorf("failed to run gRPC server: %v", err))
		}
	}()

	if err := runGatewayServer(); err != nil {
		panic(fmt.Errorf("failed to run gateway server: %v", err))
	}
}

func runGRPCServer(apiSrv events_api.EventsServiceServer) error {
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	events_api.RegisterEventsServiceServer(s, apiSrv)

	slog.Info("gRPC server listening", "addr", grpcAddr)
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

	if err := events_api.RegisterEventsServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts); err != nil {
		return err
	}

	r.Mount("/", mux)

	slog.Info("gRPC-Gateway server listening", "addr", httpAddr, "swagger", swaggerPath)
	return http.ListenAndServe(httpAddr, r)
}
