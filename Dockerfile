# Use an official Golang runtime as a parent image
FROM golang:1.18-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Use vendored dependencies
RUN go mod vendor

# Build the Go app with vendored dependencies
RUN go build -mod=vendor -o /golang_app

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD [ "/golang_app" ]
