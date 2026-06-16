# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o cli-auth ./cmd

# Runtime stage
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache libc6-compat

COPY --from=builder /app/cli-auth .

RUN mkdir -p data
RUN mkdir -p qr-codes

CMD ["./cli-auth"]