# build
FROM golang:1-alpine as builder

RUN rm -rf /var/cache/apk/* && rm -rf /tmp/*
RUN apk update
RUN apk --no-cache add -U make git

WORKDIR /go/src/github.com/denouche/plex-watcher
COPY . /go/src/github.com/denouche/plex-watcher
RUN make deps build

# run
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/denouche/plex-watcher/plex-watcher .
CMD ["/plex-watcher"]

