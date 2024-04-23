# Use golang's official docker image
FROM golang:alpine

# Install buil deps
RUN apk add build-base libx11-dev git

# Change work directory
WORKDIR /app

# Copy tlock source code
COPY . /app

# Install go deps
RUN go mod tidy

# Build
RUN CGO_ENABLED=1 go build -o /usr/bin/tlock tlock/main.go 

# Set entry point to the binary
ENTRYPOINT ["/usr/bin/tlock"]

