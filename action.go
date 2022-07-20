package main

// Action is what actor does this moment (move, build, etc) right away.
// DO NOT confuse with missions (received orders). Mission is an intent, whereas action is an activity in progress.
type action struct {
	code                     int
	targetTileX, targetTileY int
	// targetX, targetY         float64
}

const (
	ACTION_WAIT = iota
	ACTION_MOVE
)
