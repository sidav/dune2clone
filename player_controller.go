package main

import rl "github.com/gen2brain/raylib-go/raylib"

type playerController struct {
	camTopLeftX, camTopLeftY int // real coords, in pixels
	controlledFaction        *faction
	selection                actor
	mode                     int
	cursorW, cursorH         int
	scrollCooldown           int
}

func (pc *playerController) playerControl(b *battlefield) {
	pc.mode = PCMODE_NONE

	pc.scrollCooldown--
	pc.scroll(b)

	tx, ty := pc.mouseCoordsToTileCoords()
	if rl.IsMouseButtonPressed(rl.MouseRightButton) {
		if u, ok := pc.selection.(*unit); ok {
			u.currentOrder.resetOrder()
			u.currentOrder.targetTileX = tx
			u.currentOrder.targetTileY = ty
			u.currentOrder.code = ORDER_MOVE
			if u.getStaticData().maxCargoAmount > 0 && b.tiles[tx][ty].resourcesAmount > 0 {
				u.currentOrder.code = ORDER_HARVEST
			}
		}
	}
	if bld, ok := pc.selection.(*building); ok {
		built := pc.GiveOrderToBuilding(b, bld)
		if built {
			return
		}
	}

	// selection
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		actr := b.getActorAtTileCoordinates(tx, ty)
		if pc.selection != nil {
			// reset selection
			pc.selection.markSelected(false)
			pc.selection = nil
		}
		if actr != nil {
			// set selection
			actr.markSelected(true)
			pc.selection = actr
		}
	}
}

// returns true if order was given (for not auto-clicking, for example)
func (pc *playerController) GiveOrderToBuilding(b *battlefield, bld *building) bool {
	kk := rl.GetKeyPressed()
	if bld.currentAction.code == ACTION_WAIT {
		// maybe build?
		for _, code := range bld.getStaticData().builds {
			if pc.IsKeyCodeEqualToString(kk, sTableBuildings[code].hotkeyToBuild) {
				bld.currentOrder.code = ORDER_BUILD
				bld.currentOrder.targetActorCode = code
			}
		}
		// maybe product?
		for _, code := range bld.getStaticData().produces {
			if pc.IsKeyCodeEqualToString(kk, sTableUnits[code].hotkeyToBuild) {
				bld.currentOrder.code = ORDER_PRODUCE
				bld.currentOrder.targetActorCode = code
			}
		}
	}
	if bld.currentOrder.code == ORDER_WAIT_FOR_BUILDING_PLACEMENT {
		if bld.currentAction.getCompletionPercent() >= 100 {
			// if NOT building:
			if _, ok := bld.currentAction.targetActor.(*building); !ok {
				return false
			}
			pc.mode = PCMODE_PLACE_BUILDING
			pc.cursorW = bld.currentAction.targetActor.(*building).getStaticData().w
			pc.cursorH = bld.currentAction.targetActor.(*building).getStaticData().h
			tx, ty := pc.mouseCoordsToTileCoords()
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) && b.canBuildingBePlacedAt(bld.currentAction.targetActor.(*building), tx, ty, 0, false) {
				bld.currentOrder.targetTileX = tx
				bld.currentOrder.targetTileY = ty
				pc.mode = PCMODE_NONE
			}
			return true
		}
	} else {
		pc.mode = PCMODE_NONE
	}
	return false
}

func (pc *playerController) scroll(b *battlefield) {
	if pc.scrollCooldown > 0 || !rl.IsWindowFocused() || !rl.IsCursorOnScreen() {
		return
	}

	var SCROLL_MARGIN = float32(WINDOW_H / 8)
	const SCROLL_AMOUNT = TILE_SIZE_IN_PIXELS / int(5)
	const SCROLL_CD = 1

	v := rl.GetMousePosition()
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

func (pc *playerController) centerCameraAtTile(b *battlefield, tx, ty int) {
	pc.camTopLeftX = (tx - 7) * TILE_SIZE_IN_PIXELS
	pc.camTopLeftY = (ty - 7) * TILE_SIZE_IN_PIXELS
	pc.scroll(b)
}

func (pc *playerController) mouseCoordsToTileCoords() (int, int) {
	v := rl.GetMousePosition()
	return int(float32(pc.camTopLeftX)+v.X) / TILE_SIZE_IN_PIXELS, int(float32(pc.camTopLeftY)+v.Y) / TILE_SIZE_IN_PIXELS
}

func (pc *playerController) IsKeyCodeEqualToString(keyCode int32, keyString string) bool {
	//if keyCode != 0 {
	//	debugWritef("CALLED: %d - %d, diff %d\n", keyCode, int32(keyString[0]), int32(keyString[0])-keyCode)
	//}
	return int32(keyString[0])-keyCode == 0
}

const (
	PCMODE_NONE = iota
	PCMODE_PLACE_BUILDING
)
