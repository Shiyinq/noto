FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]