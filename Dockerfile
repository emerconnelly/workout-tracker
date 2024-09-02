# Start from a Go base image
FROM golang:1.23-alpine

WORKDIR /app

COPY backend/ .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

EXPOSE 8080

CMD ["./main"]
