# Use the official Golang image
FROM golang:1.25-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Install git, which is required for go modules
RUN apk add --no-cache git

# Copy go module and sum files to leverage Docker cache
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Tidy the modules to create a perfect go.sum
RUN go mod tidy

# Build the binary to ensure everything is correct
RUN CGO_ENABLED=0 go build -o ./tmp/server cmd/server.go

# Install Air for live-reloading
RUN go install github.com/cosmtrek/air@v1.51.0

# Expose port 8080
EXPOSE 8080

# Command to run the application using Air.
# Air will now use the binary we built during the image creation.
CMD ["air"]