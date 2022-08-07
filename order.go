package main

// Orders are assigned by player and are executed by assigning actions to unit.
type order struct {
	code                     orderCode
	targetTileX, targetTileY int
	targetActor              actor
	// targetX, targetY         float64
}

type orderCode uint8

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
