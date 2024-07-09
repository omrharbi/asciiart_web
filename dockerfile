# this to wath can use the os
FROM golang:alpine

WORKDIR /web-main
COPY  . .

RUN apk add  bash

RUN go build -o server .

EXPOSE 80

ENTRYPOINT ["./server"]