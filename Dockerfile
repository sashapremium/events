# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS base
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

# --- build events ---
FROM base AS events-builder
COPY events ./events
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/eventsService ./events/cmd/app

# --- build analytics ---
FROM base AS analytics-builder
COPY analytics ./analytics
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/analyticsService ./analytics/cmd/app

# --- runtime images ---
FROM alpine:3.20 AS events
WORKDIR /app
COPY --from=events-builder /out/eventsService /app/eventsService
ENTRYPOINT ["/app/eventsService"]

FROM alpine:3.20 AS analytics
WORKDIR /app
COPY --from=analytics-builder /out/analyticsService /app/analyticsService
ENTRYPOINT ["/app/analyticsService"]
