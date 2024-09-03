# Build Stage
FROM golang:1.22.5-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o server main.go

# Run Stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
COPY app.env .

EXPOSE 8080
CMD ["/app/server"]