package main

// Orders are assigned by player and are executed by assigning actions to unit.
type order struct {
	code                     int
	targetTileX, targetTileY int
	// targetX, targetY         float64
}

const (
	ORDER_NONE = iota
	ORDER_MOVE
)

