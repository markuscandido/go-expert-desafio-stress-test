FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod .
COPY main.go .

RUN go build -o stress-test main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/stress-test .

ENTRYPOINT ["./stress-test"]
