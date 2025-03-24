FROM golang:1.23

WORKDIR /songsApp

COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN go build -o app .

EXPOSE 8081

CMD ["./app"]
