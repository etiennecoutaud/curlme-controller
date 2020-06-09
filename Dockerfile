FROM golang:1.14.2-alpine3.11 AS builder

WORKDIR /go/src/github.com/etiennecoutaud/curlme-controller
COPY . .
RUN go build cmd/main.go

FROM alpine
COPY --from=builder /go/src/github.com/etiennecoutaud/curlme-controller/main /

ENTRYPOINT ["/main"]
