FROM golang:1.25-alpine AS builder


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o eventsService ./cmd/app/events

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o analyticsService ./cmd/app/analytics

FROM alpine:3.19

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/eventsService /app/eventsService
COPY --from=builder /app/analyticsService /app/analyticsService

COPY config/config.yaml /config/config.yaml

EXPOSE 8080