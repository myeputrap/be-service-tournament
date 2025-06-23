FROM golang:1.24 AS builder
WORKDIR /app

# Cache modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the full app source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./app

# Final minimal image
FROM alpine:latest
WORKDIR /app

# Copy the built binary and config file
COPY --from=builder /app/app .
COPY --from=builder /app/config.yaml .

EXPOSE 1360

CMD ["./app"]
