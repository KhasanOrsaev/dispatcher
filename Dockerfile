FROM golang:1.13
#RUN mkdir /go/src
COPY . /go/src/dispatcher_rabbit_to_dwh
WORKDIR /go/src/dispatcher_rabbit_to_dwh

RUN go env GONOSUMDB="git.fin-dev.ru" && go build -o main
CMD ["./main"]