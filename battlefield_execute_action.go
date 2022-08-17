package main

import (
	"dune2clone/geometry"
	"math"
)

func (b *battlefield) executeActionForActor(a actor) {
	switch a.getCurrentAction().code {
	case ACTION_WAIT:
		if u, ok := a.(*unit); ok {
			b.executeWaitActionForUnit(u)
		}
		if bld, ok := a.(*building); ok {
			b.executeWaitActionForBuilding(bld)
		}
	case ACTION_ROTATE:
		if u, ok := a.(*unit); ok {
			b.executeRotateActionForUnit(u)
		}
	case ACTION_MOVE:
		if u, ok := a.(*unit); ok {
			b.executeMoveActionForUnit(u)
		} else {
			panic("Is not unit!")
		}
	case ACTION_BUILD:
		b.executeBuildActionForActor(a)
	case ACTION_HARVEST:
		b.executeHarvestActionForActor(a)
	case ACTION_ENTER_BUILDING:
		if u, ok := a.(*unit); ok {
			b.executeEnterBuildingActionForUnit(u)
		}

	case ACTION_AIR_APPROACH_LAND_TILE:
		if u, ok := a.(*unit); ok && u.getStaticData().isAircraft {
			b.executeAirApproachLandTileActionForUnit(u)
		}
	case ACTION_AIR_APPROACH_ACTOR:
		if u, ok := a.(*unit); ok && u.getStaticData().isAircraft {
			b.executeAirApproachTargetActorActionForUnit(u)
		}
	case ACTION_AIR_PICK_UNIT_UP:
		if u, ok := a.(*unit); ok && u.getStaticData().isAircraft {
			b.executeAirPickUnitUpActionForUnit(u)
		}
	case ACTION_AIR_DROP_UNIT:
		if u, ok := a.(*unit); ok && u.getStaticData().isAircraft {
			b.executeAirDropActionForUnit(u)
		}

	default:
		panic("No action execution func!")
	}
}

func (b *battlefield) executeWaitActionForUnit(u *unit) {
	if u.getStaticData().isAircraft {
		if u.currentAction.targetTileX == -1 {
			u.currentAction.targetTileX, u.currentAction.targetTileY = geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
		}
		targetX, targetY := geometry.TileCoordsToPhysicalCoords(u.currentAction.targetTileX, u.currentAction.targetTileY)
		vx, vy := geometry.DegreeToUnitVector(u.chassisDegree)
		u.setPhysicalCenterCoords(
			u.centerX+u.getStaticData().movementSpeed*vx,
			u.centerY+u.getStaticData().movementSpeed*vy,
		)
		orderVectorX, orderVectorY := targetX-u.centerX, targetY-u.centerY
		if !geometry.IsVectorDegreeEqualTo(orderVectorX, orderVectorY, u.chassisDegree) {
			u.rotateChassisTowardsVector(orderVectorX, orderVectorY)
			return
		}
	} else {
		// rotate the unit to target if unit's turret can't rotate itself
		if u.turret != nil && !u.turret.canRotate() && u.turret.targetActor != nil {
			x, y := u.turret.targetActor.getPhysicalCenterCoords()
			x -= u.centerX
			y -= u.centerY
			u.rotateChassisTowardsVector(x, y)
		}
	}
}

func (b *battlefield) executeWaitActionForBuilding(bld *building) {
	if bld.getStaticData().receivesResources && bld.unitPlacedInside != nil {
		if bld.unitPlacedInside.currentCargoAmount > 0 {
			// Receive resources
			const HARVEST_PER_TICK = 11
			received := min(HARVEST_PER_TICK, bld.unitPlacedInside.currentCargoAmount) // TODO: replace
			if bld.faction.getStorageRemaining() > 0 {
				received = min(received, int(bld.faction.getStorageRemaining()))
			} else {
				return
			}
			bld.getFaction().receiveResources(float64(received), false)
			bld.unitPlacedInside.currentCargoAmount -= received
		} else {
			//unitToPlace := bld.unitPlacedInside
			//x, y := bld.topLeftX, bld.topLeftY
			//x += bld.getStaticData().unitPlacementX
			//y += bld.getStaticData().unitPlacementY
			b.addActor(bld.unitPlacedInside)
			bld.unitPlacedInside = nil
		}
	}
}

func (b *battlefield) executeEnterBuildingActionForUnit(u *unit) {
	if u.currentAction.targetActor.(*building).unitPlacedInside == nil {
		u.currentAction.code = ACTION_WAIT
		u.chassisDegree = 90
		b.removeActor(u)
		u.currentAction.targetActor.(*building).unitPlacedInside = u
	}
}

func (b *battlefield) executeHarvestActionForActor(a actor) {
	const HARVEST_PER_TICK = 3
	if u, ok := a.(*unit); ok {
		x, y := u.getPhysicalCenterCoords()
		utx, uty := geometry.TrueCoordsToTileCoords(x, y)
		if u.currentCargoAmount < u.getStaticData().maxCargoAmount && b.tiles[utx][uty].resourcesAmount > 0 {
			harvestedAmount := min(b.tiles[utx][uty].resourcesAmount, HARVEST_PER_TICK) // TODO: replace
			harvestedAmount = min(harvestedAmount, u.getStaticData().maxCargoAmount-u.currentCargoAmount)
			b.tiles[utx][uty].resourcesAmount -= harvestedAmount
			u.currentCargoAmount += harvestedAmount
		} else {
			u.currentAction.code = ACTION_WAIT
		}
	}
}

func (b *battlefield) executeRotateActionForUnit(u *unit) {
	if u.turret != nil && u.turret.canRotate() {
		if u.turret.rotationDegree == u.currentAction.targetRotation {
			u.currentAction.code = ACTION_WAIT
		} else if u.turret.targetActor == nil {
			u.turret.rotationDegree += geometry.GetDiffForRotationStep(u.turret.rotationDegree, u.currentAction.targetRotation, u.turret.getStaticData().rotateSpeed)
			u.normalizeDegrees()
		}
	} else {
		if u.chassisDegree == u.currentAction.targetRotation {
			u.currentAction.code = ACTION_WAIT
		} else {
			u.chassisDegree += geometry.GetDiffForRotationStep(u.chassisDegree, u.currentAction.targetRotation, u.getStaticData().chassisRotationSpeed)
			u.normalizeDegrees()
		}
	}
}

func (b *battlefield) executeMoveActionForUnit(u *unit) {
	x, y := u.centerX, u.centerY
	targetX, targetY := geometry.TileCoordsToPhysicalCoords(u.currentAction.targetTileX, u.currentAction.targetTileY)
	if u.getStaticData().isAircraft {
		vx, vy := geometry.DegreeToUnitVector(u.chassisDegree)
		u.setPhysicalCenterCoords(
			u.centerX+u.getStaticData().movementSpeed*vx,
			u.centerY+u.getStaticData().movementSpeed*vy,
		)
		orderVectorX, orderVectorY := targetX-u.centerX, targetY-u.centerY
		if !geometry.IsVectorDegreeEqualTo(orderVectorX, orderVectorY, u.chassisDegree) {
			u.rotateChassisTowardsVector(orderVectorX, orderVectorY)
		}
		tx, ty := geometry.TrueCoordsToTileCoords(u.getPhysicalCenterCoords())
		if tx == u.currentAction.targetTileX && ty == u.currentAction.targetTileY {
			u.currentAction.reset()
		}
	} else {
		vx, vy := targetX-x, targetY-y
		if areFloatsAlmostEqual(x, targetX) && areFloatsAlmostEqual(y, targetY) {
			u.centerX = targetX
			u.centerY = targetY
			u.currentAction.reset()
			// debugWritef("Tick %d: action finished\n", b.currentTick)
		}

		if !geometry.IsVectorDegreeEqualTo(vx, vy, u.chassisDegree) {
			u.rotateChassisTowardsVector(vx, vy)
			return
		}

		if math.Abs(vx) < u.getStaticData().movementSpeed {
			u.centerX = targetX // source of possible movement lag :(
		} else {
			u.centerX += u.getStaticData().movementSpeed * vx / math.Abs(vx)
		}

		if math.Abs(vy) < u.getStaticData().movementSpeed {
			u.centerY = targetY // source of possible movement lag :(
		} else {
			u.centerY += u.getStaticData().movementSpeed * vy / math.Abs(vy)
		}
	}
}

func (b *battlefield) executeBuildActionForActor(a actor) {
	act := a.getCurrentAction()
	moneySpent := 0.0
	// calculate spending
	if bld, ok := act.targetActor.(*building); ok {
		moneySpent = float64(bld.getStaticData().cost) /
			float64(bld.getStaticData().buildTime*(DESIRED_FPS/BUILDINGS_ACTIONS_TICK_EACH))
	}
	if unt, ok := act.targetActor.(*unit); ok {
		moneySpent = float64(unt.getStaticData().cost) /
			float64(unt.getStaticData().buildTime*(DESIRED_FPS/BUILDINGS_ACTIONS_TICK_EACH))
	}
	// spend money
	coeff := a.getFaction().getEnergyProductionMultiplier()
	moneySpent *= coeff
	if act.getCompletionPercent() < 100 && a.getFaction().getMoney() > moneySpent {
		a.getFaction().spendMoney(moneySpent)
		act.completionAmount += coeff
	}
	// if it was a unit, place it right away
	if unt, ok := act.targetActor.(*unit); ok && act.getCompletionPercent() >= 100 {
		if bld, ok := a.(*building); ok {
			for x := bld.topLeftX; x < bld.topLeftX+bld.getStaticData().w; x++ {
				// for y := bld.topLeftY-1; y <= bld.topLeftY+bld.getStaticData().h; y++ {
				y := bld.topLeftY + bld.getStaticData().h
				if b.costMapForMovement(x, y) != -1 {
					unt.centerX, unt.centerY = geometry.TileCoordsToPhysicalCoords(x, y)
					if bld.rallyTileX != -1 && unt.currentOrder.code == ORDER_NONE {
						unt.currentOrder.code = ORDER_MOVE
						unt.currentOrder.setTargetTileCoords(bld.rallyTileX, bld.rallytileY)
						if !unt.getStaticData().isAircraft {
							bld.faction.addDispatchRequest(unt, bld.rallyTileX, bld.rallytileY, ORDER_CARRY_UNIT_TO_TARGET_COORDS, b.currentTick+100)
						}
					}
					b.addActor(unt)
					bld.currentAction.reset()
					return
				}
				// }
			}
		} else {
			panic("wat")
		}
	}
}
