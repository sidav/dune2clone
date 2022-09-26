package main

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

func (ai *aiStruct) assignUnitToTaskForce(u *unit) {
	selectedTf := ai.taskForces[0]
	for _, tf := range ai.taskForces {
		if tf.getSize() < selectedTf.getSize() && tf.getSize() < tf.desiredSize {
			selectedTf = tf
		}
	}
	selectedTf.addUnit(u)
}

func (ai *aiStruct) giveOrdersToAllTaskForces(b *battlefield) {
	for _, tf := range ai.taskForces {
		if tf.nextTickToGiveOrders > b.currentTick || tf.getSize() == 0 {
			continue
		}
		tf.cleanup() // TODO: maybe call it from another place?..
		switch tf.designation {
		case AITF_DESIGNATION_ATTACK:
			ai.giveOrderToAttackTaskForce(b, tf)
		case AITF_DESIGNATION_DEFEND:
			ai.giveOrderToDefendingTaskForce(b, tf)
		}
	}
}

func (ai *aiStruct) giveOrderToAttackTaskForce(b *battlefield, tf *aiTaskForce) {
	full := tf.getSize() >= tf.desiredSize
	if full || tf.noRetreat {
		if tf.target != nil {
			tf.noRetreat = true
			for _, u := range tf.units {
				u.currentOrder.code = ORDER_ATTACK
				u.currentOrder.targetActor = tf.target
				tf.nextTickToGiveOrders = b.currentTick+5*DESIRED_TPS
			}
		} else {
			tf.target = ai.findTargetForAttack(b)
		}
	}
	if tf.target == nil {
		ai.giveRoamOrderToTaskForce(b, tf)
	}
}

func (ai *aiStruct) giveOrderToDefendingTaskForce(b *battlefield, tf *aiTaskForce) {
	const basePatrolRadius = 20
	if tf.target != nil {
		for _, u := range tf.units {
			u.currentOrder.code = ORDER_ATTACK
			u.currentOrder.targetActor = tf.target
			tf.nextTickToGiveOrders = b.currentTick + 5*DESIRED_TPS
		}
	} else {
		tf.target = ai.findTargetNearBase(b, basePatrolRadius)
	}
	if tf.target == nil || !ai.isActorInRangeFromBase(tf.target, basePatrolRadius) {
		ai.giveRoamOrderToTaskForce(b, tf)
	}
}

func (ai *aiStruct) giveRoamOrderToTaskForce(b *battlefield, tf *aiTaskForce) {
	const radius = 20
	if ai.currBaseCenterX > 0 && ai.currBaseCenterY > 0 {
		for i := 0; i < 100; i++ {
			coordX := rnd.RandInRange(ai.currBaseCenterX-radius, ai.currBaseCenterX+radius)
			coordY := rnd.RandInRange(ai.currBaseCenterY-radius, ai.currBaseCenterY+radius)
			if b.isTileClearToBeMovedInto(coordX, coordY, nil) {
				for _, u := range tf.units {
					u.currentOrder.code = ORDER_MOVE
					u.currentOrder.setTargetTileCoords(coordX, coordY)
				}
				return
			}
		}
		tf.nextTickToGiveOrders = b.currentTick + DESIRED_TPS*10
	}
}

func (ai *aiStruct) findTargetNearBase(b *battlefield, radius int) actor {
	for i := b.units.Front(); i != nil; i = i.Next() {
		currActor := i.Value.(actor)
		if currActor.getFaction() == ai.controlsFaction {
			continue
		}
		if ai.isActorInRangeFromBase(currActor, radius) {
			return currActor
		}
	}
	return nil
}

func (ai *aiStruct) findTargetForAttack(b *battlefield) actor {
	var selectedTarget actor = nil
	attackBuildings := rnd.OneChanceFrom(3)
	if attackBuildings {
		for i := b.buildings.Front(); i != nil; i = i.Next() {
			currActor := i.Value.(actor)
			if currActor.getFaction() != ai.controlsFaction {
				if selectedTarget == nil || rnd.OneChanceFrom(4) {
					selectedTarget = currActor
				}
			}
		}
	} else {
		for i := b.units.Front(); i != nil; i = i.Next() {
			currActor := i.Value.(actor)
			if currActor.getFaction() != ai.controlsFaction {
				if selectedTarget == nil || rnd.OneChanceFrom(4) {
					selectedTarget = currActor
				}
			}
		}
	}
	return selectedTarget
}
