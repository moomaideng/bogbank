# Pinned Go hot-reload image. Bump AIR_VERSION when upgrading air.
FROM golang:1.27-alpine

ARG AIR_VERSION=1.62.0

RUN go install github.com/air-verse/air@v${AIR_VERSION}

WORKDIR /app

COPY docker/go-dev-entrypoint.sh /usr/local/bin/go-dev-entrypoint.sh
RUN chmod +x /usr/local/bin/go-dev-entrypoint.sh

ENTRYPOINT ["go-dev-entrypoint.sh"]
