package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.InitWindow(int32(WINDOW_W), int32(WINDOW_H), "DAS IST KEIN DUNE 2!")
	rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyF12)

	loadResources()

	game := game{}
	game.startGame()
	fmt.Println("Yeah, I'm working")

	//for !rl.WindowShouldClose() {
	//
	//}

	rl.CloseWindow()
}
