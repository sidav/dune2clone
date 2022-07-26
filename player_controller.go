package main

import rl "github.com/gen2brain/raylib-go/raylib"

type playerController struct {
	controlledFaction *faction
	selection         actor
}

func (pc *playerController) playerControl(b *battlefield) {
	tx, ty := pc.mouseCoordsToTileCoords()
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
	if rl.IsMouseButtonPressed(rl.MouseRightButton) {
		if u, ok := pc.selection.(*unit); ok {
			u.currentOrder.targetTileX = tx
			u.currentOrder.targetTileY = ty
			u.currentOrder.code = ORDER_MOVE
		}
	}
	if b, ok := pc.selection.(*building); ok {
		pc.GiveOrderToBuilding(b)
	}
}

func (pc *playerController) GiveOrderToBuilding(b *building) {
	kk := rl.GetKeyPressed()
	if b.currentAction.code == ACTION_WAIT {
		// maybe build?
		if len(b.getStaticData().builds) > 0 {
			for _, code := range b.getStaticData().builds {
				if pc.IsKeyCodeEqualToString(kk, sTableBuildings[code].hotkeyToBuild) {
					b.currentAction.code = ACTION_BUILD
					b.currentAction.targetActor = &building{
						code:          code,
						faction:       b.faction,
					}
				}
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
	return int32(keyString[0]) - keyCode == 0
}
