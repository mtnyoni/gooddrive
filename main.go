package main

import (
	"syscall/js"
)

func main() {
	c := make(chan struct{})
	
	// Initialize the game
	game := NewGame()
	
	// Register JS functions
	js.Global().Set("startGame", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		game.Start()
		return nil
	}))
	
	js.Global().Set("updateInput", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			keyCode := args[0].String()
			keyDown := args[1].Bool()
			game.UpdateInput(keyCode, keyDown)
		}
		return nil
	}))

	println("WASM Go Initialized")
	<-c
}