FROM golang:1.15-alpine AS builder
COPY . /app
WORKDIR /app/
RUN go build -o bin/fiboser ./cmd/main/

FROM alpine:latest
WORKDIR /app/
COPY --from=builder /app/bin/fiboser .
CMD ["./fiboser", "-redisDsn", "localhost:6379/?db=0"]