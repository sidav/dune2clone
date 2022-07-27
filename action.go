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
	completionAmount   int
}

func (a *action) reset() {
	a.targetActor = nil
	a.code = ACTION_WAIT
	a.completionAmount = 0
}

func (a *action) getCompletionPercent() int {
	if a.code == ACTION_BUILD {
		if b, ok := a.targetActor.(*building); ok {
			return 100*a.completionAmount/(b.getStaticData().buildTime * (DESIRED_FPS/BUILDINGS_ACTIONS_TICK_EACH))
		} else {
		}
	}
	panic("Something is wrong in %")
}

const (
	ACTION_WAIT = iota
	// unit-only:
	ACTION_MOVE
	// production:
	ACTION_BUILD
	// unit or tower cannon
	ACTION_ROTATE
)
