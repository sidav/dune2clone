package main

import rl "github.com/gen2brain/raylib-go/raylib"

type game struct {
	battlefield battlefield
	currentTick int
}

func (g *game) startGame() {
	g.battlefield = battlefield{}
	g.battlefield.create(MAP_W, MAP_H)
	r := renderer{}

	for !rl.WindowShouldClose() {
		r.renderBattlefield(&g.battlefield)

		if g.currentTick % 2 == 0 {
			for i := range g.battlefield.units {
				g.battlefield.executeActionForUnit(g.battlefield.units[i])
			}
		}
	}
}
