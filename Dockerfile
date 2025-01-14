FROM golang:1.23-alpine AS builder
WORKDIR /app/build

COPY cmd ./
RUN go mod init main && \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o client .

# Start a new stage from scratch
FROM alpine:latest
WORKDIR /app/am2pb

COPY --from=builder /app/build/client .
EXPOSE 5000
USER 1001

ENTRYPOINT ["./client"]
