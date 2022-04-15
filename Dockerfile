FROM golang:latest

ENV GO111MODULE=on

WORKDIR /app/server
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go build

EXPOSE 8080

CMD ["./m"]