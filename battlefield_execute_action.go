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
			if u.getStaticData().isAircraft {
				b.executeAirMoveActionForUnit(u)
			} else {
				b.executeGroundMoveActionForUnit(u)
			}
		} else {
			panic("Is not unit!")
		}
	case ACTION_BUILD:
		b.executeBuildActionForActor(a)
	case ACTION_BEING_BUILT:
		b.executeBeingBuiltActionForBuilding(a.(*building))
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
		if u.turrets != nil && !u.turrets[0].canRotate() && u.turrets[0].targetActor != nil {
			x, y := u.turrets[0].targetActor.getPhysicalCenterCoords()
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
	if bld.getStaticData().repairsUnits && bld.unitPlacedInside != nil {
		if bld.unitPlacedInside.currentHitpoints < bld.unitPlacedInside.getStaticData().maxHitpoints {
			// Receive resources
			const REPAIR_PER_TICK = 2
			bld.unitPlacedInside.currentHitpoints = REPAIR_PER_TICK
			if bld.unitPlacedInside.currentHitpoints > bld.unitPlacedInside.getStaticData().maxHitpoints {
				bld.unitPlacedInside.currentHitpoints = bld.unitPlacedInside.getStaticData().maxHitpoints
			}
		} else {
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
	if u.turrets != nil && !u.turrets[0].canRotate() {
		if u.turrets[0].rotationDegree == u.currentAction.targetRotation {
			u.currentAction.code = ACTION_WAIT
		} else if u.turrets[0].targetActor == nil {
			u.turrets[0].rotationDegree += geometry.GetDiffForRotationStep(u.turrets[0].rotationDegree, u.currentAction.targetRotation, u.turrets[0].getStaticData().rotateSpeed)
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

func (b *battlefield) executeGroundMoveActionForUnit(u *unit) {
	x, y := u.centerX, u.centerY
	targetX, targetY := geometry.TileCoordsToPhysicalCoords(u.currentAction.targetTileX, u.currentAction.targetTileY)
	vx, vy := targetX-x, targetY-y
	if areFloatsRoughlyEqual(x, targetX) && areFloatsRoughlyEqual(y, targetY) {
		u.centerX = targetX
		u.centerY = targetY
		u.currentAction.reset()
		return
		// debugWritef("Tick %d: action finished\n", b.currentTick)
	}
	// rotate to target if needed
	if !geometry.IsVectorDegreeEqualTo(vx, vy, u.chassisDegree) {
		u.rotateChassisTowardsVector(vx, vy)
		return
	}

	displacementX, displacementY := vx, vy

	if math.Abs(vx) > u.getStaticData().movementSpeed {
		displacementX = u.getStaticData().movementSpeed * vx / math.Abs(vx)
	}
	if math.Abs(vy) > u.getStaticData().movementSpeed {
		displacementY = u.getStaticData().movementSpeed * vy / math.Abs(vy)
	}

	currTx, currTy := geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
	currTcx, currTcy := geometry.TileCoordsToPhysicalCoords(currTx, currTy)
	// if we're passing through tile center by our movement...
	if math.Signbit(currTcx-u.centerX) != math.Signbit(currTcx-u.centerX-displacementX) ||
		math.Signbit(currTcy-u.centerY) != math.Signbit(currTcy-u.centerY-displacementY) ||
		// ...or we're starting from this center...
		areFloatsRoughlyEqual(x, currTcx) && areFloatsRoughlyEqual(y, currTcy) {

		intVx, intVy := geometry.Float64VectorToIntDirectionVector(vx, vy)
		//debugWritef("Checking: %d, %d ", currTx+intVx, currTy+intVy)
		//debugWritef("Having %v,%v from cooords %v,%v \n", intVx, intVy, currTx, currTy)
		// ...then check if there is something on "next" tile.
		if !b.isTileClearToBeMovedInto(currTx+intVx, currTy+intVy, u) {
			// If so, stand by.
			u.centerX, u.centerY = currTcx, currTcy
			u.currentAction.reset()
			return
		}
	}

	u.centerX += displacementX
	u.centerY += displacementY
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
		if bld, ok := a.(*building); ok {
			if bld.getStaticData().buildType == BTYPE_PLACE_FIRST {
				if tBld, ok := a.(*building).currentAction.targetActor.(*building); ok {
					tBld.getCurrentAction().completionAmount = act.completionAmount
				}
			}
		}
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

func (b *battlefield) executeBeingBuiltActionForBuilding(bld *building) {
	if bld.getCurrentAction().builtAs != BTYPE_PLACE_FIRST {
		bld.currentAction.completionAmount++
	}
	if bld.currentAction.getCompletionPercent() == 100 {
		bld.currentAction.reset()
	}
}
