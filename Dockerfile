# Start from a Go base image
FROM golang:1.23

# Set the working directory in the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
