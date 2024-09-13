FROM golang:1.23-alpine

WORKDIR /app

ENV CGO_ENABLED=0

RUN go install github.com/cosmtrek/air@v1.44.0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air"]
