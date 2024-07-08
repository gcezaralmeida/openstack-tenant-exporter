# Use the official Golang image with Go 1.22.5 to build the application
FROM golang:1.22.5 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY exporter.go ./

# Build the application
RUN go build -o openstack-tenant-exporter .

# Use an Ubuntu image to run the application
FROM ubuntu:22.04

# Install necessary dependencies (if any)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/openstack-tenant-exporter /app/openstack-tenant-exporter

# Expose the port the app runs on
EXPOSE 9183

# Command to run the application
CMD ["./openstack-tenant-exporter"]
