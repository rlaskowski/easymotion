FROM golang:1.17-alpine3.15

RUN apk add --update make

WORKDIR /go/src/easymotion

COPY . .

RUN go get -d -v ./... 

RUN go install -v ./... 

RUN make build

CMD [ "dist/easymotion" ]