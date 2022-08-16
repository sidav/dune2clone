package main

import "dune2clone/geometry"

func (b *battlefield) executeWaitOrderForAircraft(u *unit) {
	if !u.getStaticData().isAircraft {
		panic("Aircraft order assigned to non-aircraft")
	}
	i := u.faction.dispatchRequests.Front()
	if i != nil {
		dr := i.Value.(*dispatchRequestStruct)
		requesterTileX, requesterTileY := geometry.TrueCoordsToTileCoords(dr.requester.getPhysicalCenterCoords())
		if geometry.GetApproxDistFromTo(dr.targetTileX, dr.targetTileY, requesterTileX, requesterTileY) > 5 {
			u.currentOrder.setTargetTileCoords(dr.targetTileX, dr.targetTileY)
			u.currentOrder.targetActor = dr.requester
			u.currentOrder.code = dr.assignedOrderCode
		}
		u.faction.removeDispatchRequest(dr)
		return
	}
}

func (b *battlefield) executeCarryUnitOrderForAircraft(carrier *unit) {
	carrierTx, carrierTy := geometry.TrueCoordsToTileCoords(carrier.centerX, carrier.centerY)
	// Order: pick targetActor up, then move it to target coords, then drop it down
	if carrier.carriedUnit == nil { // need to pick up
		targetUnit := carrier.currentOrder.targetActor.(*unit)
		targetX, targetY := geometry.TrueCoordsToTileCoords(targetUnit.centerX, targetUnit.centerY)
		if carrierTx == targetX && carrierTy == targetY { // we're over target unit
			carrier.centerX, carrier.centerY = targetUnit.getPhysicalCenterCoords()
			if carrier.chassisDegree == targetUnit.chassisDegree { // pick up
				// TODO: make this an action
				targetUnit.currentAction.reset()
				carrier.carriedUnit = targetUnit
				b.removeActor(targetUnit)
				debugWrite("PICKED UP")
			} else { // rotate properly
				debugWrite("ROTATING")
				carrier.currentAction.code = ACTION_ROTATE
				carrier.currentAction.targetRotation = targetUnit.chassisDegree
			}
		} else { // need to move to target actor
			debugWrite("MOVING OUT")
			carrier.currentAction.code = ACTION_MOVE
			carrier.currentAction.setTargetTileCoords(targetX, targetY)
		}
	} else {
		if geometry.GetApproxDistFromTo(carrierTx, carrierTy, carrier.currentOrder.targetTileX, carrier.currentOrder.targetTileY) <= 1 {
			// TODO: make this an action
			debugWrite("DROPPING 1...")
			if b.isTileClearToBeMovedInto(carrierTx, carrierTy, carrier.carriedUnit) {
				debugWrite("DROPPING 2...")
				carrier.carriedUnit.centerX, carrier.carriedUnit.centerY = geometry.TileCoordsToPhysicalCoords(carrierTx, carrierTy)
				b.addActor(carrier.carriedUnit)
				carrier.carriedUnit = nil
				carrier.currentAction.reset()
				carrier.currentOrder.resetOrder()
				debugWrite("")
			}
		} else { // need to move to target coords
			debugWrite("MOVING TO DROP")
			carrier.currentAction.code = ACTION_MOVE
			carrier.currentAction.setTargetTileCoords(carrier.currentOrder.targetTileX, carrier.currentOrder.targetTileY)
		}
	}
}
