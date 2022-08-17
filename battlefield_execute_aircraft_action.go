package main

import "dune2clone/geometry"

func (b *battlefield) executeAirApproachLandTileActionForUnit(u *unit) {
	const rangeToDropSpeed = 2
	atx, aty := geometry.TrueCoordsToTileCoords(u.getPhysicalCenterCoords())
	if geometry.GetApproxDistFromTo(atx, aty, u.currentAction.targetTileX, u.currentAction.targetTileY) >= rangeToDropSpeed {
		b.executeMoveActionForUnit(u)
		return
	}
	if u.isPresentAt(u.currentAction.targetTileX, u.currentAction.targetTileY) {
		u.currentAction.reset()
		return
	}
	newSpeed := 3 * u.getStaticData().movementSpeed / 4
	targetTrueX, targetTrueY := geometry.TileCoordsToPhysicalCoords(u.currentAction.targetTileX, u.currentAction.targetTileY)
	vx, vy := geometry.VectorToUnitVectorFloat64(targetTrueX-u.centerX, targetTrueY-u.centerY)
	u.rotateChassisTowardsVector(vx, vy)
	u.setPhysicalCenterCoords(
		u.centerX+newSpeed*vx,
		u.centerY+newSpeed*vy,
	)
}

func (b *battlefield) executeAirApproachTargetActorActionForUnit(u *unit) {
	const rangeToDropSpeed = 1
	targetTrueX, targetTrueY := u.currentAction.targetActor.getPhysicalCenterCoords()
	targetTileX, targetTileY := geometry.TrueCoordsToTileCoords(targetTrueX, targetTrueY)
	apx, apy := u.getPhysicalCenterCoords()
	atx, aty := geometry.TrueCoordsToTileCoords(u.getPhysicalCenterCoords())
	if geometry.GetApproxDistFromTo(atx, aty, targetTileX, targetTileY) >= rangeToDropSpeed {
		b.executeMoveActionForUnit(u)
		return
	}
	u.rotateChassisTowardsDegree(u.currentAction.targetActor.(*unit).chassisDegree)
	if areFloatsAlmostEqual(targetTrueX, apx) && areFloatsAlmostEqual(targetTrueY, apy) {
		u.currentAction.reset()
		return
	}
	newSpeed := 3 * u.getStaticData().movementSpeed / 4

	vx, vy := geometry.VectorToUnitVectorFloat64(targetTrueX-u.centerX, targetTrueY-u.centerY)
	u.setPhysicalCenterCoords(
		u.centerX+newSpeed*vx,
		u.centerY+newSpeed*vy,
	)
}

func (b *battlefield) executeAirPickUnitUpActionForUnit(u *unit) {
	const rangeToLockOn = 2
	debugWrite("PICKING UP")
	atx, aty := geometry.TrueCoordsToTileCoords(u.getPhysicalCenterCoords())
	ttx, tty := geometry.TrueCoordsToTileCoords(u.currentAction.targetActor.(*unit).getPhysicalCenterCoords())
	if geometry.GetApproxDistFromTo(atx, aty, ttx, tty) > rangeToLockOn {
		debugWrite("FLYING TO")
		u.currentAction.targetTileX, u.currentAction.targetTileY = ttx, tty
		b.executeAirApproachTargetActorActionForUnit(u)
		return
	}
	debugWrite("MOVING TO LOCATION")
	targetTrueX, targetTrueY := u.currentAction.targetActor.getPhysicalCenterCoords()
	if u.isPresentAt(ttx, tty) {
		u.centerX, u.centerY = targetTrueX, targetTrueY
		u.currentAction.targetActor.getCurrentAction().reset()
		u.carriedUnit = u.currentAction.targetActor.(*unit)
		b.removeActor(u.currentAction.targetActor)
		u.currentAction.reset()
	} else {
		newSpeed := u.getStaticData().movementSpeed / 2
		vx, vy := geometry.VectorToUnitVectorFloat64(targetTrueX-u.centerX, targetTrueY-u.centerY)
		u.centerX += newSpeed * vx
		u.centerY += newSpeed * vy
		u.rotateChassisTowardsDegree(u.currentAction.targetActor.(*unit).chassisDegree)
	}
}

func (b *battlefield) executeAirDropActionForUnit(u *unit) {
	debugWrite("DROP: STARTING")
	atx, aty := geometry.TrueCoordsToTileCoords(u.getPhysicalCenterCoords())
	if b.isTileClearToBeMovedInto(atx, aty, u.carriedUnit) {
		u.carriedUnit.setPhysicalCenterCoords(geometry.TileCoordsToPhysicalCoords(atx, aty))
		b.addActor(u.carriedUnit)
		u.carriedUnit = nil
		u.currentAction.reset()
	}
}
