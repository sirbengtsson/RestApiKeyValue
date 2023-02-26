FROM golang:1.19-alpine

WORKDIR /app

COPY . .
COPY go.mod .
COPY go.sum .
RUN go mod download

RUN go build -o server

CMD ["./server"]