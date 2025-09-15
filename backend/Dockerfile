# Use the official Golang image for development
FROM golang:1.22-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Install git, which is required for `go install` to fetch modules
RUN apk add --no-cache git

# Download and install Air for live-reloading
RUN go install github.com/cosmtrek/air@v1.51.0

# Copy go module and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the application using Air
CMD ["air"]
