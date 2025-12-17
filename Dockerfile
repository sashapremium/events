# syntax=docker/dockerfile:1

FROM golang:alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o eventsService ./cmd/app


FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/events /app/events

EXPOSE 8080 50051

ENV POSTGRES_DSN=""
ENV KAFKA_BROKERS=""
ENV KAFKA_TOPIC="content_events"

CMD ["/app/events"]
