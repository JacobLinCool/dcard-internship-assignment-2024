FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build .

FROM alpine AS production

WORKDIR /app

COPY --from=builder /app/ad-service .

CMD ["./ad-service"]
