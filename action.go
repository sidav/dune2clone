package main

// Action is what actor does this moment (move, build, etc) right away.
// DO NOT confuse with order. Order is an intent, whereas action is a low-level activity in progress.
type action struct {
	code                     int
	targetTileX, targetTileY int
	targetRotation           int
	// targetX, targetY         float64
	targetActor actor

	// construction-related
	moneySpentOnAction int
	completionAmount   float64
}

func (a *action) getTextDescription() string {
	switch a.code {
	case ACTION_WAIT:
		return "Awaiting orders"
	case ACTION_MOVE:
		return "Moving to position"
	case ACTION_ROTATE:
		return "Rotating"
	case ACTION_BUILD:
		return "Constructing"
	case ACTION_HARVEST:
		return "Harvesting"
	case ACTION_ENTER_BUILDING:
		return "Harvesting"

	case ACTION_AIR_APPROACH_LAND_TILE:
		return "Slowly flying to"
	case ACTION_AIR_PICK_UNIT_UP:
		return "Picking unit up"
	case ACTION_AIR_DROP_UNIT:
		return "Dropping"
	}
	return "NO DESC"
	panic("No action description!")
}

func (a *action) setTargetTileCoords(x, y int) {
	a.targetTileX, a.targetTileY = x, y
}

func (a *action) reset() {
	a.targetTileX = -1
	a.targetTileY = -1
	a.targetActor = nil
	a.code = ACTION_WAIT
	a.completionAmount = 0
}

func (a *action) getCompletionPercent() int {
	if a.code == ACTION_BUILD {
		if b, ok := a.targetActor.(*building); ok {
			return int(100*a.completionAmount) / (b.getStaticData().buildTime * (DESIRED_FPS / BUILDINGS_ACTIONS_TICK_EACH))
		}
		if b, ok := a.targetActor.(*unit); ok {
			return int(100*a.completionAmount) / (b.getStaticData().buildTime * (DESIRED_FPS / BUILDINGS_ACTIONS_TICK_EACH))
		}
	}
	return -1
	// panic("Something is wrong in %")
}

const (
	ACTION_WAIT = iota
	// unit-only:
	ACTION_MOVE
	// production:
	ACTION_BUILD
	// unit or tower cannon
	ACTION_ROTATE
	ACTION_HARVEST

	ACTION_ENTER_BUILDING

	ACTION_AIR_APPROACH_LAND_TILE
	ACTION_AIR_PICK_UNIT_UP
	ACTION_AIR_DROP_UNIT
)
