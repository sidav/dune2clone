package main

import (
	"dune2clone/geometry"
)

func (u *unit) approachTrueCoordinatesAsAir(x, y, speed float64) {
	apx, apy := u.getPhysicalCenterCoords()
	if areFloatsAlmostEqual(x, apx) && areFloatsAlmostEqual(y, apy) {
		u.currentAction.resetAction()
		return
	}
	vx, vy := geometry.VectorToUnitVectorFloat64(x-u.centerX, y-u.centerY)
	u.setPhysicalCenterCoords(
		u.centerX+speed*vx,
		u.centerY+speed*vy,
	)
}

func (b *battlefield) executeAirMoveActionForUnit(u *unit) {
	targetX, targetY := geometry.TileCoordsToTrueCoords(u.currentAction.targetTileX, u.currentAction.targetTileY)
	vx, vy := geometry.DegreeToUnitVector(u.chassisDegree)
	movementSpeed := u.getMovementSpeed()
	u.setPhysicalCenterCoords(
		u.centerX+movementSpeed*vx,
		u.centerY+movementSpeed*vy,
	)
	orderVectorX, orderVectorY := targetX-u.centerX, targetY-u.centerY
	if !geometry.IsVectorDegreeEqualTo(orderVectorX, orderVectorY, u.chassisDegree) {
		u.rotateChassisTowardsVector(orderVectorX, orderVectorY)
	}
	tx, ty := geometry.TrueCoordsToTileCoords(u.getPhysicalCenterCoords())
	if tx == u.currentAction.targetTileX && ty == u.currentAction.targetTileY {
		u.currentAction.resetAction()
	}
}

func (b *battlefield) executeAirApproachLandTileActionForUnit(u *unit) {
	const rangeToDropSpeed = 2
	atx, aty := geometry.TrueCoordsToTileCoords(u.getPhysicalCenterCoords())
	if geometry.GetApproxDistFromTo(atx, aty, u.currentAction.targetTileX, u.currentAction.targetTileY) >= rangeToDropSpeed {
		b.executeAirMoveActionForUnit(u)
		return
	}
	newSpeed := 3 * u.getMovementSpeed() / 4
	targetTrueX, targetTrueY := geometry.TileCoordsToTrueCoords(u.currentAction.targetTileX, u.currentAction.targetTileY)
	vx, vy := geometry.VectorToUnitVectorFloat64(targetTrueX-u.centerX, targetTrueY-u.centerY)
	u.rotateChassisTowardsVector(vx, vy)
	u.approachTrueCoordinatesAsAir(targetTrueX, targetTrueY, newSpeed)
}

func (b *battlefield) executeAirApproachTargetActorActionForUnit(u *unit) {
	const rangeToDropSpeed = 1
	targetTrueX, targetTrueY := u.currentAction.targetActor.getPhysicalCenterCoords()
	targetTileX, targetTileY := geometry.TrueCoordsToTileCoords(targetTrueX, targetTrueY)
	atx, aty := geometry.TrueCoordsToTileCoords(u.getPhysicalCenterCoords())
	if geometry.GetApproxDistFromTo(atx, aty, targetTileX, targetTileY) > rangeToDropSpeed {
		b.executeAirMoveActionForUnit(u)
		return
	}
	u.rotateChassisTowardsDegree(u.currentAction.targetActor.(*unit).chassisDegree)
	newSpeed := 3 * u.getMovementSpeed() / 4
	u.approachTrueCoordinatesAsAir(targetTrueX, targetTrueY, newSpeed)
}

func (b *battlefield) executeAirPickUnitUpActionForUnit(u *unit) {
	// debugWrite("PICKING UP")
	ttx, tty := geometry.TrueCoordsToTileCoords(u.currentAction.targetActor.(*unit).getPhysicalCenterCoords())
	// debugWrite("MOVING TO LOCATION")
	targetTrueX, targetTrueY := u.currentAction.targetActor.getPhysicalCenterCoords()
	if u.isPresentAt(ttx, tty) {
		u.centerX, u.centerY = targetTrueX, targetTrueY
		u.currentAction.targetActor.getCurrentAction().resetAction()
		u.carriedUnit = u.currentAction.targetActor.(*unit)
		b.removeActor(u.currentAction.targetActor)
		u.currentAction.resetAction()
	} else {
		newSpeed := u.getMovementSpeed() / 2
		vx, vy := geometry.VectorToUnitVectorFloat64(targetTrueX-u.centerX, targetTrueY-u.centerY)
		u.centerX += newSpeed * vx
		u.centerY += newSpeed * vy
		u.rotateChassisTowardsDegree(u.currentAction.targetActor.(*unit).chassisDegree)
	}
}

func (b *battlefield) executeAirDropActionForUnit(u *unit) {
	debugWrite("DROP: STARTING")
	atx, aty := geometry.TrueCoordsToTileCoords(u.getPhysicalCenterCoords())
	u.carriedUnit.setPhysicalCenterCoords(geometry.TileCoordsToTrueCoords(atx, aty))
	b.addActor(u.carriedUnit)
	u.carriedUnit = nil
	u.currentAction.resetAction()
}
