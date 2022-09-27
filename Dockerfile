FROM rlaskowski/opencv:4.5.4

ENV GOPATH /go

WORKDIR /go/src/gocv.io/x/gocv

COPY . .

RUN make build-linux

VOLUME [ "/videos" ]

CMD ["dist/easymotion", "-f", "/videos", "run"]
