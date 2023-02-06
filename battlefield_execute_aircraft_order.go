package main

import "dune2clone/geometry"

func (b *battlefield) executeWaitOrderForAircraft(u *unit) {
	if !u.getStaticData().IsAircraft {
		panic("Aircraft order assigned to non-aircraft")
	}
	// look for dispatch requests
	i := u.faction.dispatchRequests.Front()
	if i != nil {
		dr := i.Value.(*dispatchRequestStruct)
		if dr.assignedOrderCode == ORDER_CARRY_UNIT_TO_TARGET_COORDS && u.getStaticData().IsTransport {
			u.currentOrder.setTargetTileCoords(dr.targetTileX, dr.targetTileY)
			u.currentOrder.targetActor = dr.requester
			u.currentOrder.code = dr.assignedOrderCode
			u.faction.removeDispatchRequest(dr)
			return
		}
	}
}

func (b *battlefield) executeAirMoveToRepairOrder(u *unit) {
	// x, y := u.centerX, u.centerY
	utx, uty := geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
	// for this, target tile is tile to return to after repairs.
	if u.currentOrder.targetActor == nil {
		// nothing found, doing nothing
		return
	}
	orderTileX, orderTileY := u.currentOrder.targetActor.(*building).getUnitPlacementAbsoluteCoords()

	if orderTileX == utx && orderTileY == uty {
		u.currentOrder.code = ORDER_MOVE
		u.currentOrder.setTargetTileCoords(orderTileX, orderTileY+1)
		u.currentOrder.dispatchCalled = false
		u.currentAction.code = ACTION_ENTER_BUILDING
		u.currentAction.targetActor = u.currentOrder.targetActor
		return
	} else {
		u.currentAction.code = ACTION_AIR_APPROACH_LAND_TILE
		u.currentAction.setTargetTileCoords(orderTileX, orderTileY)
	}
}

func (b *battlefield) executeCarryUnitOrderForAircraft(carrier *unit) {
	// Order: pick targetActor up, then move it to target coords, then drop it down
	carrierTx, carrierTy := carrier.getTileCoords()
	targetUnit := carrier.currentOrder.targetActor.(*unit)
	targetX, targetY := targetUnit.getTileCoords()
	if carrier.carriedUnit == nil { // need to pick up
		distOfPickableToItsTarget := geometry.GetApproxDistFromTo(targetX, targetY, carrier.currentOrder.targetTileX, carrier.currentOrder.targetTileY)
		if distOfPickableToItsTarget <= 5 {
			// target is too close already, no need for transport
			carrier.currentOrder.resetOrder()
			return
		}
		const rangeToLockOn = 2
		if geometry.GetApproxDistFromTo(carrierTx, carrierTy, targetX, targetY) > rangeToLockOn {
			// debugWrite("FLYING TO")
			carrier.currentAction.setTargetTileCoords(targetX, targetY)
			carrier.currentAction.targetActor = targetUnit
			carrier.currentAction.code = ACTION_AIR_APPROACH_ACTOR
			return
		}
		// debugWrite("PICK ORDER SET")
		carrier.currentAction.code = ACTION_AIR_PICK_UNIT_UP
		carrier.currentAction.targetActor = targetUnit
	} else { // already picked up
		if carrier.isPresentAt(carrier.currentOrder.targetTileX, carrier.currentOrder.targetTileY) {
			if b.isTileClearToBeMovedInto(carrier.currentOrder.targetTileX, carrier.currentOrder.targetTileY, carrier.carriedUnit) {

				// debugWrite("DROP ACTION SET")
				carrier.currentAction.code = ACTION_AIR_DROP_UNIT
				carrier.currentAction.setTargetTileCoords(carrier.currentOrder.targetTileX, carrier.currentOrder.targetTileY)
				carrier.currentOrder.resetOrder()
			} else {
				closestX, closestY := geometry.SpiralSearchForClosestConditionFrom(
					func(x, y int) bool { return b.isTileClearToBeMovedInto(x, y, carrier.carriedUnit) },
					carrier.currentOrder.targetTileX, carrier.currentOrder.targetTileY, 5, b.tickToNonImportantRandom(4),
				)
				carrier.currentOrder.setTargetTileCoords(closestX, closestY)
			}
		} else {
			// debugWrite("APPROACH ACTION SET")
			carrier.currentAction.code = ACTION_AIR_APPROACH_LAND_TILE
			carrier.currentAction.setTargetTileCoords(carrier.currentOrder.targetTileX, carrier.currentOrder.targetTileY)
		}
	}
}
