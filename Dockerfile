# Use a base Golang image
FROM golang:1.20-alpine

# Set the working directory
WORKDIR /app

# Copy the Go source code
COPY . .

# Build the application
RUN go build -o main .

# Expose the port your application listens on (e.g., 8080)
EXPOSE 8080

# Define the command to run when the container starts
CMD ["/app/main"]