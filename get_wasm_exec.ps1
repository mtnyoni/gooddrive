# PowerShell script to copy wasm_exec.js from the Go installation

# Get the Go root directory
$goroot = & go env GOROOT

# Source path for wasm_exec.js
$sourceFile = Join-Path -Path $goroot -ChildPath "misc\wasm\wasm_exec.js"

# Copy the file to the current directory
Copy-Item -Path $sourceFile -Destination ".\wasm_exec.js"

Write-Host "Copied wasm_exec.js from $sourceFile to the current directory"
