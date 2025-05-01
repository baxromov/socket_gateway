# Use the official Go image as the base image
FROM golang:1.24

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the app source code
COPY . .

# Build the application binary
RUN go build -o app

# Expose the port (use the same port as your application)
EXPOSE 5000

# Command to run the app
CMD ["./app", "-port", "5000"]