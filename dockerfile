FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o jobqueuer main.go

# ---- Runtime ----
FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/jobqueuer .
COPY app.env .
COPY ./app/migrations ./app/migrations

EXPOSE 8080

CMD ["./jobqueuer", "runjob"]
