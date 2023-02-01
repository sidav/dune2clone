package main

import (
	"dune2clone/geometry"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"strconv"
	"strings"
	"time"
)

const UI_FONT_SIZE = 28
const BUILD_LIST_FONT_SIZE = 24
const BUILD_PANEL_WIDTH = 400
const BUILD_PANEL_HEIGHT = 500

func (r *renderer) renderUI(b *battlefield, pc *playerController) {
	r.renderResourcesUI(b, pc)
	r.renderSelectedActorUI(b, pc, 0, 3*WINDOW_H/4)
	r.drawMinimap(b, pc)
	if pc.mode == PCMODE_PLACE_BUILDING {
		r.renderBuildCursor(b, pc)
	}
	if pc.mode == PCMODE_ELASTIC_SELECTION {
		r.renderElasticSelection(b, pc)
	}
	r.renderOrderGivenAnimation(b, pc)
	// technical
	rl.DrawRectangle(0, 0, 32*WINDOW_W/100, 30*WINDOW_H/100, rl.Color{
		R: 32,
		G: 32,
		B: 32,
		A: 128,
	})
	r.drawText(fmt.Sprintf("TICK %d, frame rendered in %dms", b.currentTick, r.lastFrameRenderingTime), 0, 0, 24, rl.White)
	r.drawText(b.collectStatisticsForDebug(), 0, 28, 24, rl.White)
	for i := range r.timeDebugInfosToRender {
		color := rl.White
		ind := geometry.GetPartitionIndex(int(r.timeDebugInfosToRender[i].duration*time.Microsecond), 0, int(r.timeDebugInfosToRender[i].criticalDuration*time.Microsecond), 4)
		// debugWritef("%d from %d is %d\n", int(r.timeDebugInfosToRender[i].duration * time.Microsecond), int(r.timeDebugInfosToRender[i].criticalDuration * time.Microsecond), ind)
		if ind == 2 {
			color = rl.Yellow
		}
		if ind == 3 {
			color = rl.Red
		}
		logicNameString := fmt.Sprintf("%-18s", r.timeDebugInfosToRender[i].logicName+":")
		durationString := fmt.Sprintf("%5d (max %5d, mean %4d)",
			r.timeDebugInfosToRender[i].duration,
			r.timeDebugInfosToRender[i].maxRecordedDuration,
			r.timeDebugInfosToRender[i].calculatedMeanDuration,
		)
		r.drawText(fmt.Sprintf("%s %s", logicNameString, durationString), 0, int32(56+23*i), 18, color)
	}
}

func (r *renderer) renderOrderGivenAnimation(b *battlefield, pc *playerController) {
	const ticksForAnimation = 30
	ticksSince := b.currentTick - pc.tickOrderGiven
	completionPercent := int32(100 * ticksSince / ticksForAnimation)
	if ticksSince > ticksForAnimation {
		return
	}
	osx, osy := r.physicalToOnScreenCoords(float64(pc.orderGivenX*TILE_PHYSICAL_SIZE), float64(pc.orderGivenY*TILE_PHYSICAL_SIZE))
	if completionPercent <= 50 {
		r.drawBoldRect(
			osx+TILE_SIZE_IN_PIXELS*completionPercent/100,
			osy+TILE_SIZE_IN_PIXELS*completionPercent/100,
			TILE_SIZE_IN_PIXELS*(100-2*completionPercent)/100,
			TILE_SIZE_IN_PIXELS*(100-2*completionPercent)/100,
			3, rl.Green,
		)
	} else {
		completionPercent = 100 - completionPercent
		r.drawBoldRect(
			osx+TILE_SIZE_IN_PIXELS*completionPercent/100,
			osy+TILE_SIZE_IN_PIXELS*completionPercent/100,
			TILE_SIZE_IN_PIXELS*(100-2*completionPercent)/100,
			TILE_SIZE_IN_PIXELS*(100-2*completionPercent)/100,
			3, rl.Green,
		)
	}
}

func (r *renderer) renderResourcesUI(b *battlefield, pc *playerController) {
	factColor := pc.controlledFaction.getDarkerColor()
	// draw money
	moneyStr := fmt.Sprintf("%.f", math.Round(pc.controlledFaction.getMoney()))
	r.drawLineInfoBox(WINDOW_W-500, 0, 250, "$", moneyStr, factColor, rl.Black, rl.White)
	// draw storage
	r.drawLineInfoBox(WINDOW_W-250, 0, 250, "STRG", fmt.Sprintf("%.f", pc.controlledFaction.getStorageRemaining()),
		factColor, rl.Black, rl.White)
	// draw tech level
	techLevel := pc.controlledFaction.currTechLevel
	r.drawLineInfoBox(WINDOW_W-500, 32, 150, "TECH", fmt.Sprintf("%d", techLevel), factColor, rl.Black, rl.White)
	// draw energy
	energyStr := fmt.Sprintf("%d/%d", pc.controlledFaction.energyConsumption, pc.controlledFaction.energyProduction)
	energyBgColor := rl.Black
	energyFgColor := rl.White
	if pc.controlledFaction.getAvailableEnergy() < 0 && (b.currentTick/60)%2 == 0 {
		energyBgColor = rl.Red
		energyFgColor = rl.Black
	}
	r.drawLineInfoBox(WINDOW_W-350, 32, 350, "ENERGY", energyStr, factColor, energyBgColor, energyFgColor)
}

func (r *renderer) renderSelectedActorUI(b *battlefield, pc *playerController, x, y int32) {
	if len(pc.selection) == 0 {
		return
	}
	// draw outline
	r.drawOutlinedRect(x, y, 2*WINDOW_W/5, WINDOW_H/4, 2, pc.controlledFaction.getDarkerColor(), rl.Black)
	var lineNum int32
	r.drawText(fmt.Sprintf("%s (%d/%d)", pc.getFirstSelection().getName(),
		pc.getFirstSelection().getHitpoints(),
		pc.getFirstSelection().getMaxHitpoints()),
		x+15, y+1+lineNum*UI_FONT_SIZE, UI_FONT_SIZE, rl.Green)
	lineNum++
	r.drawText(fmt.Sprintf("%s - %s)",
		pc.getFirstSelection().getCurrentOrder().getTextDescription(),
		pc.getFirstSelection().getCurrentAction().getTextDescription()),
		x+15, y+1+lineNum*UI_FONT_SIZE, UI_FONT_SIZE, rl.Green)
	lineNum++

	expLevel := pc.getFirstSelection().getExperienceLevel()
	if expLevel > 0 {
		r.drawText(fmt.Sprintf("Experience: %d (Level %d - %s)", pc.getFirstSelection().getExperience(),
			expLevel, getExpLevelName(expLevel)),
			x+15, y+1+lineNum*UI_FONT_SIZE, UI_FONT_SIZE, rl.Green)
	}

	if u, ok := pc.getFirstSelection().(*unit); ok {
		if u.getStaticData().maxCargoAmount > 0 {
			r.drawText(fmt.Sprintf("Cargo: %d/%d", u.currentCargoAmount, u.getStaticData().maxCargoAmount),
				x+15, y+1+lineNum*UI_FONT_SIZE, UI_FONT_SIZE, rl.Green)
			lineNum++
		}
		lineNum++
		r.drawText(fmt.Sprintf("Speed %.2f Regen %d", u.getMovementSpeed(), u.getRegenAmount()),
			x+15, y+1+lineNum*UI_FONT_SIZE, UI_FONT_SIZE, rl.Green)
		lineNum++
		r.drawText(fmt.Sprintf("(%.1f, %.1f); Rotation: %d", u.centerX, u.centerY, u.chassisDegree),
			x+15, y+1+lineNum*UI_FONT_SIZE, UI_FONT_SIZE, rl.Green)
		lineNum++
	}

	if bld, ok := pc.getFirstSelection().(*building); ok {
		if len(bld.getStaticData().builds) != 0 || len(bld.getStaticData().produces) != 0 {
			r.renderSelectedBuildingUI(bld, WINDOW_W-BUILD_PANEL_WIDTH, 100)
		}
	}

}

func (r *renderer) renderBuildCursor(b *battlefield, pc *playerController) {
	if pc.getFirstSelection().(*building).currentOrder.targetActor == nil {
		return
	}
	targetBuilding := pc.getFirstSelection().(*building).currentOrder.targetActor.(*building)
	tx, ty, _ := pc.mouseCoordsToTileCoords(b)
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

func (r *renderer) renderElasticSelection(b *battlefield, pc *playerController) {
	v := rl.GetMousePosition()
	x, y := int32(pc.mouseDownCoordX), int32(pc.mouseDownCoordY)
	w, h := int32(v.X-pc.mouseDownCoordX), int32(v.Y-pc.mouseDownCoordY)
	// debugWritef("Got %d, %d, %d, %d --- ", x, y, w, h)
	if w < 0 {
		w = -w
		x -= w
	}
	if h < 0 {
		h = -h
		y -= h
	}
	// debugWritef("Drawing %d, %d, %d, %d\n", x, y, w, h)
	rl.DrawRectangleLines(x, y, w, h, rl.Gray)
	rl.DrawRectangleLines(x+1, y+1, w-2, h-2, rl.White)
	// rl.DrawLine(x+w-3, y+2, x+w-3, y+h-4, rl.Gray)
}

func (r *renderer) renderSelectedBuildingUI(bld *building, x, y int32) {
	var line int32
	r.drawOutlinedRect(x, y, BUILD_PANEL_WIDTH, BUILD_PANEL_HEIGHT, 2, rl.Green, rl.Black)
	if bld.currentAction.code == ACTION_WAIT {
		for _, code := range bld.getStaticData().builds {
			color := rl.Orange
			if !bld.faction.isTechAvailableForBuildingOfCode(code) {
				color = rl.DarkGray
			}
			r.drawText(r.collectLineForBuildMenu(sTableBuildings[code].hotkeyToBuild,
				sTableBuildings[code].displayedName, sTableBuildings[code].cost),
				x+4, y+1+BUILD_LIST_FONT_SIZE*line, BUILD_LIST_FONT_SIZE, color)
			line++
		}
		for _, code := range bld.getStaticData().produces {
			color := rl.Orange
			if !bld.faction.isTechAvailableForUnitOfCode(code) {
				color = rl.DarkGray
			}
			r.drawText(r.collectLineForBuildMenu(sTableUnits[code].hotkeyToBuild,
				sTableUnits[code].displayedName, sTableUnits[code].cost),
				x+4, y+1+BUILD_LIST_FONT_SIZE*line, BUILD_LIST_FONT_SIZE, color)
			line++
		}
	}
	if bld.currentAction.code == ACTION_BUILD {
		r.drawText(fmt.Sprintf("Builds %s (%d%%)", bld.currentAction.targetActor.getName(), bld.currentAction.getCompletionPercent()),
			x+4, y+1+UI_FONT_SIZE, UI_FONT_SIZE, rl.Orange)
	}
}

func (r *renderer) collectLineForBuildMenu(hotkey, name string, cost int) string {
	costStr := strconv.Itoa(cost)
	if len(costStr) < 5 {
		costStr += strings.Repeat(" ", 5-len(costStr))
	}
	return fmt.Sprintf("$%s - %s %s", costStr, hotkey, name)
}

func (r *renderer) renderBlinkingIconCenteredAt(iconSpriteCode string, x, y int32, blinkOrder int) {
	if (r.btl.currentTick/30)%3 == blinkOrder {
		icon := uiAtlaces[iconSpriteCode].getSpriteByFrame(0)
		rl.DrawTexture(
			icon,
			x-icon.Width/2,
			y-icon.Height/2,
			DEFAULT_TINT,
		)
	}
}
