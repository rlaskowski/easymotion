FROM rlaskowski/opencv:4.6.0

ENV GOPATH /go

ARG buildname

WORKDIR /go/src/gocv.io/x/gocv

COPY . .

RUN make build buildname=${buildname}

ENTRYPOINT ["dist/${buildname}/${buildname}", "run"]
