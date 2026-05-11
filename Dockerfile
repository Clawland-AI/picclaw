# syntax=docker/dockerfile:1

ARG GO_VERSION=1.24

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /src

RUN apk add --no-cache ca-certificates git upx

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
ARG BUILD_TIME=unknown
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}" \
    -o /out/picoclaw \
    ./cmd/picoclaw
RUN upx --best --lzma /out/picoclaw

FROM alpine:3.20 AS runtime
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S picoclaw \
    && adduser -S -G picoclaw -h /home/picoclaw picoclaw

WORKDIR /home/picoclaw
COPY --from=builder /out/picoclaw /usr/local/bin/picoclaw

ARG PICOCLAW_GATEWAY_PORT=18790
ARG PICOCLAW_EDGE_PORT=9090
ENV PICOCLAW_AGENTS_DEFAULTS_WORKSPACE=/home/picoclaw/workspace \
    PICOCLAW_GATEWAY_HOST=0.0.0.0 \
    PICOCLAW_GATEWAY_PORT=${PICOCLAW_GATEWAY_PORT} \
    PICOCLAW_EDGE_PORT=${PICOCLAW_EDGE_PORT}

EXPOSE ${PICOCLAW_GATEWAY_PORT} ${PICOCLAW_EDGE_PORT}
USER picoclaw

ENTRYPOINT ["picoclaw"]
CMD ["gateway"]
