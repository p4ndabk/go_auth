# Build stage — CGO_ENABLED=0 works because the default sqlite driver
# (glebarez/sqlite) is pure Go; no gcc/musl-dev needed.
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/migrate ./cmd/migrate
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/seed ./cmd/seed

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /out/api ./api
COPY --from=builder /out/migrate ./migrate
COPY --from=builder /out/seed ./seed

ENV DB_DRIVER=sqlite
ENV DB_PATH=/data/go_auth.db
ENV GIN_MODE=release

EXPOSE 8080

CMD ["./api"]
