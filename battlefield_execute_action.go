package main

import "math"

func (b *battlefield) executeActionForUnit(u *unit) {
	if u.currentAction.code == ACTION_WAIT {
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
			if math.Abs(tx-x) < u.getStaticData().movementSpeed {
				u.centerX = tx
			} else {
				u.centerX += u.getStaticData().movementSpeed * (tx - x) / math.Abs(tx-x)
			}
		} else if ty != y {
			if !u.rotateChassisTowardsVector(0, ty-y) {
				return
			}
			if math.Abs(ty-y) < u.getStaticData().movementSpeed {
				u.centerY = ty
			} else {
				u.centerY += u.getStaticData().movementSpeed * (ty - y) / math.Abs(ty-y)
			}
		}
		if areFloatsAlmostEqual(x, tx) && areFloatsAlmostEqual(y, ty) {
			u.centerX = tx
			u.centerY = ty
			u.currentAction.code = ACTION_WAIT
		}
	}
}
