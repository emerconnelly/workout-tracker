# Build stage
FROM golang:1.23-alpine3.20 AS builder
WORKDIR /app
COPY backend/ .
RUN go mod download
RUN apk add --no-cache upx
RUN CGO_ENABLED=0 GOOS=linux go build -o main .
RUN upx --best --lzma main

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
