package main

import rl "github.com/gen2brain/raylib-go/raylib"

type game struct {
	battlefield battlefield
}

func (g *game) startGame() {
	g.battlefield = battlefield{}
	g.battlefield.create(MAP_W, MAP_H)
	pc := playerController{}
	r := renderer{}

	for !rl.WindowShouldClose() {
		r.renderBattlefield(&g.battlefield)
		pc.playerControl(&g.battlefield)

		if g.battlefield.currentTick % 1 == 0 {
			for i := range g.battlefield.units {
				g.battlefield.executeActionForUnit(g.battlefield.units[i])
			}
		}

		g.battlefield.currentTick++
	}
}
