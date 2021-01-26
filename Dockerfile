FROM golang:latest

RUN mkdir /app
WORKDIR /app

ENV GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download && go get github.com/pilu/fresh
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/app
