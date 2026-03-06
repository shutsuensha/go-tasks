# Stage 1: build
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app cmd/app/main.go

# Stage 2: runtime
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]