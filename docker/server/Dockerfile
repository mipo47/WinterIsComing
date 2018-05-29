FROM golang:1.10.2-alpine AS builder

WORKDIR /go/src/github.com/mipo47/WinterIsComing
ADD . .
RUN go run run.go build server



FROM alpine:3.6

COPY --from=builder /go/src/github.com/mipo47/WinterIsComing/out/server /usr/bin/winter_server
ENTRYPOINT ["/usr/bin/winter_server"]
