FROM rlaskowski/opencv:4.5.4

ENV GOPATH /go

WORKDIR /go/src/gocv.io/x/gocv

COPY . .

RUN make build-raspi

VOLUME [ "/videos" ]

CMD ["dist/easymotion", "run", "-f", "/videos"]
