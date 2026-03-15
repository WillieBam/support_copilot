# ---- build stage ----
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./backend/main.go

# ---- runtime stage ----
FROM alpine:3.21

RUN addgroup -S app && adduser -S app -G app
USER app

COPY --from=builder /server /server

EXPOSE 8080
ENTRYPOINT ["/server", "server"]
