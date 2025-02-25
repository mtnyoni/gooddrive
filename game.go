package main

import (
	"math"
	"math/rand"
	"syscall/js"
	"time"
)

const (
	// Canvas dimensions
	canvasWidth  = 800
	canvasHeight = 600

	// Game physics constants
	gravity   = 9.8
	friction  = 0.95
	dragForce = 0.1

	// Road constants
	roadWidth = 600
	laneCount = 3

	// Game speed
	initialGameSpeed   = 5
	maxGameSpeed       = 15
	gameSpeedIncrement = 0.001
)

// Car represents the player's vehicle
type Car struct {
	X, Y          float64 // Position
	Width, Height float64 // Dimensions
	Speed         float64 // Current speed
	Acceleration  float64 // Current acceleration
	Steering      float64 // Steering angle
	MaxSpeed      float64 // Maximum speed
	Crashed       bool    // Crash state
}

// Obstacle represents objects on the road
type Obstacle struct {
	X, Y          float64 // Position
	Width, Height float64 // Dimensions
	Type          string  // Type of obstacle (trash, stone)
	Active        bool    // Whether obstacle is active
}

// Game holds the game state
type Game struct {
	Car       *Car
	Obstacles []*Obstacle
	Canvas    js.Value
	Context   js.Value
	GameSpeed float64
	Distance  float64
	Score     int
	GameOver  bool
	Keys      map[string]bool
	LastFrame float64
}

// NewGame creates a new game instance
func NewGame() *Game {
	// Get canvas and context
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", "game-canvas")
	context := canvas.Call("getContext", "2d")

	// Create player car
	car := &Car{
		X:            canvasWidth / 2,
		Y:            canvasHeight - 100,
		Width:        50,
		Height:       80,
		Speed:        0,
		Acceleration: 0,
		MaxSpeed:     10,
	}

	// Create initial game state
	game := &Game{
		Car:       car,
		Obstacles: make([]*Obstacle, 0),
		Canvas:    canvas,
		Context:   context,
		GameSpeed: initialGameSpeed,
		Distance:  0,
		Score:     0,
		GameOver:  false,
		Keys:      make(map[string]bool),
		LastFrame: js.Global().Get("performance").Call("now").Float(),
	}

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	return game
}

// Start begins the game loop
func (g *Game) Start() {
	if g.GameOver {
		g.Reset()
	}

	// Set up animation frame callback
	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		now := args[0].Float()
		dt := (now - g.LastFrame) / 1000 // Convert to seconds
		g.LastFrame = now

		// Update game state
		g.Update(dt)

		// Render game
		g.Render()

		// Continue animation loop if game is not over
		if !g.GameOver {
			js.Global().Call("requestAnimationFrame", renderFrame)
		}

		return nil
	})

	// Start animation loop
	js.Global().Call("requestAnimationFrame", renderFrame)
}

// Reset restarts the game
func (g *Game) Reset() {
	g.Car.X = canvasWidth / 2
	g.Car.Y = canvasHeight - 100
	g.Car.Speed = 0
	g.Car.Acceleration = 0
	g.Car.Crashed = false
	g.Obstacles = make([]*Obstacle, 0)
	g.GameSpeed = initialGameSpeed
	g.Distance = 0
	g.Score = 0
	g.GameOver = false
}

// Update updates game state based on time elapsed
func (g *Game) Update(dt float64) {
	// Skip update if game is over
	if g.GameOver {
		return
	}

	// Update car physics
	g.updateCarPhysics(dt)

	// Update obstacles
	g.updateObstacles(dt)

	// Generate new obstacles
	g.generateObstacles()

	// Check for collisions
	g.checkCollisions()

	// Update game metrics
	g.Distance += g.GameSpeed * dt
	g.Score = int(g.Distance)

	// Gradually increase game speed
	if g.GameSpeed < maxGameSpeed {
		g.GameSpeed += gameSpeedIncrement * dt
	}
}

// updateCarPhysics applies physics to car movement
func (g *Game) updateCarPhysics(dt float64) {
	car := g.Car

	// Debug values before update
	println("Before update - X:", car.X, "Speed:", car.Speed, "Steering:", car.Steering)

	// Your existing code for applying acceleration
	if g.Keys["ArrowUp"] || g.Keys["KeyW"] {
		car.Acceleration = 8.0
	} else if g.Keys["ArrowDown"] || g.Keys["KeyS"] {
		car.Acceleration = -4.0
	} else {
		car.Acceleration = 0
	}

	// Apply steering with increased effect
	if g.Keys["ArrowLeft"] || g.Keys["KeyA"] {
		car.Steering = -20.0 // Much higher value for noticeable turning
	} else if g.Keys["ArrowRight"] || g.Keys["KeyD"] {
		car.Steering = 20.0 // Much higher value
	} else {
		car.Steering = 0
	}

	// Update speed based on acceleration
	car.Speed += car.Acceleration * dt

	// Apply friction and drag
	car.Speed *= friction

	// Clamp speed to max speed
	if car.Speed > car.MaxSpeed {
		car.Speed = car.MaxSpeed
	} else if car.Speed < -car.MaxSpeed/2 {
		car.Speed = -car.MaxSpeed / 2
	}

	// Make steering work even at low speeds
	minEffectiveSpeed := 1.0
	effectiveSpeed := math.Max(math.Abs(car.Speed), minEffectiveSpeed)
	steeringEffect := car.Steering * effectiveSpeed * dt
	car.X += steeringEffect

	// Keep car on screen (your existing code)
	roadLeftEdge := float64((canvasWidth - roadWidth) / 2)
	roadRightEdge := roadLeftEdge + roadWidth

	if car.X-car.Width/2 < roadLeftEdge {
		car.X = roadLeftEdge + car.Width/2
	} else if car.X+car.Width/2 > roadRightEdge {
		car.X = roadRightEdge - car.Width/2
	}

	// Debug values after update
	println("After update - X:", car.X, "Speed:", car.Speed, "Steering:", car.Steering, "Effective steering:", steeringEffect)
}

// updateObstacles moves obstacles and removes those off-screen
func (g *Game) updateObstacles(dt float64) {
	var activeObstacles []*Obstacle

	for _, obs := range g.Obstacles {
		// Move obstacle down based on game speed
		obs.Y += g.GameSpeed * dt * 100

		// Keep active obstacles
		if obs.Y < canvasHeight+obs.Height {
			activeObstacles = append(activeObstacles, obs)
		}
	}

	g.Obstacles = activeObstacles
}

// generateObstacles creates new obstacles at random intervals
func (g *Game) generateObstacles() {
	// Chance to generate new obstacle increases with game speed
	if rand.Float64() > 0.97-(g.GameSpeed/100) {
		obstacleTypes := []string{"trash", "stone"}
		obstacleType := obstacleTypes[rand.Intn(len(obstacleTypes))]

		// Determine size based on type
		var width, height float64
		if obstacleType == "trash" {
			width = 30
			height = 30
		} else { // stone
			width = 40
			height = 35
		}

		// Determine lane position
		roadLeftEdge := (canvasWidth - roadWidth) / 2
		laneWidth := roadWidth / laneCount
		lane := rand.Intn(laneCount)
		x := float64(roadLeftEdge) + (float64(lane) * float64(laneWidth)) + float64(laneWidth/2)

		// Create obstacle
		obstacle := &Obstacle{
			X:      x,
			Y:      -height,
			Width:  width,
			Height: height,
			Type:   obstacleType,
			Active: true,
		}

		g.Obstacles = append(g.Obstacles, obstacle)
	}
}

// checkCollisions detects and handles collisions
func (g *Game) checkCollisions() {
	car := g.Car

	// Skip if already crashed
	if car.Crashed {
		g.GameOver = true
		return
	}

	// Check for collisions with obstacles
	for _, obs := range g.Obstacles {
		if obs.Active {
			// Simple rectangle collision
			if car.X+car.Width/2 > obs.X-obs.Width/2 &&
				car.X-car.Width/2 < obs.X+obs.Width/2 &&
				car.Y+car.Height/2 > obs.Y-obs.Height/2 &&
				car.Y-car.Height/2 < obs.Y+obs.Height/2 {
				// Handle collision based on obstacle type
				car.Crashed = true
				break
			}
		}
	}
}

// Render draws the game state to the canvas
func (g *Game) Render() {
	ctx := g.Context

	// Clear canvas
	ctx.Set("fillStyle", "#87CEEB") // Sky blue background
	ctx.Call("fillRect", 0, 0, canvasWidth, canvasHeight)

	// Draw road
	g.drawRoad()

	// Draw obstacles
	g.drawObstacles()

	// Draw car
	g.drawCar()

	// Draw score
	g.drawScore()

	// Draw game over screen if needed
	if g.GameOver {
		g.drawGameOver()
	}
}

// drawRoad renders the road
func (g *Game) drawRoad() {
	ctx := g.Context
	roadLeftEdge := (canvasWidth - roadWidth) / 2

	// Draw road background
	ctx.Set("fillStyle", "#333333") // Dark gray
	ctx.Call("fillRect", roadLeftEdge, 0, roadWidth, canvasHeight)

	// Draw lane markers
	ctx.Set("strokeStyle", "#FFFFFF")
	ctx.Set("lineWidth", 3)

	laneWidth := roadWidth / laneCount

	// Dash pattern for lane markers
	dashLength := 30.0
	gapLength := 20.0

	// Calculate offset for moving dashes
	dashOffset := math.Mod(g.Distance*10, dashLength+gapLength)

	// Draw lane markings
	for i := 1; i < laneCount; i++ {
		x := float64(roadLeftEdge) + float64(i)*float64(laneWidth)

		ctx.Call("beginPath")
		ctx.Call("setLineDash", []interface{}{dashLength, gapLength})
		ctx.Set("lineDashOffset", -dashOffset)
		ctx.Call("moveTo", x, 0)
		ctx.Call("lineTo", x, canvasHeight)
		ctx.Call("stroke")
	}

	// Reset line dash
	ctx.Call("setLineDash", []interface{}{})

	// Draw road edges
	ctx.Set("strokeStyle", "#FFFFFF")
	ctx.Set("lineWidth", 5)

	ctx.Call("beginPath")
	ctx.Call("moveTo", roadLeftEdge, 0)
	ctx.Call("lineTo", roadLeftEdge, canvasHeight)
	ctx.Call("stroke")

	ctx.Call("beginPath")
	ctx.Call("moveTo", roadLeftEdge+roadWidth, 0)
	ctx.Call("lineTo", roadLeftEdge+roadWidth, canvasHeight)
	ctx.Call("stroke")
}

// drawObstacles renders the obstacles
func (g *Game) drawObstacles() {
	ctx := g.Context

	for _, obs := range g.Obstacles {
		ctx.Call("save")

		if obs.Type == "trash" {
			// Draw trash (paper/garbage)
			ctx.Set("fillStyle", "#DDDDDD") // Light gray
			ctx.Call("beginPath")
			ctx.Call("ellipse",
				obs.X,
				obs.Y,
				obs.Width/2,
				obs.Height/2,
				0, 0, 2*math.Pi)
			ctx.Call("fill")

			// Add some details
			ctx.Set("strokeStyle", "#999999")
			ctx.Set("lineWidth", 2)
			ctx.Call("beginPath")
			ctx.Call("moveTo", obs.X-obs.Width/4, obs.Y)
			ctx.Call("lineTo", obs.X+obs.Width/4, obs.Y)
			ctx.Call("stroke")
		} else { // stone
			// Draw stone
			ctx.Set("fillStyle", "#777777") // Stone gray
			ctx.Call("beginPath")
			ctx.Call("ellipse",
				obs.X,
				obs.Y,
				obs.Width/2,
				obs.Height/2,
				0, 0, 2*math.Pi)
			ctx.Call("fill")

			// Add some highlights
			ctx.Set("fillStyle", "#999999")
			ctx.Call("beginPath")
			ctx.Call("ellipse",
				obs.X-obs.Width/4,
				obs.Y-obs.Height/4,
				obs.Width/6,
				obs.Height/6,
				0, 0, 2*math.Pi)
			ctx.Call("fill")
		}

		ctx.Call("restore")
	}
}

// drawCar renders the player's car
func (g *Game) drawCar() {
	ctx := g.Context
	car := g.Car

	ctx.Call("save")

	// Draw car body
	if car.Crashed {
		ctx.Set("fillStyle", "#AA0000") // Red for crashed car
	} else {
		ctx.Set("fillStyle", "#0000AA") // Blue for normal car
	}

	// Car body
	ctx.Call("beginPath")
	ctx.Call("roundRect",
		car.X-car.Width/2,
		car.Y-car.Height/2,
		car.Width,
		car.Height,
		10)
	ctx.Call("fill")

	// Windshield
	ctx.Set("fillStyle", "#CCCCFF")
	ctx.Call("beginPath")
	ctx.Call("roundRect",
		car.X-car.Width/3,
		car.Y-car.Height/3,
		car.Width*2/3,
		car.Height/4,
		5)
	ctx.Call("fill")

	// Wheels
	ctx.Set("fillStyle", "#000000")

	// Front left wheel
	ctx.Call("beginPath")
	ctx.Call("ellipse",
		car.X-car.Width/2-5,
		car.Y-car.Height/4,
		7,
		12,
		0, 0, 2*math.Pi)
	ctx.Call("fill")

	// Front right wheel
	ctx.Call("beginPath")
	ctx.Call("ellipse",
		car.X+car.Width/2+5,
		car.Y-car.Height/4,
		7,
		12,
		0, 0, 2*math.Pi)
	ctx.Call("fill")

	// Rear left wheel
	ctx.Call("beginPath")
	ctx.Call("ellipse",
		car.X-car.Width/2-5,
		car.Y+car.Height/4,
		7,
		12,
		0, 0, 2*math.Pi)
	ctx.Call("fill")

	// Rear right wheel
	ctx.Call("beginPath")
	ctx.Call("ellipse",
		car.X+car.Width/2+5,
		car.Y+car.Height/4,
		7,
		12,
		0, 0, 2*math.Pi)
	ctx.Call("fill")

	ctx.Call("restore")

	// If crashed, add some effects
	if car.Crashed {
		// Smoke/dust particles
		ctx.Set("fillStyle", "rgba(200,200,200,0.5)")

		for i := 0; i < 10; i++ {
			size := rand.Float64()*20 + 10
			offsetX := (rand.Float64() - 0.5) * car.Width * 1.5
			offsetY := (rand.Float64() - 0.5) * car.Height * 1.5

			ctx.Call("beginPath")
			ctx.Call("ellipse",
				car.X+offsetX,
				car.Y+offsetY,
				size/2,
				size/2,
				0, 0, 2*math.Pi)
			ctx.Call("fill")
		}
	}
}

// drawScore displays the current score
func (g *Game) drawScore() {
	ctx := g.Context

	ctx.Set("font", "24px Arial")
	ctx.Set("fillStyle", "#000000")
	ctx.Call("fillText", "Score: "+js.ValueOf(g.Score).String(), 20, 30)
	ctx.Call("fillText", "Speed: "+js.ValueOf(int(g.GameSpeed*10)/10).String()+"x", 20, 60)
}

// drawGameOver displays game over message
func (g *Game) drawGameOver() {
	ctx := g.Context

	// Semi-transparent overlay
	ctx.Set("fillStyle", "rgba(0,0,0,0.7)")
	ctx.Call("fillRect", 0, 0, canvasWidth, canvasHeight)

	// Game over text
	ctx.Set("font", "48px Arial")
	ctx.Set("fillStyle", "#FFFFFF")
	ctx.Set("textAlign", "center")
	ctx.Call("fillText", "GAME OVER", canvasWidth/2, canvasHeight/2-50)

	// Score text
	ctx.Set("font", "32px Arial")
	ctx.Call("fillText", "Final Score: "+js.ValueOf(g.Score).String(), canvasWidth/2, canvasHeight/2+20)

	// Restart instruction
	ctx.Set("font", "24px Arial")
	ctx.Call("fillText", "Press SPACE to restart", canvasWidth/2, canvasHeight/2+80)

	// Reset text alignment
	ctx.Set("textAlign", "left")
}

// UpdateInput processes player input
func (g *Game) UpdateInput(keyCode string, keyDown bool) {
	// Update key state
	g.Keys[keyCode] = keyDown

	// Debug logging to console
	println("Go received key:", keyCode, keyDown)

	// Check for restart
	if keyCode == "Space" && keyDown && g.GameOver {
		g.Reset()
		g.Start()
	}
}
