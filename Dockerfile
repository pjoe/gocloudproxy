FROM alpine:edge

ENV GOPATH=/go PATH=/go/bin:$PATH

RUN apk add --no-cache ca-certificates \
    && apk --no-cache add --virtual build-dependencies musl-dev go git \
    && go get -u github.com/pjoe/gocloudproxy \
    && go clean --cache \
    && apk del --purge build-dependencies \
    && rm -rf /go/pkg /go/src

CMD /go/bin/gocloudproxy