FROM golang:1.20-alpine

WORKDIR /app

# Install git for private repos
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application for Linux
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o bin/route-master ./cmd

# Expose port
EXPOSE 8081

# Run the application
CMD ["/app/bin/route-master"]
