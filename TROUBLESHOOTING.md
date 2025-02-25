# Troubleshooting Guide

## Common Error: "imports syscall/js: build constraints exclude all Go files"

If you see this error:
```
imports syscall/js: build constraints exclude all Go files in C:\Program Files\Go\src\syscall\js
```

### What's happening?
This error occurs because you're trying to run the Go code directly with `go run main.go` on your local machine. The code imports the `syscall/js` package, which is only available when compiling Go code for WebAssembly.

### Solution:

**Don't use `go run main.go`**. This code is meant to be compiled to WebAssembly and run in a browser.

#### For Windows (PowerShell):

1. Run the included PowerShell script:
   ```powershell
   .\build.ps1
   ```

2. If you don't have the wasm_exec.js file, run:
   ```powershell
   .\get_wasm_exec.ps1
   ```

3. Serve the files with a local web server (e.g., using Python):
   ```powershell
   python -m http.server
   ```

4. Open your browser and go to: http://localhost:8000

#### For Linux/macOS:

1. Run the included bash script:
   ```bash
   chmod +x build.sh
   ./build.sh
   ```

2. If you don't have the wasm_exec.js file, run:
   ```bash
   chmod +x get_wasm_exec.sh
   ./get_wasm_exec.sh
   ```

3. Serve the files with a local web server:
   ```bash
   python -m http.server
   ```

4. Open your browser and go to: http://localhost:8000

### Explanation:

The `syscall/js` package provides functions for interacting with JavaScript from Go code when compiled to WebAssembly. It's not available for normal Go programs that run directly on your operating system.

When you set the environment variables `GOOS=js` and `GOARCH=wasm` (which the build scripts do automatically), you're telling the Go compiler to target WebAssembly for web browsers instead of your local operating system.