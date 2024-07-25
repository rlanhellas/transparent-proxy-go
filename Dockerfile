FROM golang:1.22.5-alpine3.20

WORKDIR app

ADD main.go .
ADD go.mod .
RUN GOOS=linux CGO_ENABLED=0 go build -o main

EXPOSE 8080

CMD ["./main"]
