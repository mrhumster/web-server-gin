FROM golang:1.25-alpine AS builder
ARG VERSION=1.1.8 
ARG BUILD_DATE=24.09.2025 

WORKDIR /app
COPY go.mod ./ 

RUN if [ -f go.sum ]; then cp go.sum .; fi
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-w -s -X main.version=$VERSION -X main.buildDate=$BUILD_DATE" \
  -o server ./cmd/main.go

RUN go install github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2

FROM alpine:3.18
ARG VERSION=1.1.8
ARG BUILD_DATE=24.09.2025
LABEL version=$VERSION \
  build-date=$BUILD_DATE \
  maintainer="me@xomrkob.ru"
# for HTTPS
# RUN apk --no-cache add ca-certificates
RUN addgroup -g 1000 appgroup && \
  adduser -D -u 1000 -G appgroup appuser
WORKDIR /app 
COPY --from=builder --chown=appuser:appgroup /app/server .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/migrations /app/migrations
USER appuser
EXPOSE 8080

COPY docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod +x /app/docker-entrypoint.sh

ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["/app/server"]
