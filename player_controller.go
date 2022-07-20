package main

import rl "github.com/gen2brain/raylib-go/raylib"

type playerController struct {
	selection actor
}

func (pc *playerController) playerControl(b *battlefield) {
	tx, ty := pc.mouseCoordsToTileCoords()
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		actr := b.getActorAtTileCoordinates(tx, ty)
		if u, ok := actr.(*unit); ok {
			u.isSelected = true
			pc.selection = u
		} else {
			if pc.selection != nil {
				pc.selection.markSelected(false)
			}
			pc.selection = nil
		}
	}
	if rl.IsMouseButtonPressed(rl.MouseRightButton) {
		if u, ok := pc.selection.(*unit); ok {
			u.currentAction.targetTileX = tx
			u.currentAction.targetTileY = ty
			u.currentAction.code = ACTION_MOVE
		}
	}
}

func (pc *playerController) mouseCoordsToTileCoords() (int, int) {
	v := rl.GetMousePosition()
	return int(v.X) / TILE_SIZE_IN_PIXELS, int(v.Y) / TILE_SIZE_IN_PIXELS
}
