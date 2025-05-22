# Build binary
FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN go build -o deploy-action main.go

# Runtime image
FROM alpine:latest
RUN apk add --no-cache openssh-client
COPY --from=builder /app/deploy-action /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/deploy-action"]