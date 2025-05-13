FROM golang:1.24.1-alpine AS builder

COPY . /github.com/nogavadu/articles-service
WORKDIR /github.com/nogavadu/articles-service

RUN go mod download
RUN go build -o ./bin/articles-service cmd/http_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/nogavadu/articles-service/bin/articles-service .

CMD ["./articles-service"]