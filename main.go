package main

import (
	"dune2clone/fibrandom"
	"dune2clone/map_generator"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var rnd fibrandom.FibRandom

func main() {
	// geometry.SetDegreesInCircleAmount(100)
	rl.InitWindow(int32(WINDOW_W), int32(WINDOW_H), "DAS IST KEIN DUNE 2!")
	rl.SetTargetFPS(DESIRED_FPS)
	rl.SetWindowState(rl.FlagWindowResizable)
	rl.SetExitKey(rl.KeyEscape)

	rnd.InitDefault()

	// TODO: remove
	map_generator.SetRandom(&rnd)
	generatedMap := &map_generator.GameMap{}
	generatedMap.Init(64, 64)
	generatedMap.Generate()
	for {
		rl.BeginDrawing()
		rl.EndDrawing()
		drawGeneratedMap(generatedMap)
		if rl.IsKeyDown(rl.KeySpace) || rl.IsKeyDown(rl.KeyEscape) {
			break
		} else if rl.GetKeyPressed() != 0 {
			generatedMap.Init(64, 64)
			generatedMap.Generate()
		}
	}

	loadResources()

	//for i := 0; i <= 360; i+=10 {
	//	debugWritef("%ddeg is %d sector\n", i, degreeToRotationFrameNumber(i))
	//}
	game := game{}
	game.startGame()


	rl.CloseWindow()
}
