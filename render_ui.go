package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (r *renderer) renderUI(b *battlefield, pc *playerController) {
	r.renderResourcesUI(b, pc)
	r.renderSelectedActorUI(b, pc, 0, 3*WINDOW_H/4)
}

func (r *renderer) renderResourcesUI(b *battlefield, pc *playerController) {
	rl.DrawText(fmt.Sprintf("TICK %d", b.currentTick), 0, 0, 24, rl.White)
	rl.DrawText(fmt.Sprintf("%d$", pc.controlledFaction.money), WINDOW_W/2, 0, 36, rl.White)
}

func (r *renderer) renderSelectedActorUI(b *battlefield, pc *playerController, x, y int32) {
	// draw outline
	r.drawOutlinedRect(x, y, int32(WINDOW_W)/3, WINDOW_H/4, 2, rl.Green, rl.Black)
	if pc.selection == nil {
		return
	}
	rl.DrawText(pc.selection.getName(), x+15, y+1, 32, rl.Green)
	// if u, ok := pc.selection.(*unit); ok {}
	if bld, ok := pc.selection.(*building); ok {
		r.renderSelectedBuildingUI(bld, x, y)
	}

}

func (r *renderer) renderSelectedBuildingUI(bld *building, x, y int32) {
	for i, code := range bld.getStaticData().builds {
		rl.DrawText(fmt.Sprintf("%s - Build %s ($%d)", sTableBuildings[code].hotkeyToBuild,
			sTableBuildings[code].displayedName, sTableBuildings[code].cost),
			x+4, y+1+32+32*int32(i), 32, rl.Orange)
	}
}

func (r *renderer) drawOutlinedRect(x, y, w, h, outlineThickness int32, outlineColor, fillColor rl.Color) {
	// draw outline
	for i := int32(0); i < outlineThickness; i++ {
		rl.DrawRectangleLines(x+i, y+i, w-i*outlineThickness, h-i*outlineThickness, outlineColor)
	}
	rl.DrawRectangle(x+outlineThickness, y+outlineThickness, w-outlineThickness*2, h-outlineThickness*2, fillColor)
}
