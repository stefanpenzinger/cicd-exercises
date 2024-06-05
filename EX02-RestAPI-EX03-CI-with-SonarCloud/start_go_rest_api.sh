#!/bin/bash

# Check if the file "cicd-exercises" exists
if [ -f "cicd-exercises" ]; then
    # If the file exists, delete it
    echo "Deleting existing 'cicd-exercises' file..."
    rm cicd-exercises
fi

# Build the Go application
echo "Building the Go application..."
go build

# Run the built application
echo "Running 'cicd-exercises'..."
./cicd-exercises