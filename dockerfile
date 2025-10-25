# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o jobqueuer .

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/jobqueuer .
EXPOSE 8080
CMD ["./jobqueuer"]
