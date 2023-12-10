FROM golang:1.16

WORKDIR /app

COPY . .

RUN go build -o main .

EXPOSE 8000

CMD ["./main"]
