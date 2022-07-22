package main

// Orders are assigned by player and are executed by assigning actions to unit.
type order struct {
	code                     orderCode
	targetTileX, targetTileY int
	// targetX, targetY         float64
}

type orderCode uint8

const (
	ORDER_NONE orderCode = iota
	ORDER_MOVE
)

