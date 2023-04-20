FROM golang:1.20.1-bullseye

WORKDIR /app

COPY . ./app

RUN go mod download 

RUN go build -o main 

CMD ["./app/main"]