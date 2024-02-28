FROM golang:latest

WORKDIR /hezzl
COPY . .

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go build -o main ./cmd/hezzl

CMD ["/hezzl/main"]

EXPOSE $PORT