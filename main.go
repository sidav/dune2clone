package main

import (
	"dune2clone/fibrandom"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var rnd fibrandom.FibRandom

func main() {
	rl.InitWindow(int32(WINDOW_W), int32(WINDOW_H), "DAS IST KEIN DUNE 2!")
	rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyEscape)

	rnd.InitDefault()
	loadResources()

	//for i := 0; i <= 360; i+=10 {
	//	debugWritef("%ddeg is %d sector\n", i, degreeToRotationFrameNumber(i))
	//}

	game := game{}
	game.startGame()
	fmt.Println("Yeah, I'm working")

	rl.CloseWindow()
}
