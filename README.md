# Infinite Driving Game

An infinite driving game built with Go and WebAssembly, featuring car physics, obstacle avoidance, and collision detection.

## Features

- Play directly in the browser using WebAssembly
- Realistic car physics with acceleration, friction, and steering
- Obstacles (trash and stones) to avoid
- Distance-based scoring system
- Visual effects for crashes
- Keyboard controls

## How to Build

To compile the Go code to WebAssembly, follow these steps:

1. Make sure you have Go 1.13 or later installed.

2. Build the WebAssembly binary:
   ```bash
   GOOS=js GOARCH=wasm go build -o main.wasm
   ```

3. Make sure you have the wasm_exec.js file (included in this repo), or copy it from your Go installation:
   ```bash
   cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
   ```

## How to Run

You need a web server to serve the files because browsers won't load WebAssembly from local file:// URLs.

You can use any web server. Here's a simple way with Python:

```bash
# Python 3
python -m http.server

# Python 2
python -m SimpleHTTPServer
```

Then open your browser and navigate to: http://localhost:8000

## How to Play

- **Arrow Keys or WASD:** Control the car
- **Up Arrow / W:** Accelerate
- **Down Arrow / S:** Brake/Reverse
- **Left Arrow / A:** Steer left
- **Right Arrow / D:** Steer right
- **Space:** Restart after game over

The goal is to survive as long as possible by avoiding obstacles. Your score increases based on the distance traveled.

## Game Design

- The game uses a simple physics model including acceleration, friction, and steering.
- As you progress, the game speed gradually increases, making it more challenging.
- Colliding with obstacles (trash or stones) causes the car to crash, ending the game.

## Implementation Details

The game is implemented in Go and compiled to WebAssembly so it can run in a browser:

- `main.go` - Entry point, handling WebAssembly initialization and JavaScript interop
- `game.go` - Game logic, physics, rendering, and state management
- `index.html` - HTML interface and JavaScript for loading the WASM module
- `wasm_exec.js` - Go's WebAssembly runtime support

The game uses the HTML5 Canvas API for rendering, accessed via Go's syscall/js package.