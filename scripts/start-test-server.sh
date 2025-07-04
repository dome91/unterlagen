#!/bin/bash

# Set the required environment variables
export UNTERLAGEN_SERVER_SESSION_KEY=${UNTERLAGEN_SERVER_SESSION_KEY:-"test-session-key"}
export UNTERLAGEN_SERVER_PORT=${UNTERLAGEN_SERVER_PORT:-8080}

# Build the application
echo "Building Unterlagen application..."
go build -o ./tmp/unterlagen-test cmd/unterlagen.go

# Start the server in test mode
echo "Starting Unterlagen server on port $UNTERLAGEN_SERVER_PORT..."
./tmp/unterlagen-test

# This script will keep running until the server is terminated
