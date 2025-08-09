# Multi-stage build
FROM golang:1.22-alpine AS builder
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
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/api

FROM alpine:3.19
WORKDIR /app
ENV PORT=5001 \
    HTTP_PROXY= \
    HTTPS_PROXY= \
    http_proxy= \
    https_proxy=
# Add non-root user
RUN adduser -D -g '' appuser
COPY --from=builder /app/server /app/server
COPY --from=builder /app/.env.example /app/.env
USER appuser
EXPOSE 5001
ENTRYPOINT ["/app/server"]
