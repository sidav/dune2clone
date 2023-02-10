package main

import "dune2clone/geometry"

func (ai *aiStruct) cleanupAiTaskForces() {
	for _, tf := range ai.taskForces {
		tf.cleanup()
	}
}

func (ai *aiStruct) isUnitInAnyTaskForce(u *unit) bool {
	for _, tf := range ai.taskForces {
		if tf.doesHaveUnit(u) {
			return true
		}
	}
	return false
}

func (ai *aiStruct) areAllTaskForcesFull() bool {
	for _, tf := range ai.taskForces {
		if !tf.isFull() {
			return false
		}
	}
	return true
}

func (ai *aiStruct) assignUnitToTaskForce(u *unit) {
	const minSpeedForBeingRecon = 0.14
	isUnitSuitableForRecon := u.getMovementSpeed() >= minSpeedForBeingRecon
	selectedTf := ai.taskForces[0]
	for _, tf := range ai.taskForces {
		if tf.getFullnessPercent() < selectedTf.getFullnessPercent() && !tf.isFull() ||
			// don't assign slow units to recon
			(selectedTf.mission == AITF_MISSION_RECON && !isUnitSuitableForRecon) {

			selectedTf = tf
		}
	}
	ai.debugWritef("adding %s to task force %d/%d\n", u.getName(), selectedTf.getSize(), selectedTf.desiredSize)
	selectedTf.addUnit(u)
}

func (ai *aiStruct) giveOrdersToAllTaskForces(b *battlefield) {
	for _, tf := range ai.taskForces {
		if tf.nextTickToGiveOrders > b.currentTick || tf.getSize() == 0 {
			continue
		}
		tf.cleanup() // TODO: maybe call it from another place?..
		switch tf.mission {
		case AITF_MISSION_ATTACK:
			ai.giveOrderToAttackTaskForce(b, tf)
		case AITF_MISSION_DEFEND:
			ai.giveOrderToDefendingTaskForce(b, tf)
		case AITF_MISSION_RECON:
			ai.giveReconOrderToTaskForce(b, tf)
		default:
			panic("No such task force mission!")
		}
	}
}

func (ai *aiStruct) giveOrderToAttackTaskForce(b *battlefield, tf *aiTaskForce) {
	if tf.shouldBeRetreated() {
		ai.debugWritef("Attack TF: size %d, retreat!\n", tf.getSize())
		tf.target = nil
	}
	if tf.isFull() || tf.target != nil {
		if tf.target != nil {
			ai.debugWritef("Attack TF: attack %s!\n", tf.target.getName())
			for _, u := range tf.units {
				if rnd.OneChanceFrom(4) {
					u.currentOrder.code = ORDER_ATTACK
				} else {
					u.currentOrder.code = ORDER_ATTACK_MOVE
					u.currentOrder.setTargetTileCoords(geometry.TrueCoordsToTileCoords(tf.target.getPhysicalCenterCoords()))
				}
				u.currentOrder.targetActor = tf.target
				tf.nextTickToGiveOrders = b.currentTick + 15*config.Engine.TicksPerNominalSecond
			}
		} else {
			tf.target = ai.findVisibleTargetForAttack(b)
		}
	}
	if tf.target == nil {
		if tf.isFull() {
			ai.debugWritef("Attack TF: battle recon!\n")
			ai.giveReconOrderToTaskForce(b, tf)
		} else {
			ai.giveRoamNearBaseOrderToTaskForce(b, tf)
		}
	}
}

func (ai *aiStruct) giveOrderToDefendingTaskForce(b *battlefield, tf *aiTaskForce) {
	const basePatrolRadius = 20
	if tf.target != nil {
		for _, u := range tf.units {
			u.currentOrder.code = ORDER_ATTACK_MOVE
			u.currentOrder.setTargetTileCoords(geometry.TrueCoordsToTileCoords(tf.target.getPhysicalCenterCoords()))
			u.currentOrder.targetActor = tf.target
			tf.nextTickToGiveOrders = b.currentTick + 10*config.Engine.TicksPerNominalSecond
		}
	} else {
		tf.target = ai.findVisibleTargetNearBase(b, basePatrolRadius)
	}
	if tf.target == nil || !ai.isActorInRangeFromBase(tf.target, basePatrolRadius) {
		ai.giveRoamNearBaseOrderToTaskForce(b, tf)
	}
}

func (ai *aiStruct) giveRoamNearBaseOrderToTaskForce(b *battlefield, tf *aiTaskForce) {
	const radius = 20
	if ai.currBaseCenterX > 0 && ai.currBaseCenterY > 0 {
		for i := 0; i < 100; i++ {
			coordX := rnd.RandInRange(ai.currBaseCenterX-radius, ai.currBaseCenterX+radius)
			coordY := rnd.RandInRange(ai.currBaseCenterY-radius, ai.currBaseCenterY+radius)
			if b.isTileClearToBeMovedInto(coordX, coordY, nil) {
				for _, u := range tf.units {
					u.currentOrder.resetOrder()
					u.currentOrder.code = ORDER_ATTACK_MOVE
					u.currentOrder.setTargetTileCoords(coordX, coordY)
				}
				return
			}
		}
		tf.nextTickToGiveOrders = b.currentTick + config.Engine.TicksPerNominalSecond*10
	}
}

func (ai *aiStruct) giveReconOrderToTaskForce(b *battlefield, tf *aiTaskForce) {
	minReconRange := 20
	maxReconRange, _ := b.getSize()
	coordX, coordY := -1, -1
	if rnd.OneChanceFrom(4) {
		coordX, coordY = geometry.SpiralSearchForFarthestConditionFrom(
			func(x, y int) bool {
				return b.areTileCoordsValid(x, y) && !ai.controlsFaction.seesTileAtCoords(x, y)
			},
			ai.currBaseCenterX, ai.currBaseCenterY, rnd.RandInRange(minReconRange, maxReconRange), rnd.Rand(4),
		)
	} else {
		// search close to base
		coordX, coordY = geometry.SpiralSearchForFarthestConditionFrom(
			func(x, y int) bool {
				return b.areTileCoordsValid(x, y) && !ai.controlsFaction.hasTileAtCoordsExplored(x, y)
			},
			ai.currBaseCenterX, ai.currBaseCenterY, rnd.RandInRange(minReconRange, maxReconRange), rnd.Rand(4))
		// if nowhere to search, search farther, but select closest unexplored tile
		//if coordX == -1 && coordY == -1 {
		//	coordX, coordY = geometry.SpiralSearchForClosestConditionFrom(
		//		func(x, y int) bool {
		//			return b.areTileCoordsValid(x, y) && !ai.controlsFaction.hasTileAtCoordsExplored(x, y)
		//		},
		//		ai.currBaseCenterX, ai.currBaseCenterY, rnd.RandInRange(reconRange/3, reconRange), rnd.Rand(4))
		//}
	}
	if coordX != -1 && coordY != -1 {
		for _, u := range tf.units {
			u.currentOrder.resetOrder()
			u.currentOrder.code = ORDER_MOVE
			u.currentOrder.setTargetTileCoords(coordX, coordY)
		}
	}
	tf.nextTickToGiveOrders = b.currentTick + 10*config.Engine.TicksPerNominalSecond
}

func (ai *aiStruct) findVisibleTargetNearBase(b *battlefield, radius int) actor {
	for i := b.units.Front(); i != nil; i = i.Next() {
		currActor := i.Value.(actor)
		if currActor.getFaction() == ai.controlsFaction {
			continue
		}
		if ai.isActorInRangeFromBase(currActor, radius) && b.canFactionSeeActor(ai.controlsFaction, currActor) {
			return currActor
		}
	}
	return nil
}

func (ai *aiStruct) findVisibleTargetForAttack(b *battlefield) actor {
	var selectedTarget actor = nil
	attackBuildings := rnd.OneChanceFrom(2)
	if attackBuildings {
		for i := b.buildings.Front(); i != nil; i = i.Next() {
			currActor := i.Value.(actor)
			if currActor.getFaction() != ai.controlsFaction && b.hasFactionExploredBuilding(ai.controlsFaction, currActor.(*building)) {
				if selectedTarget == nil || rnd.OneChanceFrom(4) {
					selectedTarget = currActor
				}
			}
		}
	} else {
		for i := b.units.Front(); i != nil; i = i.Next() {
			currActor := i.Value.(actor)
			if currActor.getFaction() != ai.controlsFaction && b.canFactionSeeActor(ai.controlsFaction, currActor) {
				if selectedTarget == nil || rnd.OneChanceFrom(4) {
					selectedTarget = currActor
				}
			}
		}
	}
	return selectedTarget
}
