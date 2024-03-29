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
			if u.getStaticData().IsAircraft {
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
		if u, ok := a.(*unit); ok && u.getStaticData().IsAircraft {
			b.executeAirApproachLandTileActionForUnit(u)
		}
	case ACTION_AIR_APPROACH_ACTOR:
		if u, ok := a.(*unit); ok && u.getStaticData().IsAircraft {
			b.executeAirApproachTargetActorActionForUnit(u)
		}
	case ACTION_AIR_PICK_UNIT_UP:
		if u, ok := a.(*unit); ok && u.getStaticData().IsAircraft {
			b.executeAirPickUnitUpActionForUnit(u)
		}
	case ACTION_AIR_DROP_UNIT:
		if u, ok := a.(*unit); ok && u.getStaticData().IsAircraft {
			b.executeAirDropActionForUnit(u)
		}
	case ACTION_DEPLOY:
		if u, ok := a.(*unit); ok {
			b.executeBeingDeployedActionForUnit(u)
		}

	default:
		panic("No action execution func!")
	}
}

func (b *battlefield) executeBuildingSelfRepair(bld *building) {
	cost := float64(bld.getStaticData().Cost/2) / float64(bld.getStaticData().MaxHitpoints)
	if bld.currentHitpoints >= bld.getMaxHitpoints() || bld.faction.getMoney() < cost {
		bld.isRepairingSelf = false
		return
	}
	bld.currentHitpoints++
	bld.faction.spendMoney(cost)
}

func (b *battlefield) executeWaitActionForUnit(u *unit) {
	if u.getStaticData().IsAircraft {
		if u.currentAction.targetTileX == -1 {
			u.currentAction.targetTileX, u.currentAction.targetTileY = geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
		}
		targetX, targetY := geometry.TileCoordsToTrueCoords(u.currentAction.targetTileX, u.currentAction.targetTileY)
		vx, vy := geometry.DegreeToUnitVector(u.chassisDegree)
		speed := u.getMovementSpeed()
		u.setPhysicalCenterCoords(
			u.centerX+speed*vx,
			u.centerY+speed*vy,
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
		// idle rotating
		if u.turrets != nil && u.turrets[0].canRotate() && u.turrets[0].targetActor == nil && (b.currentTick/config.Engine.UnitsActionPeriod)%25 == 0 {
			rotSpeed := u.turrets[0].getStaticData().RotateSpeed
			if abs(geometry.GetDiffForRotationStep(u.chassisDegree, u.turrets[0].rotationDegree, 90)) >= 10 {
				u.turrets[0].rotationDegree += geometry.GetDiffForRotationStep(u.turrets[0].rotationDegree, u.chassisDegree, rotSpeed)
			} else if rnd.OneChanceFrom(15) {
				// here be glitches
				u.turrets[0].rotationDegree += rnd.RandInRange(-2, 2) * 25
			}
			u.turrets[0].normalizeDegrees()
		}
	}
}

func (b *battlefield) executeWaitActionForBuilding(bld *building) {
	if bld.getStaticData().ReceivesResources && bld.unitPlacedInside != nil {
		if bld.unitPlacedInside.currentCargoAmount > 0 {
			// Receive resources
			received := min(config.Economy.HarvesterUnloadSpeed, bld.unitPlacedInside.currentCargoAmount)
			if bld.faction.getStorageRemaining() > 0 {
				received = min(received, int(bld.faction.getStorageRemaining()))
			} else {
				return
			}
			bld.getFaction().receiveResources(float64(received), false)
			bld.unitPlacedInside.currentCargoAmount -= received
		} else {
			b.addActor(bld.unitPlacedInside)
			bld.unitPlacedInside = nil
		}
	}
	if bld.getStaticData().RepairsUnits && bld.unitPlacedInside != nil {
		if bld.unitPlacedInside.getHitpointsPercentage() < 100 {
			const REPAIR_PER_TICK = 2
			bld.unitPlacedInside.receiveHealing(REPAIR_PER_TICK)
		} else {
			// force clear tile so that the entry point will be cleared for airplanes too
			upx, upy := bld.getUnitPlacementAbsoluteCoords()
			b.clearTilesOccupationInRect(upx, upy, 1, 1)
			b.addActor(bld.unitPlacedInside)
			bld.unitPlacedInside = nil
		}
	}
	if bld.turret != nil && bld.turret.targetActor == nil && bld.getFaction().getAvailableEnergy() >= 0 && rnd.OneChanceFrom(50) {
		bld.turret.rotationDegree += rnd.RandInRange(-1, 1) * 45
		bld.turret.normalizeDegrees()
	}
}

func (b *battlefield) executeEnterBuildingActionForUnit(u *unit) {
	if u.currentAction.targetActor.(*building).unitPlacedInside == nil {
		u.currentAction.code = ACTION_WAIT
		ptx, pty := geometry.TileCoordsToTrueCoords(u.currentAction.targetActor.(*building).getUnitPlacementAbsoluteCoords())
		u.setPhysicalCenterCoords(ptx, pty)
		u.chassisDegree = 90
		for i := range u.turrets {
			u.turrets[i].rotationDegree = 90
		}
		ux, uy := u.getTileCoords()
		b.removeActor(u)
		b.setTilesOccupiedByActor(ux, uy, 1, 1, u) // WORKAROUND so that the entry space will be occupied
		u.currentAction.targetActor.(*building).unitPlacedInside = u
	}
}

func (b *battlefield) executeHarvestActionForActor(a actor) {
	if u, ok := a.(*unit); ok {
		x, y := u.getPhysicalCenterCoords()
		utx, uty := geometry.TrueCoordsToTileCoords(x, y)
		if u.currentCargoAmount < u.getStaticData().MaxCargoAmount && b.tiles[utx][uty].resourcesAmount > 0 {
			harvestedAmount := min(config.Economy.HarvestingSpeed, b.tiles[utx][uty].resourcesAmount)
			harvestedAmount = min(harvestedAmount, u.getStaticData().MaxCargoAmount-u.currentCargoAmount)
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
			u.turrets[0].rotationDegree += geometry.GetDiffForRotationStep(u.turrets[0].rotationDegree, u.currentAction.targetRotation, u.turrets[0].getStaticData().RotateSpeed)
			u.normalizeDegrees()
		}
	} else {
		if u.chassisDegree == u.currentAction.targetRotation {
			u.currentAction.code = ACTION_WAIT
		} else {
			u.chassisDegree += geometry.GetDiffForRotationStep(u.chassisDegree, u.currentAction.targetRotation, u.getStaticData().ChassisRotationSpeed)
			u.normalizeDegrees()
		}
	}
}

func (b *battlefield) executeGroundMoveActionForUnit(u *unit) {
	x, y := u.centerX, u.centerY
	currTx, currTy := geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
	targetX, targetY := geometry.TileCoordsToTrueCoords(u.currentAction.targetTileX, u.currentAction.targetTileY)
	// interrupting movement, if another order was given
	if u.currentAction.interruptMovement {
		// setting target to the nearest tile
		intVx, intVy := geometry.Float64VectorToIntUnitVector(targetX-x, targetY-y)
		u.currentAction.targetTileX, u.currentAction.targetTileY = currTx+intVx, currTy+intVy
		targetX, targetY = geometry.TileCoordsToTrueCoords(u.currentAction.targetTileX, u.currentAction.targetTileY)
		u.currentAction.interruptMovement = false
	}
	vx, vy := targetX-x, targetY-y
	if areFloatsRoughlyEqual(x, targetX) && areFloatsRoughlyEqual(y, targetY) {
		u.centerX = targetX
		u.centerY = targetY
		u.currentAction.resetAction()
		return
		// debugWritef("Tick %d: action finished\n", b.currentTick)
	}
	// rotate to target if needed
	if !geometry.IsVectorDegreeEqualTo(vx, vy, u.chassisDegree) {
		u.rotateChassisTowardsVector(vx, vy)
		return
	}

	displacementX, displacementY := vx, vy
	speed := u.getMovementSpeed()

	if math.Abs(vx) > speed {
		displacementX = speed * vx / math.Abs(vx)
	}
	if math.Abs(vy) > speed {
		displacementY = speed * vy / math.Abs(vy)
	}

	currTcx, currTcy := geometry.TileCoordsToTrueCoords(currTx, currTy)
	// if we're passing through tile center by our movement...
	if math.Signbit(currTcx-u.centerX) != math.Signbit(currTcx-u.centerX-displacementX) ||
		math.Signbit(currTcy-u.centerY) != math.Signbit(currTcy-u.centerY-displacementY) ||
		// ...or we're starting from this center...
		areFloatsRoughlyEqual(x, currTcx) && areFloatsRoughlyEqual(y, currTcy) {

		intVx, intVy := geometry.Float64VectorToIntUnitVector(vx, vy)
		//debugWritef("Checking: %d, %d ", currTx+intVx, currTy+intVy)
		//debugWritef("Having %v,%v from cooords %v,%v \n", intVx, intVy, currTx, currTy)
		// ...then check if there is something on "next" tile.
		nextTx, nextTy := currTx+intVx, currTy+intVy
		if !b.isTileClearToBeMovedInto(nextTx, nextTy, u) {
			// If so, stand by, and increase action failure counter
			u.centerX, u.centerY = currTcx, currTcy
			u.currentAction.fail(true) // action may be reset if failed continuously
			return

			// additional check, so that the next tile won't be occupied if no further inter-tile movement is needed
		} else if currTx != u.currentAction.targetTileX || currTy != u.currentAction.targetTileY {
			b.tiles[currTx][currTy].isOccupiedByActor = nil
			b.setTilesOccupiedByActor(nextTx, nextTy, 1, 1, u)
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
		moneySpent = float64(bld.getStaticData().Cost) /
			float64(bld.getStaticData().BuildTime*(config.Engine.TicksPerNominalSecond/config.Engine.BuildingsActionPeriod))
	}
	if unt, ok := act.targetActor.(*unit); ok {
		moneySpent = float64(unt.getStaticData().Cost) /
			float64(unt.getStaticData().BuildTime*(config.Engine.TicksPerNominalSecond/config.Engine.BuildingsActionPeriod))
	}
	// initial percent
	initalPercent := act.getCompletionPercent()
	// spend money
	coeff := a.getFaction().getEnergyProductionMultiplier()
	moneySpent *= coeff
	if act.getCompletionPercent() < 100 && a.getFaction().getMoney() > moneySpent {
		a.getFaction().spendMoney(moneySpent)
		act.completionAmount += coeff
		if bld, ok := a.(*building); ok {
			if bld.getStaticData().BuildType == BTYPE_PLACE_FIRST {
				if tBld, ok := a.(*building).currentAction.targetActor.(*building); ok {
					tBld.getCurrentAction().completionAmount = act.completionAmount
					// increase bld curr hitpoints
					// > 0 because the building has already received starting HP at 0%
					if initalPercent > 0 && act.getCompletionPercent() > initalPercent {
						diff := act.getCompletionPercent() - initalPercent
						tBld.currentHitpoints += diff * tBld.getStaticData().MaxHitpoints / 100
					}
				}
			}
		}
	}
	// if it was a unit, place it right away
	if unt, ok := act.targetActor.(*unit); ok && act.getCompletionPercent() >= 100 {
		if bld, ok := a.(*building); ok {
			for x := bld.topLeftX; x < bld.topLeftX+bld.getStaticData().W; x++ {
				// for y := bld.topLeftY-1; y <= bld.topLeftY+bld.getStaticData().h; y++ {
				y := bld.topLeftY + bld.getStaticData().H
				if unt.getStaticData().IsAircraft || b.costMapForMovement(x, y) != -1 {
					unt.centerX, unt.centerY = geometry.TileCoordsToTrueCoords(x, y)
					if bld.rallyTileX != -1 && unt.currentOrder.code == ORDER_NONE {
						unt.currentOrder.code = ORDER_MOVE
						unt.currentOrder.setTargetTileCoords(bld.rallyTileX, bld.rallytileY)
						if !unt.getStaticData().IsAircraft {
							bld.faction.addDispatchRequest(unt, bld.rallyTileX, bld.rallytileY, ORDER_CARRY_UNIT_TO_TARGET_COORDS, b.currentTick+100)
						}
					}
					b.addActor(unt)
					bld.currentAction.resetAction()
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
	if bld.currentAction.getCompletionPercent() >= 100 {
		bld.currentAction.resetAction()
	}
}

func (b *battlefield) executeBeingDeployedActionForUnit(u *unit) {
	if u.currentAction.completionAmount == u.currentAction.maxCompletionAmount {
		tx, ty := geometry.TrueCoordsToTileCoords(u.getPhysicalCenterCoords())
		if b.canUnitBeDeployedAt(u, tx, ty) {
			targetBld := u.currentAction.targetActor.(*building)
			targetBld.topLeftX = tx
			targetBld.topLeftY = ty
			targetBld.currentAction.code = ACTION_BEING_BUILT
			targetBld.currentAction.builtAs = BTYPE_BUILD_FIRST
			targetBld.currentAction.maxCompletionAmount = float64(config.Engine.BuildingAnimationTicks)
			b.removeActor(u)
			b.addActor(targetBld)
			return
		}
		u.currentAction.resetAction()
		return
	}
	u.currentAction.completionAmount++
}

func (b *battlefield) canUnitsActionBeInterrupted(u *unit) bool {
	act := u.getCurrentAction()
	if u.isInAir() {
		return true
	}
	switch act.code {
	case ACTION_WAIT, ACTION_ROTATE, ACTION_HARVEST:
		return true
	case ACTION_MOVE:
		return act.failedContinuously
	default:
		return false
	}
}
