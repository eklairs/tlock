# Use golang's official docker image
FROM golang:alpine

# Enable colors
ENV TERM=xterm-256color

# Install buil deps
RUN apk add build-base libx11-dev git

# Change work directory
WORKDIR /app

# Copy tlock source code
COPY . /app

# Install go deps
RUN go mod tidy

# Build
RUN go build -ldflags "-X github.com/eklairs/tlock/tlock-internal/constants.VERSION=docker -w -s" -o /usr/bin/tlock tlock/main.go

# Set entry point to the binary
ENTRYPOINT ["/usr/bin/tlock"]

