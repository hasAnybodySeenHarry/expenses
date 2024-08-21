FROM golang:1.22-alpine AS builder

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /bin/app ./cmd/api

FROM scratch

WORKDIR /app

COPY --from=builder /bin/app /app/app

EXPOSE 8080
EXPOSE 50051

CMD ["./app"]