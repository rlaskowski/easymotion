FROM golang:1.17.8-alpine3.15

WORKDIR /usr/local/go/easymotion

COPY dist/easymotion .

CMD [ "/usr/local/go/easymotion/easymotion", "run" ]