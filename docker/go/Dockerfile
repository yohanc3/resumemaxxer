# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /bin/app ./cmd/${SERVICE}

FROM gcr.io/distroless/static-debian12

COPY --from=builder --chmod=755 /bin/app /app
COPY --from=builder /app/cmd/dbsetup/migration_files /migration_files

ENTRYPOINT ["/app"]

