package main

import (
	"dune2clone/fibrandom"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var rnd fibrandom.FibRandom

func main() {
	// geometry.SetDegreesInCircleAmount(100)
	rl.InitWindow(int32(WINDOW_W), int32(WINDOW_H), "DAS IST KEIN DUNE 2!")
	rl.SetTargetFPS(RENDERER_DESIRED_FPS)
	rl.SetWindowState(rl.FlagWindowResizable)
	rl.SetExitKey(rl.KeyEscape)

	rnd.InitDefault()
	loadResources()

	//for i := 0; i <= 360; i+=10 {
	//	debugWritef("%ddeg is %d sector\n", i, degreeToRotationFrameNumber(i))
	//}
	game := game{}
	game.startGame()


	rl.CloseWindow()
}
