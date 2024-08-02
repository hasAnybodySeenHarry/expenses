FROM golang:1.20-alpine AS builder

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /bin/app /app/app

EXPOSE 8080

CMD ["./app"]