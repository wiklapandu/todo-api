FROM golang:1.23-alpine

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod .
COPY go.sum .
COPY ./vendor/ ./vendor

RUN go build -mod=vendor -o main .

COPY . .

CMD ["air"]