package main

import (
	"dune2clone/map_generator"
	rl "github.com/gen2brain/raylib-go/raylib"
	"strings"
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
			ch := make(chan bool, 1)
			go generatedMap.Generate(w, h, currSelectedPatternIndex, ch)
			generated := false
			for !generated {
				dots := (generatedMap.GenerationTries / 10) % 4
				g.render.drawLoadingScreen("Generating" + strings.Repeat(".", dots) + strings.Repeat(" ", 4-dots))
				if len(ch) > 0 {
					generated = <-ch
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
		reGenerate = true
		g.render.drawGeneratedMap(generatedMap, currSelectedPatternIndex)
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
	scs := g.render.drawStartSelectionMenu(len(generatedMap.StartPoints))
	g.battlefield.initFromRandomMap(generatedMap, scs)
}
