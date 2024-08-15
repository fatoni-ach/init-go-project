FROM golang:1.22-alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod download && go mod tidy

RUN go build cmd/server/main.go 

CMD ["/app/main"]

EXPOSE 8080