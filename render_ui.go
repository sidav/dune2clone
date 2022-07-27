package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

func (r *renderer) renderUI(b *battlefield, pc *playerController) {
	r.renderResourcesUI(b, pc)
	r.renderSelectedActorUI(b, pc, 0, 3*WINDOW_H/4)
	if pc.mode == PCMODE_PLACE_BUILDING {
		r.renderBuildCursor(b, pc)
	}
}

func (r *renderer) renderResourcesUI(b *battlefield, pc *playerController) {
	rl.DrawText(fmt.Sprintf("TICK %d", b.currentTick), 0, 0, 24, rl.White)
	rl.DrawText(fmt.Sprintf("%.f$", math.Round(pc.controlledFaction.money)), WINDOW_W/2, 0, 36, rl.White)
}

func (r *renderer) renderSelectedActorUI(b *battlefield, pc *playerController, x, y int32) {
	// draw outline
	r.drawOutlinedRect(x, y, 2*int32(WINDOW_W)/5, WINDOW_H/4, 2, rl.Green, rl.Black)
	if pc.selection == nil {
		return
	}
	rl.DrawText(pc.selection.getName(), x+15, y+1, 32, rl.Green)
	// if u, ok := pc.selection.(*unit); ok {}
	if bld, ok := pc.selection.(*building); ok {
		r.renderSelectedBuildingUI(bld, x, y)
	}

}

func (r *renderer) renderBuildCursor(b *battlefield, pc *playerController) {
	tx, ty := pc.mouseCoordsToTileCoords()
	for i := 0; i < pc.cursorW; i++ {
		for j := 0; j < pc.cursorH; j++ {
			color := rl.Red
			if b.isRectClearForBuilding(tx+i, ty+j, 1, 1) {
				color = rl.Green
			}
			r.drawDitheredRect(int32((tx+i)*TILE_SIZE_IN_PIXELS), int32((ty+j) * TILE_SIZE_IN_PIXELS), TILE_SIZE_IN_PIXELS, TILE_SIZE_IN_PIXELS, color)
		}
	}
}

func (r *renderer) renderSelectedBuildingUI(bld *building, x, y int32) {
	var line int32
	if bld.currentAction.code == ACTION_WAIT {
		for _, code := range bld.getStaticData().builds {
			rl.DrawText(fmt.Sprintf("%s - Build %s ($%d)", sTableBuildings[code].hotkeyToBuild,
				sTableBuildings[code].displayedName, sTableBuildings[code].cost),
				x+4, y+1+32+32*line, 32, rl.Orange)
			line++
		}
		for _, code := range bld.getStaticData().produces {
			rl.DrawText(fmt.Sprintf("%s - Make %s ($%d)", sTableUnits[code].hotkeyToBuild,
				sTableUnits[code].displayedName,sTableUnits[code].cost),
				x+4, y+1+32+32*line, 32, rl.Orange)
			line++
		}
	}
	if bld.currentAction.code == ACTION_BUILD {
			rl.DrawText(fmt.Sprintf("Builds %s (%d%%)", bld.currentAction.targetActor.getName(), bld.currentAction.getCompletionPercent()),
				x+4, y+1+32+32, 32, rl.Orange)
	}
}