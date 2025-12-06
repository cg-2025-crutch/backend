FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/grpc

FROM alpine:latest

COPY --from=builder /app/.env .
COPY --from=builder /app/main /main
COPY --from=builder /app/migrations ./migrations

CMD ["/main"]
