package main

// Orders are assigned by player and are executed by assigning actions to unit.
type order struct {
	code                       orderCode
	targetTileX, targetTileY   int
	targetTile2X, targetTile2Y int // when only single target tile coords are not enough
	targetActor                actor
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
	}
	panic("No action description!")
}

const (
	ORDER_NONE orderCode = iota
	ORDER_MOVE
	ORDER_HARVEST
	ORDER_RETURN_TO_REFINERY
)
