#!/bin/bash
# Shell script to copy wasm_exec.js from the Go installation

# Get the Go root directory
GOROOT=$(go env GOROOT)

# Source path for wasm_exec.js
SOURCE_FILE="${GOROOT}/misc/wasm/wasm_exec.js"

# Copy the file to the current directory
cp "$SOURCE_FILE" "./wasm_exec.js"

echo "Copied wasm_exec.js from $SOURCE_FILE to the current directory"