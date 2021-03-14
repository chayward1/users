FROM golang:1.13.8-alpine

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go build -o main main.go

EXPOSE 8080

CMD ["main"]