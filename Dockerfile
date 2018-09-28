FROM golang:1.11-alpine
COPY ./ /go/src/cweagans/hookreceiver
WORKDIR /go/src/cweagans/hookreceiver
RUN go build
CMD ["/go/src/cweagans/hookreceiver/hookreceiver"]
