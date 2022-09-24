package main

import (
	"dune2clone/map_generator"
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
)

func (g *game) selectMapToGenerateBattlefield() {
	map_generator.SetRandom(&rnd)
	generatedMap := &map_generator.GeneratedMap{}
	w := 64
	h := 64
	generatedMap.Generate(w, h)
	for {
		rl.BeginDrawing()
		drawGeneratedMap(generatedMap)
		rl.EndDrawing()
		time.Sleep(100 * time.Millisecond)
		if rl.IsKeyDown(rl.KeyEnter) || rl.IsKeyDown(rl.KeyEscape) {
			break
		} else if rl.IsKeyDown(rl.KeySpace) {
			generatedMap.Generate(w, h)
		} else if rl.IsKeyDown(rl.KeyRight) {
			w += 16
			h += 16
			generatedMap.Generate(w, h)
		} else if rl.IsKeyDown(rl.KeyLeft) {
			w -= 16
			h -= 16
			if w < 32 {
				w = 32
			}
			if h < 32 {
				h = 32
			}
			generatedMap.Generate(w, h)
		}
	}
	g.battlefield.initFromRandomMap(generatedMap)
}

