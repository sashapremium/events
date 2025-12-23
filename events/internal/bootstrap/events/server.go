package events

import (
	"context"
	"log"
	"net"
	"net/http"

	api "github.com/sashapremium/events/events/internal/api/events_service_api"
	"github.com/sashapremium/events/events/internal/pb/events_api"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RunServer(apiSrv *api.EventsServiceAPI) error {
	go func() {
		if err := runGRPC(apiSrv); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	return runHTTPGateway()
}

func runGRPC(apiSrv events_api.EventsServiceServer) error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	events_api.RegisterEventsServiceServer(s, apiSrv)

	log.Println("gRPC listening on :50051")
	return s.Serve(lis)
}

func runHTTPGateway() error {
	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := events_api.RegisterEventsServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts); err != nil {
		return err
	}

	r := chi.NewRouter()
	r.Mount("/", mux)

	addr := ":8080"
	log.Println("HTTP gateway listening on", addr)
	return http.ListenAndServe(addr, r)
}
