# FROM golang:1.22.3 as builder

# WORKDIR /app

# COPY . .

# RUN CGO_ENABLED=0 go build -o mailerApp ./cmd/api

# RUN chmod +x /app/mailerApp

FROM alpine:latest

WORKDIR /app

# COPY --from=builder /app/mailerApp /app
COPY mailerApp .
COPY templates ./templates

CMD ["/app/mailerApp"]
