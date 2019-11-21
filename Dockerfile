FROM golang:1.13
#RUN mkdir /go/src
ARG GIT_LOGIN

ADD . /go/src/dispatcher_rabbit_to_dwh
WORKDIR /go/src/dispatcher_rabbit_to_dwh


RUN go env GONOSUMDB="git.fin-dev.ru" && echo "${GIT_LOGIN}" > ~/.netrc && go build -o main

CMD ["./main"]