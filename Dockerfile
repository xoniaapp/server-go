
FROM golang:alpine

RUN go build . -o server

RUN chmod +x ./server

EXPOSE 5000

CMD ["./server"]
