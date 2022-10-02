package main

import (
	"dune2clone/map_generator"
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
)

func (g *game) selectMapToGenerateBattlefield() {
	map_generator.SetRandom(&rnd)
	generatedMap := &map_generator.GeneratedMap{}
	currSelectedPatternIndex := 0
	w := map_generator.GetPatternByIndex(currSelectedPatternIndex).MinWidth
	h := map_generator.GetPatternByIndex(currSelectedPatternIndex).MinHeight
	reGenerate := true
	for {
		if reGenerate {
			drawLoadingScreen("Generating...")
			generatedMap.Generate(w, h, currSelectedPatternIndex)
			time.Sleep(100 * time.Millisecond)
		}
		reGenerate = true
		rl.BeginDrawing()
		drawGeneratedMap(generatedMap, currSelectedPatternIndex)
		rl.EndDrawing()
		time.Sleep(100 * time.Millisecond)
		if rl.IsKeyDown(rl.KeyEnter) || rl.IsKeyDown(rl.KeyEscape) {
			break
		} else if rl.IsKeyDown(rl.KeySpace) {
		} else if rl.IsKeyDown(rl.KeyRight) {
			currSelectedPatternIndex++
		} else if rl.IsKeyDown(rl.KeyLeft) {
			currSelectedPatternIndex--
		} else if rl.IsKeyDown(rl.KeyDown) {
			w -= 16
			h -= 16
		} else if rl.IsKeyDown(rl.KeyUp) {
			w += 16
			h += 16
		} else {
			reGenerate = false
		}

		if w < map_generator.GetPatternByIndex(currSelectedPatternIndex).MinWidth {
			w = map_generator.GetPatternByIndex(currSelectedPatternIndex).MinWidth
		}
		if h < map_generator.GetPatternByIndex(currSelectedPatternIndex).MinHeight {
			h = map_generator.GetPatternByIndex(currSelectedPatternIndex).MinHeight
		}
	}
	g.battlefield.initFromRandomMap(generatedMap)
}

