package main

import "math"

func (b *battlefield) executeActionForUnit(u *unit) {
	if u.currentAction == nil {
		return // TODO: look around for targets etc
	}

	x, y := u.centerX, u.centerY

	switch u.currentAction.code {
	case ACTION_MOVE:
		// todo: pathfinding
		tx, ty := tileCoordsToPhysicalCoords(u.currentAction.targetTileX, u.currentAction.targetTileY)
		if (tx-x)*(tx-x) + (ty-y)*(ty-y) <= u.getStaticData().speed * u.getStaticData().speed {
			u.centerX, u.centerY = tx, ty
			u.currentAction = nil
			return
		}
		if tx != x {
			u.centerX += u.getStaticData().speed * (tx-x)/math.Abs(tx-x)
		} else if ty != y {
			u.centerY += u.getStaticData().speed * (ty-y)/math.Abs(ty-y)
		}
	}
}
