# Multi-stage build
FROM golang:1.23-alpine AS builder
ARG HTTP_PROXY
ARG HTTPS_PROXY
ENV HTTP_PROXY=${HTTP_PROXY} \
    HTTPS_PROXY=${HTTPS_PROXY} \
    http_proxy=${HTTP_PROXY} \
    https_proxy=${HTTPS_PROXY}
WORKDIR /app
# Install build deps
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/api

# Production stage for API
FROM alpine:3.18 AS production
WORKDIR /app
ENV PORT=5004 \
    HTTP_PROXY= \
    HTTPS_PROXY= \
    http_proxy= \
    https_proxy=
# Install curl for health checks
RUN apk add --no-cache curl ca-certificates
# Add non-root user
RUN adduser -D -g '' appuser
COPY --from=builder /app/server /app/server
COPY --from=builder /app/.env.example /app/.env
USER appuser
EXPOSE 5004
ENTRYPOINT ["/app/server"]

# gRPC stage (if needed)
FROM alpine:3.18 AS grpc
WORKDIR /app
ENV GRPC_PORT=9090 \
    HTTP_PROXY= \
    HTTPS_PROXY= \
    http_proxy= \
    https_proxy=
# Install grpc health probe
RUN apk add --no-cache curl ca-certificates
RUN adduser -D -g '' appuser
COPY --from=builder /app/server /app/server
COPY --from=builder /app/.env.example /app/.env
USER appuser
EXPOSE 9090
ENTRYPOINT ["/app/server"]
