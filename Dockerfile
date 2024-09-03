# Build Stage
FROM golang:1.22.5-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o server main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz

# Run Stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder ./app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY ./db/migration ./migration

EXPOSE 8080
CMD ["/app/server"]
ENTRYPOINT ["/app/start.sh"]