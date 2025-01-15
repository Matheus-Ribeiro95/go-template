FROM golang:1.23-alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY main.go ./
COPY html ./html
COPY lib ./lib
RUN go build main.go

CMD ["./main"]
