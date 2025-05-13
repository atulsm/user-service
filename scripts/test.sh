#!/bin/bash

# Set default values
COVERAGE=0
VERBOSE=0
RACE=0
SHORT=0
PACKAGE=""

# Print usage information
usage() {
    echo "Usage: $0 [-c] [-v] [-r] [-s] [-p package]"
    echo "Options:"
    echo "  -c    Generate coverage report"
    echo "  -v    Run tests in verbose mode"
    echo "  -r    Run tests with race detector"
    echo "  -s    Run tests in short mode"
    echo "  -p    Run tests for specific package (e.g., ./internal/handlers)"
    echo "  -h    Show this help message"
    exit 1
}

# Parse command line options
while getopts "cvrsp:h" opt; do
    case $opt in
        c) COVERAGE=1 ;;
        v) VERBOSE=1 ;;
        r) RACE=1 ;;
        s) SHORT=1 ;;
        p) PACKAGE=$OPTARG ;;
        h) usage ;;
        ?) usage ;;
    esac
done

# Build the test command
CMD="go test"

# Add options based on flags
if [ $VERBOSE -eq 1 ]; then
    CMD="$CMD -v"
fi

if [ $RACE -eq 1 ]; then
    CMD="$CMD -race"
fi

if [ $SHORT -eq 1 ]; then
    CMD="$CMD -short"
fi

if [ $COVERAGE -eq 1 ]; then
    # Create coverage directory if it doesn't exist
    mkdir -p coverage
    CMD="$CMD -coverprofile=coverage/coverage.out"
fi

# Add package or default to all packages
if [ -n "$PACKAGE" ]; then
    CMD="$CMD $PACKAGE"
else
    CMD="$CMD ./..."
fi

# Print the command being run
echo "Running: $CMD"

# Run the tests
eval $CMD

# Generate coverage report if requested
if [ $COVERAGE -eq 1 ]; then
    echo "Generating coverage report..."
    go tool cover -html=coverage/coverage.out -o coverage/coverage.html
    echo "Coverage report generated at coverage/coverage.html"
    
    # Calculate and display coverage percentage
    COVERAGE_PCT=$(go tool cover -func=coverage/coverage.out | grep total | awk '{print $3}')
    echo "Total coverage: $COVERAGE_PCT"
fi 