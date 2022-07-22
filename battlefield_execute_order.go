package main

import "math"

func (b *battlefield) executeOrderForUnit(u *unit) {
	if u.currentAction.code != ACTION_WAIT {
		return // execute the order only after finishing the action
	}
	if u.currentOrder.code == ORDER_NONE {
		return // TODO: look around for targets etc
	}

	switch u.currentOrder.code {
	case ORDER_MOVE:
		b.executeMoveOrder(u)
	}
}

func (b *battlefield) executeMoveOrder(u *unit) {
	x, y := u.centerX, u.centerY
	// todo: pathfinding
	tx, ty := tileCoordsToPhysicalCoords(u.currentAction.targetTileX, u.currentAction.targetTileY)

	if tx != x {
		if !u.rotateChassisTowardsVector(tx-x, 0) {
			return
		}
		if math.Abs(tx-x) < u.getStaticData().movementSpeed {
			u.centerX = tx
			return
		}
		u.centerX += u.getStaticData().movementSpeed * (tx-x)/math.Abs(tx-x)
	} else if ty != y {
		if !u.rotateChassisTowardsVector(0, ty-y) {
			return
		}
		if math.Abs(ty-y) < u.getStaticData().movementSpeed {
			u.centerY = ty
			return
		}
		u.centerY += u.getStaticData().movementSpeed * (ty-y)/math.Abs(ty-y)
	}
	if areFloatsAlmostEqual(x, tx) && areFloatsAlmostEqual(y, ty) {
		u.centerX = tx
		u.centerY = ty
		u.currentAction.code = ACTION_WAIT
	}
}
