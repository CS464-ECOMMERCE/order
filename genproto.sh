#!/bin/bash

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check if protoc exists
if ! command_exists protoc; then
    echo "protoc is not installed. Please install Protocol Buffers compiler."
    exit 1
fi

# Check if protoc-gen-go exists
if ! command_exists protoc-gen-go; then
    echo "protoc-gen-go is not installed. Installing..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Generate Go code from proto files
echo "Generating Go code from proto files..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/*.proto

echo "Done." 