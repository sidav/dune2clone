package main

import rl "github.com/gen2brain/raylib-go/raylib"

type playerController struct {
	controlledFaction *faction
	selection         actor
	mode              int
	cursorW, cursorH  int
}

func (pc *playerController) playerControl(b *battlefield) {
	pc.mode = PCMODE_NONE
	tx, ty := pc.mouseCoordsToTileCoords()
	if rl.IsMouseButtonPressed(rl.MouseRightButton) {
		if u, ok := pc.selection.(*unit); ok {
			u.currentOrder.targetTileX = tx
			u.currentOrder.targetTileY = ty
			u.currentOrder.code = ORDER_MOVE
		}
	}
	if bld, ok := pc.selection.(*building); ok {
		pc.GiveOrderToBuilding(b, bld)
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

func (pc *playerController) GiveOrderToBuilding(b *battlefield, bld *building) {
	kk := rl.GetKeyPressed()
	if bld.currentAction.code == ACTION_WAIT {
		// maybe build?
		for _, code := range bld.getStaticData().builds {
			if pc.IsKeyCodeEqualToString(kk, sTableBuildings[code].hotkeyToBuild) {
				bld.currentAction.code = ACTION_BUILD
				bld.currentAction.targetActor = &building{
					code:    code,
					faction: bld.faction,
				}
			}
		}
		// maybe product?
		for _, code := range bld.getStaticData().produces {
			if pc.IsKeyCodeEqualToString(kk, sTableUnits[code].hotkeyToBuild) {
				bld.currentAction.code = ACTION_BUILD
				bld.currentAction.targetActor = createUnit(code, 0, 0, bld.faction)
			}
		}
	}
	if bld.currentAction.code == ACTION_BUILD {
		if bld.currentAction.getCompletionPercent() >= 100 {
			pc.mode = PCMODE_PLACE_BUILDING
			pc.cursorW = bld.currentAction.targetActor.(*building).getStaticData().w
			pc.cursorH = bld.currentAction.targetActor.(*building).getStaticData().h
			tx, ty := pc.mouseCoordsToTileCoords()
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) && b.isRectClearForBuilding(tx, ty, pc.cursorW, pc.cursorH) {
				targetBld := bld.currentAction.targetActor.(*building)
				targetBld.topLeftX = tx
				targetBld.topLeftY = ty
				b.buildings = append(b.buildings, targetBld)
				pc.mode = PCMODE_NONE
				bld.currentAction.reset()
			}
		}
	}
}

func (pc *playerController) mouseCoordsToTileCoords() (int, int) {
	v := rl.GetMousePosition()
	return int(v.X) / TILE_SIZE_IN_PIXELS, int(v.Y) / TILE_SIZE_IN_PIXELS
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
