FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/api ./cmd/api/*.go

EXPOSE 3000

CMD ["./bin/api"]