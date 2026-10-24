# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /calculator-api ./cmd/api

# Stage 2: Minimal runtime image
FROM alpine:3.20

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /calculator-api /app/calculator-api

EXPOSE 8080

ENV PORT=8080

CMD ["/app/calculator-api"]
