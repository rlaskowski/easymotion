FROM golang:1.17-alpine3.15

WORKDIR /go/src/easymotion

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

CMD [ "./easymotion" ]