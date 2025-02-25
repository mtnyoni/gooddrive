# PowerShell script for compiling Go to WebAssembly

# Set environment variables for WebAssembly compilation
$env:GOOS = "js"
$env:GOARCH = "wasm"

# Compile the Go code to WebAssembly
go build -o main.wasm

Write-Host "WebAssembly build completed. Output file: main.wasm"
Write-Host "To run the application, serve the directory with a HTTP server and open index.html"