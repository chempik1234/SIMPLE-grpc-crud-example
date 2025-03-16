FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go test -c -o integration-tests ./tests/integrationTests

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/integration-tests /app/integration-tests

CMD ["./integration-tests"]