#!/bin/bash
# Build script for compiling Go to WebAssembly

# Set environment variables for WebAssembly compilation
export GOOS=js
export GOARCH=wasm

# Compile the Go code to WebAssembly
go build -o main.wasm

echo "WebAssembly build completed. Output file: main.wasm"
echo "To run the application, serve the directory with a HTTP server and open index.html"