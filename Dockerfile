FROM golang:1.22 as builder

WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -v -o webhook-router

FROM debian:bookworm-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/webhook-router /app/webhook-router

EXPOSE 80
CMD ["/app/webhook-router"]
