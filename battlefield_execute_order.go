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
	orderTileX, orderTileY := u.currentOrder.targetTileX, u.currentOrder.targetTileY
	path := b.findPathForUnitTo(u, orderTileX, orderTileY)
	vx, vy := path.GetNextStepVector()

	if vx*vy != 0 {
		if b.currentTick % 2 == 0 {
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
