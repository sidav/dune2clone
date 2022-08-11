package main

import (
	"dune2clone/map_generator"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func drawGeneratedMap(gm *map_generator.GameMap) {
	rl.BeginDrawing()
	const tileSize = 8
	rl.ClearBackground(rl.Black)
	for x := range gm.Tiles {
		for y := range gm.Tiles[x] {
			color := rl.Magenta
			switch gm.Tiles[x][y] {
			case map_generator.SAND:
				color = rl.Orange
			case map_generator.RESOURCES:
				color = rl.Red
			default:
				color = rl.Brown
			}
			rl.DrawRectangle(int32(x*tileSize), int32(y*tileSize), tileSize, tileSize, color)
		}
	}

	rl.EndDrawing()
}
