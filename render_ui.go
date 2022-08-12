package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

const UI_FONT_SIZE = 28

func (r *renderer) renderUI(b *battlefield, pc *playerController) {
	r.renderResourcesUI(b, pc)
	r.renderSelectedActorUI(b, pc, 0, 3*WINDOW_H/4)
	r.drawMinimap(b, pc, WINDOW_W-256, WINDOW_H-256, 256, 256)
	if pc.mode == PCMODE_PLACE_BUILDING {
		r.renderBuildCursor(b, pc)
	}
}

func (r *renderer) renderResourcesUI(b *battlefield, pc *playerController) {
	rl.DrawText(fmt.Sprintf("TICK %d", b.currentTick), 0, 0, 24, rl.White)
	// draw money
	r.drawOutlinedRect(WINDOW_W-400, 0, 400, 36, 2, rl.Green, rl.Black)
	moneyStr := fmt.Sprintf("%.f/%.f", math.Round(pc.controlledFaction.getMoney()), pc.controlledFaction.maxResources)
	rl.DrawText("$", WINDOW_W-392, 0, 36, rl.White)
	rl.DrawText(moneyStr, WINDOW_W-int32((2*36/3)*len(moneyStr)), 0, 36, rl.White)
	// draw energy
	r.drawOutlinedRect(WINDOW_W-300, 36, 300, 36, 2, rl.Green, rl.Black)
	energyStr := fmt.Sprintf("%d/%d", pc.controlledFaction.requiredEnergy, pc.controlledFaction.currentEnergy)
	rl.DrawText("ENERGY", WINDOW_W-292, 36, 36, rl.White)
	rl.DrawText(energyStr, WINDOW_W-int32((2*36/3)*len(energyStr)), 36, 36, rl.White)
}

func (r *renderer) renderSelectedActorUI(b *battlefield, pc *playerController, x, y int32) {
	// draw outline
	r.drawOutlinedRect(x, y, 2*int32(WINDOW_W)/5, WINDOW_H/4, 2, rl.Green, rl.Black)
	if pc.selection == nil {
		return
	}
	if u, ok := pc.selection.(*unit); ok {
		rl.DrawText(fmt.Sprintf("%s (%s-%s)", u.getName(), u.currentOrder.getTextDescription(), u.getCurrentAction().getTextDescription()),
			x+15, y+1, UI_FONT_SIZE, rl.Green)
	} else {
		rl.DrawText(fmt.Sprintf("%s (%s)", pc.selection.getName(), pc.selection.getCurrentAction().getTextDescription()),
			x+15, y+1, UI_FONT_SIZE, rl.Green)
	}

	if u, ok := pc.selection.(*unit); ok {
		if u.getStaticData().maxCargoAmount > 0 {
			rl.DrawText(fmt.Sprintf("Cargo: %d/%d", u.currentCargoAmount, u.getStaticData().maxCargoAmount),
				x+15, y+UI_FONT_SIZE+1, UI_FONT_SIZE, rl.Green)
		}
	}

	if bld, ok := pc.selection.(*building); ok {
		r.renderSelectedBuildingUI(bld, x, y)
	}

}

func (r *renderer) renderBuildCursor(b *battlefield, pc *playerController) {
	if pc.selection.(*building).currentAction.targetActor == nil {
		return
	}
	targetBuilding := pc.selection.(*building).currentAction.targetActor.(*building)
	tx, ty := pc.mouseCoordsToTileCoords()
	_, _, w, h := targetBuilding.getDimensionsForConstructon()
	color := rl.Red
	if b.canBuildingBePlacedAt(targetBuilding, tx, ty, 0, false) {
		color = rl.Green
	}
	r.drawDitheredRect(int32((tx)*TILE_SIZE_IN_PIXELS)-r.camTopLeftX, int32((ty)*TILE_SIZE_IN_PIXELS)-r.camTopLeftY,
		int32(w)*TILE_SIZE_IN_PIXELS, int32(h)*TILE_SIZE_IN_PIXELS, color)
	//for i := 0; i < pc.cursorW; i++ {
	//	for j := 0; j < pc.cursorH; j++ {
	//		color := rl.Red
	//		if b.isRectClearForBuilding(tx+i, ty+j, 1, 1, pc.controlledFaction) {
	//			color = rl.Green
	//		}
	//		r.drawDitheredRect(int32((tx+i)*TILE_SIZE_IN_PIXELS)-r.camTopLeftX, int32((ty+j)*TILE_SIZE_IN_PIXELS)-r.camTopLeftY,
	//			TILE_SIZE_IN_PIXELS, TILE_SIZE_IN_PIXELS, color)
	//	}
	//}
}

func (r *renderer) renderSelectedBuildingUI(bld *building, x, y int32) {
	var line int32
	if bld.currentAction.code == ACTION_WAIT {
		for _, code := range bld.getStaticData().builds {
			rl.DrawText(fmt.Sprintf("%s - Build %s ($%d)", sTableBuildings[code].hotkeyToBuild,
				sTableBuildings[code].displayedName, sTableBuildings[code].cost),
				x+4, y+1+UI_FONT_SIZE+UI_FONT_SIZE*line, UI_FONT_SIZE, rl.Orange)
			line++
		}
		for _, code := range bld.getStaticData().produces {
			rl.DrawText(fmt.Sprintf("%s - Make %s ($%d)", sTableUnits[code].hotkeyToBuild,
				sTableUnits[code].displayedName, sTableUnits[code].cost),
				x+4, y+1+UI_FONT_SIZE+UI_FONT_SIZE*line, UI_FONT_SIZE, rl.Orange)
			line++
		}
	}
	if bld.currentAction.code == ACTION_BUILD {
		rl.DrawText(fmt.Sprintf("Builds %s (%d%%)", bld.currentAction.targetActor.getName(), bld.currentAction.getCompletionPercent()),
			x+4, y+1+UI_FONT_SIZE+UI_FONT_SIZE, UI_FONT_SIZE, rl.Orange)
	}
}
