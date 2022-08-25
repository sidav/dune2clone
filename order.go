package main

// Orders are assigned by player and are executed by assigning actions to unit.
type order struct {
	code                       orderCode
	targetTileX, targetTileY   int
	targetTile2X, targetTile2Y int // when only single target tile coords are not enough
	targetActor                actor
	targetActorCode            int

	dispatchCalled bool
	// targetX, targetY         float64
}

type orderCode uint8

func (o *order) resetOrder() {
	o.code = ORDER_NONE
	o.targetTileX = -1
	o.targetTileY = -1
	o.targetTile2X = -1
	o.targetTile2Y = -1
	o.targetActor = nil
	o.dispatchCalled = false
}

func (o *order) setTargetTileCoords(x, y int) {
	o.targetTileX, o.targetTileY = x, y
}

func (o *order) getTextDescription() string {
	switch o.code {
	case ORDER_NONE:
		return "Standing by"
	case ORDER_MOVE:
		return "Moving"
	case ORDER_HARVEST:
		return "Harvesting"
	case ORDER_RETURN_TO_REFINERY:
		return "Returning with cargo"
	case ORDER_BUILD:
		return "Building"
	case ORDER_PRODUCE:
		return "Training"
	case ORDER_WAIT_FOR_BUILDING_PLACEMENT:
		return "Waiting for placement"
	case ORDER_MOVE_TO_REPAIR:
		return "Moving for repairs"
	case ORDER_CARRY_UNIT_TO_TARGET_COORDS:
		return "Transporting"
	case ORDER_CANCEL_BUILD:
		return ""
	}
	panic("No order description!")
}

const (
	ORDER_NONE orderCode = iota
	ORDER_MOVE
	ORDER_HARVEST
	ORDER_RETURN_TO_REFINERY
	ORDER_MOVE_TO_REPAIR
	ORDER_BUILD
	ORDER_PRODUCE
	ORDER_WAIT_FOR_BUILDING_PLACEMENT
	ORDER_CANCEL_BUILD
	// aircraft orders
	ORDER_CARRY_UNIT_TO_TARGET_COORDS
)
