FROM golang:1.12-stretch

WORKDIR /go/src/app
COPY . .

RUN go build ./src/app.go

CMD ./app

