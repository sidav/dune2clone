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

		if tx != x {
			if !u.rotateChassisTowardsVector(tx-x, 0) {
				return
			}
			u.centerX += u.getStaticData().movementSpeed * (tx-x)/math.Abs(tx-x)
			if areFloatsAlmostEqual(x, tx) {
				u.centerX = tx
			}
		} else if ty != y {
			if !u.rotateChassisTowardsVector(0, ty-y) {
				return
			}
			u.centerY += u.getStaticData().movementSpeed * (ty-y)/math.Abs(ty-y)
			if areFloatsAlmostEqual(y, ty) {
				u.centerY = ty
			}
		}
	}
}
