package main

// Mission is a received order. Missions are assigned by player and are executed by assigning actions to unit.
type mission struct {
	code                     int
	targetTileX, targetTileY int
	// targetX, targetY         float64
}

const (
	MISSION_WAIT = iota
	MISSION_MOVE
)

