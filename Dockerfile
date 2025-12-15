# Multi-stage build for a small production image
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o inventory-service .

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/inventory-service ./inventory-service
EXPOSE 8080
ENTRYPOINT ["/app/inventory-service"]
