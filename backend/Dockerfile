# ---- build stage ----
FROM docker.io/library/golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY db/ ./db/
COPY backend/ ./backend/
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./backend/main.go

# ---- runtime stage ----
FROM alpine:3.21

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /server ./server
RUN chown -R app:app /app
USER app

EXPOSE 8080
ENTRYPOINT ["./server", "server"]
