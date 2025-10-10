FROM golang:1.25-alpine AS builder
ARG VERSION=1.1.41
ARG BUILD_DATE=10.10.2025

WORKDIR /app
COPY go.mod ./ 

RUN if [ -f go.sum ]; then cp go.sum .; fi
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-w -s -X main.version=$VERSION -X main.buildDate=$BUILD_DATE" \
  -o server ./cmd/main.go

FROM alpine:3.18
ARG VERSION=1.1.41
ARG BUILD_DATE=10.10.2025
LABEL version=$VERSION \
  build-date=$BUILD_DATE \
  maintainer="me@xomrkob.ru"
RUN addgroup -g 1000 appgroup && \
  adduser -D -u 1000 -G appgroup appuser
WORKDIR /app 
COPY --from=builder --chown=appuser:appgroup /app/server .
EXPOSE 8080

COPY docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod +x /app/docker-entrypoint.sh

USER appuser
CMD ["/app/server"]
