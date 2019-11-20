FROM golang:1.13
#RUN mkdir /go/src
ARG DOCKER_NETRC

ADD . /go/src/dispatcher_rabbit_to_dwh
WORKDIR /go/src/dispatcher_rabbit_to_dwh

RUN echo "${DOCKER_NETRC}" > ~/.netrc
RUN go env GONOSUMDB="git.fin-dev.ru" && go build -o main

CMD ["./main"]