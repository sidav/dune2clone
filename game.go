package main

import rl "github.com/gen2brain/raylib-go/raylib"

type game struct {
	battlefield battlefield
}

func (g *game) startGame() {
	g.battlefield = battlefield{}
	g.battlefield.create(10, 10)
	r := renderer{}

	for !rl.WindowShouldClose() {
		r.renderBattlefield(&g.battlefield)
	}
}
