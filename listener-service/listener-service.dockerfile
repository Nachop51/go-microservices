# FROM golang:1.22.3 as builder

# WORKDIR /app

# COPY . .

# RUN CGO_ENABLED=0 go build -o listenerApp ./cmd/api

# RUN chmod +x /app/listenerApp

FROM alpine:latest

WORKDIR /app

# COPY --from=builder /app/listenerApp /app
COPY listenerApp .

CMD ["/app/listenerApp"]
