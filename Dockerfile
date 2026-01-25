# Build Stage
FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -o vigil-server ./cmd/server

# Final Stage
FROM alpine:latest
WORKDIR /app

# Copy the binary
COPY --from=builder /app/vigil-server .

# Copy the web folder (CRITICAL for the new UI)
COPY --from=builder /app/web ./web

EXPOSE 8090
CMD ["./vigil-server"]