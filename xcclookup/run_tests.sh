#!/bin/bash

# Test runner script for xcclookup package
set -e

echo "Running unit tests for xcclookup package..."

# Run tests for each package
echo "Testing services package..."
go test -v ./pkgs/services/...

echo "Testing handlers package..."
go test -v ./pkgs/handlers/...

echo "Testing claims package..."
go test -v ./pkgs/claims/...

echo "Testing mocks package..."
go test -v ./pkgs/mocks/...

echo "Running integration tests..."
go test -v ./pkgs/integration/...

echo "Running tests with coverage..."
go test -v -coverprofile=coverage.out ./pkgs/...

echo "Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html

echo "All tests completed successfully!"
echo "Coverage report generated: coverage.html" 