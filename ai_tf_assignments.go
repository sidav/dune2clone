package main

func (ai *aiStruct) cleanupAiTaskForces() {
	for _, tf := range ai.taskForces {
		tf.cleanDead()
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

func (ai *aiStruct) giveOrdersToTaskForces(b *battlefield) {
	for _, tf := range ai.taskForces {
		switch tf.designation {
		case AITF_DESIGNATION_ATTACK:
			ai.giveOrderToAttackTaskForce(tf)
		case AITF_DESIGNATION_DEFEND:
			ai.giveOrderToDefendingTaskForce(tf)
		}
	}
}

func (ai *aiStruct) giveOrderToAttackTaskForce(tf *aiTaskForce) {

}

func (ai *aiStruct) giveOrderToDefendingTaskForce(tf *aiTaskForce) {

}
