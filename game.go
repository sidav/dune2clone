package main

import rl "github.com/gen2brain/raylib-go/raylib"

type game struct {
	battlefield battlefield
}

func (g *game) startGame() {
	g.battlefield = battlefield{}
	g.battlefield.create(MAP_W, MAP_H)
	r := renderer{}

	for !rl.WindowShouldClose() {
		r.renderBattlefield(&g.battlefield)

		if g.battlefield.currentTick % 1 == 0 {
			for i := range g.battlefield.units {
				g.battlefield.executeActionForUnit(g.battlefield.units[i])
			}

			// TODO: remove this block
			for i := range g.battlefield.units {
				if g.battlefield.units[i].currentAction == nil {
					g.battlefield.units[i].currentAction = &action{
						code:        ACTION_MOVE,
						targetTileX: rnd.Rand(MAP_W),
						targetTileY: rnd.Rand(MAP_H),
					}
				}
			}
		}

		g.battlefield.currentTick++
	}
}
