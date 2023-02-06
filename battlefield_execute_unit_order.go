package main

import (
	"dune2clone/geometry"
	"math"
)

func (b *battlefield) executeOrderForUnit(u *unit) {
	if u.currentOrder.isNew {
		if u.currentAction.code == ACTION_MOVE {
			u.currentAction.interruptMovement = true
		}
		u.currentOrder.isNew = false
	}
	if !b.canUnitsActionBeInterrupted(u) {
		return
	}
	if u.currentOrder.code == ORDER_NONE {
		if u.getStaticData().IsAircraft {
			b.executeWaitOrderForAircraft(u)
		}
		return // TODO: look around for targets etc
	}

	switch u.currentOrder.code {
	case ORDER_MOVE:
		b.executeMoveOrder(u)
	case ORDER_ATTACK:
		b.executeAttackOrder(u)
	case ORDER_HARVEST:
		b.executeHarvestOrder(u)
	case ORDER_RETURN_TO_REFINERY:
		b.executeReturnResourcesOrder(u)
	case ORDER_MOVE_TO_REPAIR:
		if u.getStaticData().IsAircraft {
			b.executeAirMoveToRepairOrder(u)
		} else {
			b.executeGroundMoveToRepairOrder(u)
		}
	case ORDER_DEPLOY:
		b.executeDeployOrderForUnit(u)

	case ORDER_CARRY_UNIT_TO_TARGET_COORDS:
		b.executeCarryUnitOrderForAircraft(u)
	default:
		panic("No order execution func!")
	}
}

func (b *battlefield) executeMoveOrder(u *unit) {
	// x, y := u.centerX, u.centerY
	utx, uty := geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
	orderTileX, orderTileY := u.currentOrder.targetTileX, u.currentOrder.targetTileY
	if utx == orderTileX && uty == orderTileY || u.currentOrder.failedToExecute {
		u.currentOrder.resetOrder()
		return
	}
	if u.getStaticData().IsAircraft {
		u.currentAction.code = ACTION_MOVE
		u.currentAction.targetTileX, u.currentAction.targetTileY = orderTileX, orderTileY
	} else {
		b.SetActionForUnitForPathTo(u, orderTileX, orderTileY)
	}
}

func (b *battlefield) executeAttackOrder(u *unit) {
	if !u.currentOrder.targetActor.isAlive() {
		u.currentOrder.resetOrder()
		return
	}
	if len(u.turrets) != 0 {
		u.turrets[0].targetActor = u.currentOrder.targetActor
	}
	utx, uty := geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
	orderTileX, orderTileY := geometry.TrueCoordsToTileCoords(u.currentOrder.targetActor.getPhysicalCenterCoords())
	if geometry.GetApproxDistFromTo(utx, uty, orderTileX, orderTileY) <= u.getMainTurretRange() {
		return
	} else {
		if u.getStaticData().IsAircraft {
			u.currentAction.code = ACTION_MOVE
			u.currentAction.targetTileX, u.currentAction.targetTileY = orderTileX, orderTileY
		} else {
			// find the closest tile, from which the target can be reached
			orderTileX, orderTileY := geometry.SpiralSearchForLowestScoreFrom(
				func(x, y int) int { return geometry.GetApproxDistFromTo(x, y, utx, uty) },
				func(x, y int) bool {
					return b.areTileCoordsValid(x, y) &&
						geometry.GetApproxDistFromTo(x, y, orderTileX, orderTileY) <= u.getMainTurretRange()
				},
				orderTileX, orderTileY, u.getMainTurretRange(), b.tickToNonImportantRandom(4),
			)
			b.SetActionForUnitForPathTo(u, orderTileX, orderTileY)
		}
	}
}

func (b *battlefield) executeHarvestOrder(u *unit) {
	if u.currentCargoAmount >= u.getStaticData().MaxCargoAmount || u.currentOrder.failedToExecute {
		u.currentOrder.dispatchCalled = false
		u.currentOrder.failedToExecute = false
		u.currentOrder.code = ORDER_RETURN_TO_REFINERY
	}

	utx, uty := geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
	orderTileX, orderTileY := u.currentOrder.targetTileX, u.currentOrder.targetTileY
	if orderTileX == -1 || orderTileY == -1 {
		u.currentOrder.targetTileX, u.currentOrder.targetTileY = utx, uty
		orderTileX, orderTileY = utx, uty
	}

	if b.tiles[orderTileX][orderTileY].resourcesAmount == 0 || !b.isTileClearToBeMovedInto(orderTileX, orderTileY, u) {
		// find resources coords
		rx, ry := b.getCoordsOfClosestEmptyTileWithResourcesTo(orderTileX, orderTileY)
		if rx < 0 || ry < 0 {
			if u.currentCargoAmount > 0 {
				u.currentOrder.code = ORDER_RETURN_TO_REFINERY
			} else {
				u.currentOrder.code = ORDER_MOVE // TODO: remove since resources now grow?
			}
			return
		}
		u.currentOrder.targetTileX, u.currentOrder.targetTileY = rx, ry
		orderTileX, orderTileY = rx, ry
	}

	if orderTileX != utx || orderTileY != uty {
		b.SetActionForUnitForPathTo(u, orderTileX, orderTileY)
	}

	if utx == orderTileX && uty == orderTileY {
		u.currentAction.code = ACTION_HARVEST
	} else if !u.currentOrder.dispatchCalled {
		u.faction.addDispatchRequest(u, orderTileX, orderTileY, ORDER_CARRY_UNIT_TO_TARGET_COORDS, b.currentTick+100)
		u.currentOrder.dispatchCalled = true
	}
}

func (b *battlefield) executeReturnResourcesOrder(u *unit) {
	// x, y := u.centerX, u.centerY
	utx, uty := geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
	// for this, target tile is resource tile. Target refinery is in targetActor.
	if (b.currentTick/30)%2 == 0 || u.currentOrder.targetActor == nil ||
		u.currentOrder.targetActor.(*building).unitPlacedInside != nil {

		u.currentOrder.targetActor = b.getClosestEmptyFactionRefineryFromCoords(u.faction, utx, uty)
	}
	if u.currentOrder.targetActor == nil {
		// nothing found, doing nothing
		return
	}
	orderTileX, orderTileY := u.currentOrder.targetActor.(*building).getUnitPlacementAbsoluteCoords()
	if !u.currentOrder.dispatchCalled {
		u.faction.addDispatchRequest(u, orderTileX, orderTileY, ORDER_CARRY_UNIT_TO_TARGET_COORDS, b.currentTick+100)
		u.currentOrder.dispatchCalled = true
	}

	if orderTileX == utx && orderTileY == uty {
		u.currentOrder.code = ORDER_HARVEST
		u.currentOrder.dispatchCalled = false
		u.currentAction.code = ACTION_ENTER_BUILDING
		u.currentAction.targetActor = u.currentOrder.targetActor
		return
	}

	b.SetActionForUnitForPathTo(u, orderTileX, orderTileY)
}

func (b *battlefield) executeGroundMoveToRepairOrder(u *unit) {
	// x, y := u.centerX, u.centerY
	utx, uty := geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
	// for this, target tile is tile to return to after repairs.
	if u.currentOrder.targetActor == nil || !u.currentOrder.targetActor.isAlive() {
		depot := b.getClosestEmptyFactionRepairDepotFromCoords(u.faction, utx, uty)
		if depot == nil {
			u.currentOrder.resetOrder()
			return
		}
		u.currentOrder.targetActor = depot
	}
	orderTileX, orderTileY := u.currentOrder.targetActor.(*building).getUnitPlacementAbsoluteCoords()
	if !u.currentOrder.dispatchCalled {
		u.faction.addDispatchRequest(u, orderTileX, orderTileY, ORDER_CARRY_UNIT_TO_TARGET_COORDS, b.currentTick+100)
		u.currentOrder.dispatchCalled = true
	}

	if orderTileX == utx && orderTileY == uty {
		u.currentOrder.code = ORDER_MOVE
		u.currentOrder.setTargetTileCoords(orderTileX, orderTileY+1)
		u.currentOrder.dispatchCalled = false
		u.currentAction.code = ACTION_ENTER_BUILDING
		u.currentAction.targetActor = u.currentOrder.targetActor
		return
	}

	b.SetActionForUnitForPathTo(u, orderTileX, orderTileY)
}

func (b *battlefield) executeDeployOrderForUnit(u *unit) {
	if u.getStaticData().CanBeDeployed {
		u.currentAction.code = ACTION_DEPLOY
		u.currentAction.targetActor = createBuilding(u.getStaticData().DeploysInto, 0, 0, u.getFaction())
		u.currentAction.maxCompletionAmount = 2 * (DESIRED_TPS / UNIT_ACTIONS_TICK_EACH)
	}
}

func (b *battlefield) SetActionForUnitForPathTo(u *unit, tx, ty int) {
	utx, uty := geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)

	path := b.findPathForUnitTo(u, tx, ty, false)
	vx, vy := path.GetNextStepVector()
	nowhereToMove := vx == 0 && vy == 0

	// creating BIG move action for several same-vector path cells
	currPathChild := path.Child
	multiplier := 1
	for currPathChild != nil {
		vxNext, vyNext := currPathChild.GetNextStepVector()
		if vx == vxNext && vy == vyNext {
			multiplier++
		} else {
			break
		}
		currPathChild = currPathChild.Child
	}
	vx *= multiplier
	vy *= multiplier
	ttx, tty := utx+vx, uty+vy

	if nowhereToMove || u.currentAction.targetTileX == ttx && u.currentAction.targetTileY == tty {
		if nowhereToMove || u.currentAction.failedContinuously {
			u.currentAction.resetAction()
			u.currentOrder.failedToExecute = true
			return
		}
	}
	u.currentAction.code = ACTION_MOVE
	u.currentAction.failedContinuously = false
	u.currentAction.setTargetTileCoords(ttx, tty)
}

func (b *battlefield) getClosestEmptyFactionRefineryFromCoords(f *faction, x, y int) actor {
	var selected actor = nil
	closestDist := math.MaxInt64
	for i := b.buildings.Front(); i != nil; i = i.Next() {
		bld := i.Value.(*building)
		if bld.faction != f || !bld.getStaticData().receivesResources || bld.unitPlacedInside != nil {
			continue
		}
		bldCX, bldCY := bld.getUnitPlacementAbsoluteCoords()
		distFromBld := geometry.GetApproxDistFromTo(x, y, bldCX, bldCY)
		if selected == nil || distFromBld < closestDist {
			closestDist = distFromBld
			selected = bld
		}
	}
	return selected
}

func (b *battlefield) getClosestEmptyFactionRepairDepotFromCoords(f *faction, x, y int) actor {
	var selected actor = nil
	closestDist := math.MaxInt64
	for i := b.buildings.Front(); i != nil; i = i.Next() {
		bld := i.Value.(*building)
		if bld.faction != f || !bld.getStaticData().repairsUnits || bld.unitPlacedInside != nil {
			continue
		}
		bldCX, bldCY := bld.getUnitPlacementAbsoluteCoords()
		distFromBld := geometry.GetApproxDistFromTo(x, y, bldCX, bldCY)
		if selected == nil || distFromBld < closestDist {
			closestDist = distFromBld
			selected = bld
		}
	}
	return selected
}
