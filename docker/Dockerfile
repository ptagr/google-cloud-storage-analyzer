FROM alpine

MAINTAINER Punit Agrawal <punagrawal@ebay.com>

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

ADD svc /svc

ENTRYPOINT ["/svc"]
