package main

import rl "github.com/gen2brain/raylib-go/raylib"

type playerController struct {
	selection actor
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
}

func (pc *playerController) mouseCoordsToTileCoords() (int, int) {
	v := rl.GetMousePosition()
	return int(v.X) / TILE_SIZE_IN_PIXELS, int(v.Y) / TILE_SIZE_IN_PIXELS
}
