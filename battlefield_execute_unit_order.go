package main

import "dune2clone/geometry"

func (b *battlefield) executeOrderForUnit(u *unit) {
	if u.currentAction.code != ACTION_WAIT {
		return // execute the order only after finishing the action
	}
	if u.currentOrder.code == ORDER_NONE {
		return // TODO: look around for targets etc
	}

	switch u.currentOrder.code {
	case ORDER_MOVE:
		b.executeMoveOrder(u)
	case ORDER_HARVEST:
		b.executeHarvestOrder(u)
	case ORDER_RETURN_TO_REFINERY:
		b.executeReturnResourcesOrder(u)
	}
}

func (b *battlefield) executeMoveOrder(u *unit) {
	// x, y := u.centerX, u.centerY
	utx, uty := geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
	orderTileX, orderTileY := u.currentOrder.targetTileX, u.currentOrder.targetTileY
	if utx == orderTileX && uty == orderTileY {
		u.currentOrder.code = ORDER_NONE
		return
	}
	if u.getStaticData().isAircraft {
		u.currentAction.code = ACTION_MOVE
		u.currentAction.targetTileX, u.currentAction.targetTileY = orderTileX, orderTileY
	} else {
		b.SetActionForUnitForPathTo(u, orderTileX, orderTileY)
	}
}

func (b *battlefield) executeHarvestOrder(u *unit) {
	if u.currentCargoAmount >= u.getStaticData().maxCargoAmount {
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
				u.currentOrder.code = ORDER_MOVE
			}
			return
		}
		u.currentOrder.targetTileX, u.currentOrder.targetTileY = rx, ry
		orderTileX, orderTileY = rx, ry
	}

	b.SetActionForUnitForPathTo(u, orderTileX, orderTileY)

	if utx == orderTileX && uty == orderTileY {
		u.currentAction.code = ACTION_HARVEST
	}
}

func (b *battlefield) executeReturnResourcesOrder(u *unit) {
	// x, y := u.centerX, u.centerY
	utx, uty := geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
	// for this, target tile is resource tile. Target refinery is in targetActor.
	if (b.currentTick/30)%2 == 0 || u.currentOrder.targetActor == nil ||
		u.currentOrder.targetActor.(*building).unitPlacedInside != nil {

		u.currentOrder.targetActor = b.getClosestEmptyFactionRefineryFromCoords(u.faction, u.centerX, u.centerY)
	}
	if u.currentOrder.targetActor == nil {
		// nothing found, doing nothing
		return
	}
	orderTileX, orderTileY := u.currentOrder.targetActor.(*building).getUnitPlacementCoords()

	if orderTileX == utx && orderTileY == uty {
		u.currentOrder.code = ORDER_HARVEST
		u.currentAction.code = ACTION_ENTER_BUILDING
		u.currentAction.targetActor = u.currentOrder.targetActor
		return
	}

	b.SetActionForUnitForPathTo(u, orderTileX, orderTileY)
}

func (b *battlefield) SetActionForUnitForPathTo(u *unit, tx, ty int) {
	utx, uty := geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)

	path := b.findPathForUnitTo(u, tx, ty, true)
	vx, vy := path.GetNextStepVector()

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

	u.currentAction.code = ACTION_MOVE
	u.currentAction.targetTileX = utx + vx
	u.currentAction.targetTileY = uty + vy
}
