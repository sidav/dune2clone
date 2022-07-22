package main

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
	}
}

func (b *battlefield) executeMoveOrder(u *unit) {
	// x, y := u.centerX, u.centerY
	utx, uty := trueCoordsToTileCoords(u.centerX, u.centerY)
	// todo: pathfinding
	otx, oty := u.currentOrder.targetTileX, u.currentOrder.targetTileY
	u.currentAction.code = ACTION_MOVE
	u.currentAction.targetTileX = utx
	u.currentAction.targetTileY = uty

	if abs(otx-utx) > abs(oty-uty) {
		u.currentAction.targetTileX = utx + sign(otx-utx)
	} else if oty != uty {
		u.currentAction.targetTileY = uty + sign(oty-uty)
	}
	// debugWritef("Order: %+v, action: %+v\n", u.currentOrder, u.currentAction)
	if utx == otx && uty == oty {
		u.currentOrder.code = ORDER_NONE
	}
}
