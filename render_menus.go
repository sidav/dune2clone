package main

import (
	"dune2clone/map_generator"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"strconv"
	"time"
)

func (r *renderer) drawGeneratedMap(gm *map_generator.GeneratedMap, patternIndex int) {
	rl.BeginDrawing()
	rl.DrawText("Select map. SPACE to generate new, ENTER to select current.", 0, 0, 28, rl.White)
	rl.DrawText(fmt.Sprintf("UP and DOWN to change size (current: %dx%d)", len(gm.Tiles), len(gm.Tiles[0])), 0, 30, 28, rl.White)
	rl.DrawText("<- and -> to change map pattern", 0, 60, 28, rl.White)
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

func (r *renderer) drawStartSelectionMenu(startPointsNum int) []*startConditions {
	const menuEntrySize = 32
	scs := make([]*startConditions, startPointsNum)
	for i := range scs {
		scs[i] = &startConditions{
			aiType:              "random",
			factionName:         "Random",
			resourcesMultiplier: 4,
		}
	}
	scs[0].aiType = "player"
	scs[0].resourcesMultiplier = 1

	allAiPersonalities := getListOfAIPersonalities()
	allAiPersonalities = append([]string{"player", "random"}, allAiPersonalities...)
	allFactions := []string{
		"Random",
		"Commonwealth",
		"BetaCorp",
	}
	cursor := 0

	helpStrs := []string{
		"Use LEFT and RIGHT to change resource bonus amount",
		"Use A to change between player/AIs",
		"Use F to change between factions",
	}

	for {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawText("Select game start conditions.", 0, 0, 36, rl.White)
		for i, s := range helpStrs {
			rl.DrawText(s, 0, int32(rl.GetScreenHeight()-36*(i+1)), 36, rl.White)
		}
		for i, sc := range scs {
			aiString := sc.aiType
			if sc.aiType != "player" {
				aiString = "AI " + aiString
			}

			totalString := fmt.Sprintf("%d: %s (%s) - resources x%.1f", i+1, sc.factionName, aiString, sc.resourcesMultiplier)
			textColor := rl.White
			if i == cursor {
				width := rl.MeasureText(totalString, menuEntrySize)
				rl.DrawRectangle(3, 40+(menuEntrySize-2)*int32(i), width+10, menuEntrySize+4, rl.Red)
				textColor = rl.Black
			}
			rl.DrawText(totalString, 7, 40+int32(menuEntrySize*i), menuEntrySize, textColor)
		}
		rl.EndDrawing()

		key := rl.GetKeyPressed()
		switch key {
		case rl.KeyEnter:
			return scs
		case rl.KeyF:
			scs[cursor].factionName = selectStringInArrayAfter(allFactions, scs[cursor].factionName)
		case rl.KeyA:
			scs[cursor].aiType = selectStringInArrayAfter(allAiPersonalities, scs[cursor].aiType)
		case rl.KeyUp:
			cursor--
		case rl.KeyDown:
			cursor++
		case rl.KeyLeft:
			scs[cursor].resourcesMultiplier -= 0.5
			if scs[cursor].resourcesMultiplier < 0.5 {
				scs[cursor].resourcesMultiplier = 0.5
			}
		case rl.KeyRight:
			scs[cursor].resourcesMultiplier += 0.5
		}
		if cursor < 0 {
			cursor = startPointsNum - 1
		}
		if cursor >= startPointsNum {
			cursor = 0
		}
		time.Sleep(100 * time.Millisecond)
	}
	panic("Wat")
}

func (r *renderer) drawLoadingScreen(msg string) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawText(fmt.Sprintf("%30s", msg), 0, WINDOW_H/2-40, 80, rl.White)
	rl.EndDrawing()
}

// Very bad code :(
func selectStringInArrayAfter(arr []string, curr string) string {
	for i := range arr {
		if arr[i] == curr {
			return arr[(i+1)%len(arr)]
		}
	}
	panic("Wat")
}
