package main

// Action is what actor does this moment (move, build, etc) right away.
// DO NOT confuse with order. Order is an intent, whereas action is a low-level activity in progress.
type action struct {
	code                     int
	targetTileX, targetTileY int
	targetRotation           int
	// targetX, targetY         float64
	targetActor actor

	currentFailuresCount uint8
	failedContinuously   bool
	interruptMovement    bool

	// construction-related
	moneySpentOnAction                    int
	completionAmount, maxCompletionAmount float64
	builtAs                               buildTypeCode
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
		return "Entering"
	case ACTION_DEPLOY:
		return "Deploying"

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

func (a *action) fail(resetIfFailedContinuously bool) {
	a.currentFailuresCount++
	if a.currentFailuresCount > 25 {
		a.failedContinuously = true
		if resetIfFailedContinuously {
			a.code = ACTION_WAIT
		}
	}
}

func (a *action) resetAction() {
	a.targetTileX = -1
	a.targetTileY = -1
	a.targetActor = nil
	a.maxCompletionAmount = 0
	a.code = ACTION_WAIT
	a.completionAmount = 0
	a.currentFailuresCount = 0
	a.failedContinuously = false
}

func (a *action) getCompletionPercent() int {
	if a.code == ACTION_BUILD {
		if b, ok := a.targetActor.(*building); ok {
			if b.getStaticData().buildType == BTYPE_PLACE_FIRST {
				return a.targetActor.getCurrentAction().getCompletionPercent()
				// int(100*a.targetActor.(*building).currentAction.completionAmount) / a.targetActor.(*building).getStaticData().maxHitpoints
			}
			return int(100*a.completionAmount) / (b.getStaticData().buildTime * (config.TargetTPS / config.Engine.BuildingsActionPeriod))
		}
		if b, ok := a.targetActor.(*unit); ok {
			return int(100*a.completionAmount) / (b.getStaticData().BuildTime * (config.TargetTPS / config.Engine.BuildingsActionPeriod))
		}
	}
	if a.maxCompletionAmount > 0 {
		return int((100 * a.completionAmount) / a.maxCompletionAmount)
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
	ACTION_BEING_BUILT
	// unit or tower cannon
	ACTION_ROTATE
	ACTION_HARVEST

	ACTION_ENTER_BUILDING
	ACTION_DEPLOY

	ACTION_AIR_APPROACH_LAND_TILE
	ACTION_AIR_APPROACH_ACTOR
	ACTION_AIR_PICK_UNIT_UP
	ACTION_AIR_DROP_UNIT
)
