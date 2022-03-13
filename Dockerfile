FROM gocv/opencv:4.5.4

ENV GOPATH /go

WORKDIR /go/src/gocv.io/x/gocv

COPY . .

RUN go build -o dist/easymotion cmd/easymotion/main.go

CMD ["dist/easymotion", "run"]
