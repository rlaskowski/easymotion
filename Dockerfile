FROM gocv/opencv:4.5.4

ENV GOPATH /go

WORKDIR /go/src/gocv.io/x/gocv

COPY . .

RUN make build-linux

CMD ["dist/easymotion", "run"]
