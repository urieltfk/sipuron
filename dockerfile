# Start by building the application
FROM golang:1.19 AS build

# Set the Current Working Directory inside the container
WORKDIR /.

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Start a new stage from scratch
FROM golang:1.19

# Copy the Pre-built binary file from the previous stage
COPY --from=build /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
