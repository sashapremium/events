package events_service_api

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func userHashFromMetadata(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	// выбери один ключ и везде используй его
	if v := md.Get("x-user-hash"); len(v) > 0 {
		return v[0]
	}
	if v := md.Get("user-hash"); len(v) > 0 {
		return v[0]
	}

	return ""
}
