package main

import (
	"math"
)

func (b *battlefield) executeActionForUnit(u *unit) {
	switch u.currentAction.code {
	case ACTION_MOVE:
		b.executeMoveActionForUnit(u)
	}
}

func (b *battlefield) executeMoveActionForUnit(u *unit) {
	x, y := u.centerX, u.centerY
	tx, ty := tileCoordsToPhysicalCoords(u.currentAction.targetTileX, u.currentAction.targetTileY)

	if tx != x {
		if !isVectorDegreeEqualTo(tx-x, 0, u.chassisDegree) {
			u.rotateChassisTowardsVector(tx-x, 0)
			return
		}
		if math.Abs(tx-x) < u.getStaticData().movementSpeed {
			u.centerX = tx // source of movement lag :(
		} else {
			u.centerX += u.getStaticData().movementSpeed * (tx - x) / math.Abs(tx-x)
		}
	} else if ty != y {
		if !isVectorDegreeEqualTo(0, ty-y, u.chassisDegree) {
			u.rotateChassisTowardsVector(0, ty-y)
			return
		}
		if math.Abs(ty-y) < u.getStaticData().movementSpeed {
			u.centerY = ty // source of movement lag :(
		} else {
			u.centerY += u.getStaticData().movementSpeed * (ty - y) / math.Abs(ty-y)
		}
	}
	if areFloatsAlmostEqual(x, tx) && areFloatsAlmostEqual(y, ty) {
		u.centerX = tx
		u.centerY = ty
		u.currentAction.code = ACTION_WAIT
		// debugWritef("Tick %d: action finished\n", b.currentTick)
	}
}
