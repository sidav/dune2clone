package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

type playerController struct {
	camTopLeftX, camTopLeftY int // real coords, in pixels
	controlledFaction        *faction
	selection                []actor
	mode                     int
	scrollCooldown           int

	// elastic frame related
	mouseDownForTicks                int
	mouseDownCoordX, mouseDownCoordY float32

	// for drawing that "order given thing"
	tickOrderGiven, orderGivenX, orderGivenY int
}

func (pc *playerController) getFirstSelection() actor {
	if len(pc.selection) > 0 {
		return pc.selection[0]
	}
	return nil
}

func (pc *playerController) playerControl(b *battlefield) {
	// pc.mode = PCMODE_NONE

	pc.scrollCooldown--
	pc.scroll(b)

	tx, ty, fromMinimap := pc.mouseCoordsToTileCoords(b)
	if rl.IsMouseButtonPressed(rl.MouseRightButton) {
		pc.rightClickWithActorSelected(b, tx, ty)
	}
	if bld, ok := pc.getFirstSelection().(*building); ok {
		built := pc.GiveOrderToBuilding(b, bld)
		if built {
			return
		}
	}

	// selection
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		pc.mouseDownForTicks = 0
		if fromMinimap {
			pc.centerCameraAtTile(b, tx, ty)
			return
		}
		actr := b.getActorAtTileCoordinates(tx, ty)
		pc.deselect()
		if actr != nil {
			// set selection
			actr.markSelected(true)
			pc.selection = []actor{actr}
		}
	}
	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		// need to enable elastic selection?
		if pc.mouseDownForTicks > 0 {
			// debugWritef("Mouse down for %d", pc.mouseDownForTicks)
			if pc.isMouseMovedFromDownCoordinates() {
				pc.changeMode(PCMODE_ELASTIC_SELECTION)
				return
			}
		} else {
			pc.mouseDownCoordX, pc.mouseDownCoordY = rl.GetMousePosition().X, rl.GetMousePosition().Y
		}
		pc.mouseDownForTicks++
	}
	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		if pc.mode == PCMODE_ELASTIC_SELECTION {
			pc.elasticFrameSelect(b)
			pc.changeMode(PCMODE_NONE)
		}
	}
}

// returns true if order was given (for not auto-clicking, for example)
func (pc *playerController) GiveOrderToBuilding(b *battlefield, bld *building) bool {
	kk := rl.GetKeyPressed()
	if pc.IsKeyCodeEqualToString(kk, "R", true) && bld.getHitpointsPercentage() < 100 {
		bld.isRepairingSelf = true
	}
	if pc.IsKeyCodeEqualToString(kk, "S", true) {
		// TODO: appropriate selling of buildings
		bld.currentHitpoints = 0
	}
	if bld.currentOrder.code == ORDER_WAIT_FOR_BUILDING_PLACEMENT {
		if bld.currentAction.getCompletionPercent() >= 100 || bld.getStaticData().BuildType == BTYPE_PLACE_FIRST {
			// if NOT building:
			if _, ok := bld.currentAction.targetActor.(*building); !ok && bld.getStaticData().BuildType != BTYPE_PLACE_FIRST {
				return false
			}
			pc.changeMode(PCMODE_PLACE_BUILDING)
			tx, ty, _ := pc.mouseCoordsToTileCoords(b)
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) && b.canBuildingBePlacedAt(bld.currentOrder.targetActor.(*building), tx, ty, 0, false) {
				bld.currentOrder.targetTileX = tx
				bld.currentOrder.targetTileY = ty
				pc.changeMode(PCMODE_NONE)
			} else if rl.IsMouseButtonPressed(rl.MouseRightButton) {
				pc.changeMode(PCMODE_NONE)
				if bld.getStaticData().BuildType == BTYPE_PLACE_FIRST {
					bld.currentOrder.resetOrder()
				}
				pc.deselect()
			}
			return true
		}
	} else {
		if bld.currentAction.code == ACTION_WAIT {
			// maybe build?
			for _, code := range bld.getStaticData().Builds {
				if pc.IsKeyCodeEqualToString(kk, sTableBuildings[code].HotkeyToBuild, false) && bld.faction.isTechAvailableForBuildingOfCode(code) {
					bld.currentOrder.code = ORDER_BUILD
					bld.currentOrder.targetActorCode = int(code)
				}
			}
			// maybe product?
			for _, code := range bld.getStaticData().Produces {
				if pc.IsKeyCodeEqualToString(kk, sTableUnits[code].HotkeyToBuild, false) && bld.faction.isTechAvailableForUnitOfCode(code) {
					bld.currentOrder.code = ORDER_PRODUCE
					bld.currentOrder.targetActorCode = code
				}
			}
		}
		pc.changeMode(PCMODE_NONE)
	}
	return false
}

func (pc *playerController) scroll(b *battlefield) {
	if pc.scrollCooldown > 0 || pc.mode == PCMODE_ELASTIC_SELECTION || !rl.IsWindowFocused() || !rl.IsCursorOnScreen() {
		return
	}

	var SCROLL_MARGIN = float32(WINDOW_H / 8)
	const SCROLL_AMOUNT = TILE_SIZE_IN_PIXELS / int(5)
	const SCROLL_CD = 1

	v := rl.GetMousePosition()
	if areScreenCoordsOnMinimap(int32(v.X), int32(v.Y), b) {
		return
	}
	if v.X < SCROLL_MARGIN {
		pc.camTopLeftX -= SCROLL_AMOUNT
	}
	if v.X > float32(WINDOW_W)-SCROLL_MARGIN {
		pc.camTopLeftX += SCROLL_AMOUNT
	}
	if v.Y < SCROLL_MARGIN {
		pc.camTopLeftY -= SCROLL_AMOUNT
	}
	if v.Y > float32(WINDOW_H)-SCROLL_MARGIN {
		pc.camTopLeftY += SCROLL_AMOUNT
	}
	if pc.camTopLeftX < 0 {
		pc.camTopLeftX = 0
	}
	mapW, mapH := b.getSize()
	if pc.camTopLeftX > (mapW-10)*TILE_SIZE_IN_PIXELS {
		pc.camTopLeftX = (mapW - 10) * TILE_SIZE_IN_PIXELS
	}
	if pc.camTopLeftY > (mapH-10)*TILE_SIZE_IN_PIXELS {
		pc.camTopLeftY = (mapH - 10) * TILE_SIZE_IN_PIXELS
	}
	if pc.camTopLeftY < 0 {
		pc.camTopLeftY = 0
	}

	pc.scrollCooldown = SCROLL_CD
}

func (pc *playerController) rightClickWithActorSelected(b *battlefield, tx, ty int) {
	if !b.areTileCoordsValid(tx, ty) {
		return
	}
	aac := b.getActorAtTileCoordinates(tx, ty)
	for i := range pc.selection {
		if u, ok := pc.selection[i].(*unit); ok {
			pc.orderGivenX, pc.orderGivenY = tx, ty
			pc.tickOrderGiven = b.currentTick
			u.currentOrder.resetOrder()
			u.currentOrder.targetTileX = tx
			u.currentOrder.targetTileY = ty
			u.currentOrder.code = ORDER_MOVE
			if aac != nil && aac.getFaction() != u.getFaction() {
				u.currentOrder.code = ORDER_ATTACK
				u.currentOrder.targetActor = aac
				return
			}
			if len(pc.selection) == 1 && aac == pc.selection[0] && u.getStaticData().CanBeDeployed {
				u.currentOrder.code = ORDER_DEPLOY
			}
			if u.getStaticData().MaxCargoAmount > 0 && b.tiles[tx][ty].resourcesAmount > 0 {
				u.currentOrder.code = ORDER_HARVEST
			}
			if bld, ok := aac.(*building); ok {
				if bld.getStaticData().RepairsUnits {
					u.currentOrder.targetActor = bld
					u.currentOrder.code = ORDER_MOVE_TO_REPAIR
				}
				if bld.getStaticData().ReceivesResources && u.getStaticData().MaxCargoAmount > 0 {
					u.currentOrder.targetActor = bld
					u.currentOrder.code = ORDER_RETURN_TO_REFINERY
				}
			}
		}
	}
	if bld, ok := pc.getFirstSelection().(*building); ok {
		// set rally
		if bld.getStaticData().Produces != nil {
			bld.rallyTileX, bld.rallytileY = tx, ty
		}
	}
}

func (pc *playerController) centerCameraAtTile(b *battlefield, tx, ty int) {
	pc.camTopLeftX = (tx - 7) * TILE_SIZE_IN_PIXELS
	pc.camTopLeftY = (ty - 7) * TILE_SIZE_IN_PIXELS
	pc.scroll(b)
}

func (pc *playerController) elasticFrameSelect(b *battlefield) {
	pc.deselect()
	v := rl.GetMousePosition()
	x, y := (pc.camTopLeftX+int(pc.mouseDownCoordX))/TILE_SIZE_IN_PIXELS, (pc.camTopLeftY+int(pc.mouseDownCoordY))/TILE_SIZE_IN_PIXELS
	w, h := int(v.X-pc.mouseDownCoordX)/TILE_SIZE_IN_PIXELS, int(v.Y-pc.mouseDownCoordY)/TILE_SIZE_IN_PIXELS
	// debugWritef("Got %d, %d, %d, %d --- ", x, y, w, h)
	if w < 0 {
		w = -w
		x -= w
	}
	if h < 0 {
		h = -h
		y -= h
	}

	actrs := b.getListOfActorsInTilesRect(x, y, w, h)
	for i := actrs.Front(); i != nil; i = i.Next() {
		actr := i.Value.(actor)
		if actr.getFaction() == pc.controlledFaction {
			if _, ok := actr.(*unit); ok {
				i.Value.(actor).markSelected(true)
				pc.selection = append(pc.selection, i.Value.(actor))
			}
		}
	}
}

func (pc *playerController) deselect() {
	if pc.selection == nil {
		return
	}
	for i := range pc.selection {
		pc.selection[i].markSelected(false)
	}
	pc.selection = []actor{}
}

func (pc *playerController) mouseCoordsToTileCoords(b *battlefield) (int, int, bool) {
	v := rl.GetMousePosition()
	if areScreenCoordsOnMinimap(int32(v.X), int32(v.Y), b) {
		tx, ty := screenCoordsToMinimapTileCoords(int32(v.X), int32(v.Y), b)
		return tx, ty, true
	}
	return int(float32(pc.camTopLeftX)+v.X) / TILE_SIZE_IN_PIXELS, int(float32(pc.camTopLeftY)+v.Y) / TILE_SIZE_IN_PIXELS, false
}

func (pc *playerController) IsKeyCodeEqualToString(keyCode int32, keyString string, withShift bool) bool {
	//if keyCode != 0 {
	//	debugWritef("CALLED: %d - %d, diff %d\n", keyCode, int32(keyString[0]), int32(keyString[0])-keyCode)
	//}
	shiftPressed := rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)
	return int32(keyString[0])-keyCode == 0 && withShift == shiftPressed
}

func (pc *playerController) isMouseMovedFromDownCoordinates() bool {
	v := rl.GetMousePosition()
	return math.Abs(float64(pc.mouseDownCoordX-v.X)) > 0.01 && math.Abs(float64(pc.mouseDownCoordY-v.Y)) > 0.01
}

func (pc *playerController) changeMode(newMode int) {
	switch newMode {
	case PCMODE_NONE:
		pc.mouseDownForTicks = 0
	case PCMODE_PLACE_BUILDING:
		pc.mouseDownForTicks = 0
	}
	pc.mode = newMode
}

const (
	PCMODE_NONE = iota
	PCMODE_PLACE_BUILDING
	PCMODE_ELASTIC_SELECTION
)
