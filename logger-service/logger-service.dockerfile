# FROM golang:1.22.3 as builder

# WORKDIR /app

# COPY . .

# RUN CGO_ENABLED=0 go build -o loggerApp ./cmd/api

# RUN chmod +x /app/loggerApp

FROM alpine:latest

WORKDIR /app

# COPY --from=builder /app/loggerApp /app
COPY loggerApp .

CMD ["/app/loggerApp"]
