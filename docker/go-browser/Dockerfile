  # docker/go-browser/Dockerfile
  FROM golang:1.25-alpine AS builder

  WORKDIR /app
  COPY go.mod go.sum ./
  RUN go mod download
  COPY . .

  ARG SERVICE
  RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd/${SERVICE}

  FROM debian:bookworm-slim

  RUN apt-get update \
      && apt-get install -y --no-install-recommends chromium ca-certificates \
      && rm -rf /var/lib/apt/lists/*

  ENV CHROME_BIN=/usr/bin/chromium

  COPY --from=builder --chmod=755 /bin/app /app
  ENTRYPOINT ["/app"]


