FROM golang:alpine

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...

RUN go build -o app

EXPOSE 5000

CMD ["./app"]
