# FROM golang:1.22.3 as builder

# WORKDIR /app

# COPY . .

# RUN CGO_ENABLED=0 go build -o authApp ./cmd/api

# RUN chmod +x /app/authApp

FROM alpine:latest

WORKDIR /app

# COPY --from=builder /app/authApp /app
COPY authApp .

CMD ["/app/authApp"]

