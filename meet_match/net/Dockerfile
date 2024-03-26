# Use the official Golang image as the base image
FROM golang:1.18

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY backend/go.mod ./
COPY backend/go.sum ./
RUN go mod download

# Copy the server code to the container
COPY . .

# Build the Go application
RUN cd backend
RUN go build -o main .
RUN cd ..

# Expose port 8080 to the outside world
#EXPOSE 8080

# Command to run the executable
CMD ["./main"]