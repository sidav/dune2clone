package main

import "dune2clone/geometry"

func (b *battlefield) executeWaitOrderForAircraft(u *unit) {
	if !u.getStaticData().isAircraft {
		panic("Aircraft order assigned to non-aircraft")
	}
	i := u.faction.dispatchRequests.Front()
	if i != nil {
		dr := i.Value.(*dispatchRequestStruct)
		u.currentOrder.setTargetTileCoords(dr.targetTileX, dr.targetTileY)
		u.currentOrder.targetActor = dr.requester
		u.currentOrder.code = dr.assignedOrderCode
		u.faction.removeDispatchRequest(dr)
		return
	}
}

func (b *battlefield) executeCarryUnitOrderForAircraft(carrier *unit) {
	// carrierTx, carrierTy := geometry.TrueCoordsToTileCoords(carrier.centerX, carrier.centerY)
	// Order: pick targetActor up, then move it to target coords, then drop it down
	if carrier.carriedUnit == nil { // need to pick up
		targetUnit := carrier.currentOrder.targetActor.(*unit)
		targetX, targetY := geometry.TrueCoordsToTileCoords(targetUnit.centerX, targetUnit.centerY)
		if geometry.GetApproxDistFromTo(targetX, targetY, carrier.currentOrder.targetTileX, carrier.currentOrder.targetTileY) < 5 {
			// target is too close already, no need for transport
			carrier.currentOrder.resetOrder()
			return
		}
		debugWrite("PICK ORDER SET")
		carrier.currentAction.code = ACTION_AIR_PICK_UNIT_UP
		carrier.currentAction.targetActor = targetUnit
	} else { // already picked up
		if carrier.isPresentAt(carrier.currentOrder.targetTileX, carrier.currentOrder.targetTileY) {
			debugWrite("DROP ACTION SET")
			carrier.currentAction.code = ACTION_AIR_DROP_UNIT
			carrier.currentAction.setTargetTileCoords(carrier.currentOrder.targetTileX, carrier.currentOrder.targetTileY)
			carrier.currentOrder.resetOrder()
		} else {
			debugWrite("APPROACH ACTION SET")
			carrier.currentAction.code = ACTION_AIR_APPROACH_LAND_TILE
			carrier.currentAction.setTargetTileCoords(carrier.currentOrder.targetTileX, carrier.currentOrder.targetTileY)
		}
	}
}
