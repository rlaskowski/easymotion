FROM rlaskowski/opencv:4.6.0

ENV GOPATH /go

WORKDIR /go/src/gocv.io/x/gocv

COPY . .

RUN make build

VOLUME [ "/videos" ]

ENTRYPOINT ["dist/easymotion", "run"]
