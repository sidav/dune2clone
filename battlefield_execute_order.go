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
	path := b.findPathForUnitTo(u, orderTileX, orderTileY, false)
	vx, vy := path.GetNextStepVector()

	if vx*vy != 0 {
		if b.currentTick%2 == 0 {
			vx = 0
		} else {
			vy = 0
		}
	}

	u.currentAction.code = ACTION_MOVE
	u.currentAction.targetTileX = utx + vx
	u.currentAction.targetTileY = uty + vy

	if utx == orderTileX && uty == orderTileY {
		u.currentOrder.code = ORDER_NONE
	}
}

func (b *battlefield) executeHarvestOrder(u *unit) {
	if u.currentCargoAmount >= u.getStaticData().maxCargoAmount {
		u.currentOrder.code = ORDER_RETURN_TO_REFINERY
	}

	utx, uty := geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
	orderTileX, orderTileY := u.currentOrder.targetTileX, u.currentOrder.targetTileY

	if b.tiles[orderTileX][orderTileY].resourcesAmount == 0 {
		// find resources coords
		rx, ry := b.getCoordsOfClosestEmptyTileWithResourcesTo(orderTileX, orderTileY)
		if rx < 0 || ry < 0 {
			return
		}
		u.currentOrder.targetTileX, u.currentOrder.targetTileY = rx, ry
		orderTileX, orderTileY = rx, ry
	}

	path := b.findPathForUnitTo(u, orderTileX, orderTileY, false)
	vx, vy := path.GetNextStepVector()

	if vx*vy != 0 {
		if b.currentTick%2 == 0 {
			vx = 0
		} else {
			vy = 0
		}
	}

	u.currentAction.code = ACTION_MOVE
	u.currentAction.targetTileX = utx + vx
	u.currentAction.targetTileY = uty + vy

	if utx == orderTileX && uty == orderTileY {
		u.currentAction.code = ACTION_HARVEST
	}
}

func (b *battlefield) executeReturnResourcesOrder(u *unit) {
	// x, y := u.centerX, u.centerY
	utx, uty := geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
	// for this, target tile is resource tile. Target refinery is in targetActor.
	if u.currentOrder.targetActor == nil {
		// find refinery TODO: find closest one
		for i := b.buildings.Front(); i != nil; i = i.Next() {
			if bld, ok := i.Value.(*building); ok {
				if bld.getStaticData().receivesResources && bld.getFaction() == u.getFaction() {
					u.currentOrder.targetActor = bld
					break
				}
			}
		}
	}
	if u.currentOrder.targetActor == nil {
		// nothing found, doing nothing
		return
	}
	orderTileX, orderTileY := u.currentOrder.targetActor.(*building).topLeftX, u.currentOrder.targetActor.(*building).topLeftY
	orderTileX += u.currentOrder.targetActor.(*building).getStaticData().unitPlacementX
	orderTileY += u.currentOrder.targetActor.(*building).getStaticData().unitPlacementY

	if orderTileX == utx && orderTileY == uty {
		u.currentOrder.code = ORDER_HARVEST
		u.currentAction.code = ACTION_ENTER_BUILDING
		u.currentAction.targetActor = u.currentOrder.targetActor
		return
	}

	path := b.findPathForUnitTo(u, orderTileX, orderTileY, true)
	vx, vy := path.GetNextStepVector()

	if vx*vy != 0 {
		if b.currentTick%2 == 0 {
			vx = 0
		} else {
			vy = 0
		}
	}

	u.currentAction.code = ACTION_MOVE
	u.currentAction.targetTileX = utx + vx
	u.currentAction.targetTileY = uty + vy
}
