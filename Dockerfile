FROM golang:1.13.8-alpine

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go build main.go -o main

EXPOSE 8080

CMD ["main"]