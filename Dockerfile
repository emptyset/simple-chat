FROM golang:1.10.1 as builder
WORKDIR /go/src/github.com/emptyset/simple-chat
ADD . /go/src/github.com/emptyset/simple-chat
RUN CGO_ENABLED=0 GOOS=linux go build -a -o ./bin/app ./cmd/server/main.go

FROM alpine:latest as release
WORKDIR /root/
COPY --from=builder /go/src/github.com/emptyset/simple-chat/bin/app .
CMD ["./app"]
