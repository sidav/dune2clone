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

	default:
		panic("No action execution func!")
	}
}

func (b *battlefield) executeWaitActionForUnit(u *unit) {
	// rotate the unit to target if unit's turret can't rotate itself
	if u.turret != nil && !u.turret.canRotate() && u.turret.targetActor != nil {
		x, y := u.turret.targetActor.getPhysicalCenterCoords()
		x -= u.centerX
		y -= u.centerY
		u.rotateChassisTowardsVector(x, y)
	}
}

func (b *battlefield) executeWaitActionForBuilding(bld *building) {
	if bld.getStaticData().receivesResources && bld.unitPlacedInside != nil {
		if bld.unitPlacedInside.currentCargoAmount > 0 {
			const HARVEST_PER_TICK = 11
			received := min(HARVEST_PER_TICK, bld.unitPlacedInside.currentCargoAmount) // TODO: replace
			bld.getFaction().money += float64(received)
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
	if u.turret.canRotate() {
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
	tx, ty := geometry.TileCoordsToPhysicalCoords(u.currentAction.targetTileX, u.currentAction.targetTileY)

	if areFloatsAlmostEqual(x, tx) && areFloatsAlmostEqual(y, ty) {
		u.centerX = tx
		u.centerY = ty
		u.currentAction.reset()
		// debugWritef("Tick %d: action finished\n", b.currentTick)
	}

	vx, vy := tx-x, ty-y

	if !geometry.IsVectorDegreeEqualTo(vx, vy, u.chassisDegree) {
		u.rotateChassisTowardsVector(vx, vy)
		return
	}

	if math.Abs(vx) < u.getStaticData().movementSpeed {
		u.centerX = tx // source of possible movement lag :(
	} else {
		u.centerX += u.getStaticData().movementSpeed * vx / math.Abs(vx)
	}

	if math.Abs(vy) < u.getStaticData().movementSpeed {
		u.centerY = ty // source of possible movement lag :(
	} else {
		u.centerY += u.getStaticData().movementSpeed * vy / math.Abs(vy)
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
	if act.getCompletionPercent() < 100 && a.getFaction().money > moneySpent {
		a.getFaction().money -= moneySpent
		act.completionAmount++
	}
	// if it was a unit, place it right away
	if unt, ok := act.targetActor.(*unit); ok && act.getCompletionPercent() >= 100 {
		if bld, ok := a.(*building); ok {
			for x := bld.topLeftX; x < bld.topLeftX+bld.getStaticData().w; x++ {
				// for y := bld.topLeftY-1; y <= bld.topLeftY+bld.getStaticData().h; y++ {
				y := bld.topLeftY + bld.getStaticData().h
				if b.costMapForMovement(x, y) != -1 {
					unt.centerX, unt.centerY = geometry.TileCoordsToPhysicalCoords(x, y)
					// debugWritef("+%v", unt)
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
