# Build stage: resolve deps and compile inside the container (no host go/mod needed)
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Copy full source; deps resolution and build happen only inside the container
COPY . .
RUN go mod tidy && CGO_ENABLED=0 go build -o /server ./cmd/server

# Runtime stage
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
