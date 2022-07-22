package main

// Action is what actor does this moment (move, build, etc) right away.
// DO NOT confuse with order. Order is an intent, whereas action is a low-level activity in progress.
type action struct {
	code                     int
	targetTileX, targetTileY int
	targetRotation           int
	// targetX, targetY         float64
}

const (
	ACTION_WAIT = iota
	ACTION_ROTATE
	ACTION_MOVE
)
