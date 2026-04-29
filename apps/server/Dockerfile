# Stage 1: Builder
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY apps/server/go.mod apps/server/go.sum ./
RUN go mod download

COPY apps/server/. .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/dist/main ./cmd/main.go

# Stage 2: Runtime
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/dist/main /app/main

EXPOSE 8080

CMD ["/app/main"]