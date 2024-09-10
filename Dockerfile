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
RUN addgroup -g 65532 nonroot && adduser -DHs /sbin/nologin -u 65532 -G nonroot nonroot
COPY --from=builder /app/main .
RUN chown nonroot:nonroot /app/main
USER nonroot:nonroot
EXPOSE 8080
ENV IS_CONTAINER=true
CMD ["./main"]
