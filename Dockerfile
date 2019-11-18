FROM golang:1.13
#RUN mkdir /go/src

ADD . /go/src/dispatcher_rabbit_to_dwh
WORKDIR /go/src/dispatcher_rabbit_to_dwh

RUN go build -o main
CMD ["./main"]