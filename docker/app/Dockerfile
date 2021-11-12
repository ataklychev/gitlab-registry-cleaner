FROM golang:1.17

WORKDIR /go/src/app
COPY ./* .

RUN chmod 755 ./docker/app/entrypoint.sh

RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT ["./docker/app/entrypoint.sh"]