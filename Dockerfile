FROM golang:1.19-alpine

WORKDIR /app

RUN apk add git bash curl

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . . 

RUN go build -o entrypoint .

ENTRYPOINT ["/app/entrypoint"]