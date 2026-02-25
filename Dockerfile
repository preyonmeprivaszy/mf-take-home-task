FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
# Download dependencies and build the binary named 'app'
RUN go mod download
RUN go build -o app ./cmd/main.go

# Run stage (keeps the image small)
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app .
# We will copy migrations later, but keep this in mind
COPY --from=builder /app/migrations ./migrations

ENTRYPOINT ["./app"]
