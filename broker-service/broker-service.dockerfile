# FROM golang:1.22.3 as builder

# WORKDIR /app

# COPY . .

# RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

# RUN chmod +x /app/brokerApp

FROM alpine:latest

WORKDIR /app

# COPY --from=builder /app/brokerApp /app
COPY brokerApp .

CMD ["/app/brokerApp"]
