package main

import (
	"dune2clone/map_generator"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"strconv"
)

func drawGeneratedMap(gm *map_generator.GeneratedMap, patternIndex int) {
	rl.BeginDrawing()
	rl.DrawText("Select map. SPACE to generate new, ENTER to select current.", 0, 0, 28, rl.White)
	rl.DrawText("UP and DOWN to change map pattern", 0, 30, 28, rl.White)
	rl.DrawText(fmt.Sprintf("<- and -> to change size (current: %dx%d)", len(gm.Tiles), len(gm.Tiles[0])), 0, 60, 28, rl.White)
	rl.DrawText(fmt.Sprintf("%35s", map_generator.GetPatternByIndex(patternIndex).Name), 0, 90, 36, rl.Gold)
	offset := int32(128)
	var tileSize = int((WINDOW_H - offset) / int32(len(gm.Tiles[0])))
	rl.ClearBackground(rl.Black)
	for x := range gm.Tiles {
		for y := range gm.Tiles[x] {
			color := rl.Magenta
			switch gm.Tiles[x][y] {
			case map_generator.SAND:
				color = rl.Orange
			case map_generator.RESOURCE_VEIN:
				color = rl.DarkGreen
			case map_generator.POOR_RESOURCES, map_generator.MEDIUM_RESOURCES:
				color = rl.Purple
			case map_generator.RICH_RESOURCES:
				color = rl.DarkPurple
			case map_generator.ROCKS:
				color = rl.DarkBrown
			default:
				color = rl.Brown
			}
			rl.DrawRectangle(int32(x*tileSize), offset+int32(y*tileSize), int32(tileSize), int32(tileSize), color)
		}
	}
	for sp := range gm.StartPoints {
		spSize := int32(tileSize) * 4
		rl.DrawRectangle(int32(tileSize*gm.StartPoints[sp][0])-spSize/3, offset+int32(tileSize*gm.StartPoints[sp][1]), spSize, spSize, rl.Black)
		rl.DrawText(strconv.Itoa(sp+1), int32(tileSize*gm.StartPoints[sp][0]), offset+int32(tileSize*gm.StartPoints[sp][1]), spSize, factionColors[sp])
	}

	rl.EndDrawing()
}

func drawLoadingScreen(msg string) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawText(fmt.Sprintf("%30s", msg), 0, WINDOW_H/2-40, 80, rl.White)
	rl.EndDrawing()
}
