FROM golang:alpine3.13

WORKDIR /service-download

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["download-service"]