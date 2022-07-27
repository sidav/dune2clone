package main

import (
	"math"
)

func (b *battlefield) executeActionForActor(a actor) {
	switch a.getCurrentAction().code {
	case ACTION_MOVE:
		if u, ok := a.(*unit); ok {
			b.executeMoveActionForUnit(u)
		} else {
			panic("Is not unit!")
		}
	case ACTION_BUILD:
		b.executeBuildActionForActor(a)
	}
}

func (b *battlefield) executeMoveActionForUnit(u *unit) {
	x, y := u.centerX, u.centerY
	tx, ty := tileCoordsToPhysicalCoords(u.currentAction.targetTileX, u.currentAction.targetTileY)

	if tx != x {
		if !isVectorDegreeEqualTo(tx-x, 0, u.chassisDegree) {
			u.rotateChassisTowardsVector(tx-x, 0)
			return
		}
		if math.Abs(tx-x) < u.getStaticData().movementSpeed {
			u.centerX = tx // source of movement lag :(
		} else {
			u.centerX += u.getStaticData().movementSpeed * (tx - x) / math.Abs(tx-x)
		}
	} else if ty != y {
		if !isVectorDegreeEqualTo(0, ty-y, u.chassisDegree) {
			u.rotateChassisTowardsVector(0, ty-y)
			return
		}
		if math.Abs(ty-y) < u.getStaticData().movementSpeed {
			u.centerY = ty // source of movement lag :(
		} else {
			u.centerY += u.getStaticData().movementSpeed * (ty - y) / math.Abs(ty-y)
		}
	}
	if areFloatsAlmostEqual(x, tx) && areFloatsAlmostEqual(y, ty) {
		u.centerX = tx
		u.centerY = ty
		u.currentAction.reset()
		// debugWritef("Tick %d: action finished\n", b.currentTick)
	}
}

func (b *battlefield) executeBuildActionForActor(a actor) {
	act := a.getCurrentAction()
	moneySpent := 0.0
	// calculate spending
	if bld, ok := act.targetActor.(*building); ok {
		moneySpent = float64(bld.getStaticData().cost) /
			float64(bld.getStaticData().buildTime*(DESIRED_FPS/BUILDINGS_ACTIONS_TICK_EACH))
	}
	if unt, ok := act.targetActor.(*unit); ok {
		moneySpent = float64(unt.getStaticData().cost) /
			float64(unt.getStaticData().buildTime*(DESIRED_FPS/BUILDINGS_ACTIONS_TICK_EACH))
	}
	// spend money
	if act.getCompletionPercent() < 100 && a.getFaction().money > moneySpent {
		a.getFaction().money -= moneySpent
		act.completionAmount++
	}
	// if it was a unit, place it right away
	if unt, ok := act.targetActor.(*unit); ok && act.getCompletionPercent() >= 100 {
		if bld, ok := a.(*building); ok {
			for x := bld.topLeftX - 1; x <= bld.topLeftX+bld.getStaticData().w; x++ {
				// for y := bld.topLeftY-1; y <= bld.topLeftY+bld.getStaticData().h; y++ {
				y := bld.topLeftY+bld.getStaticData().h
				if b.costMapForMovement(x, y) != -1 {
					unt.centerX, unt.centerY = tileCoordsToPhysicalCoords(x, y)
					// debugWritef("+%v", unt)
					b.addActor(unt)
					bld.currentAction.reset()
					return
				}
				// }
			}
		} else {
			panic("wat")
		}
	}
}
