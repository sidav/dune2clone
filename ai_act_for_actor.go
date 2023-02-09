package main

func (ai *aiStruct) actForUnit(b *battlefield, u *unit) {
	if u.getStaticData().MaxCargoAmount > 0 {
		if ai.shouldUnitBeSentForRepairs(u) {
			ai.sendUnitForRepairs(u)
			return
		}
	} else {
		if ai.shouldUnitBeSentForRepairs(u) && rnd.OneChanceFrom(3) {
			ai.sendUnitForRepairs(u)
		}
	}
	if u.currentAction.code != ACTION_WAIT {
		return
	}
	if u.currentOrder.code != ORDER_NONE {
		return
	}
	u.currentOrder.resetOrder()
	if u.getStaticData().IsTransport {
		return // don't touch transports, they're automated!
	}
	if u.getStaticData().CanBeDeployed {
		ai.deployDeployableUnitSomewhere(b, u)
		return
	}
	if u.getStaticData().MaxCargoAmount > 0 {
		u.currentOrder.code = ORDER_HARVEST
		return
	}
	if len(u.turrets) > 0 {
		if !ai.isUnitInAnyTaskForce(u) {
			ai.assignUnitToTaskForce(u)
		}
		// other orders are in task-force based decisions
	}
}

func (ai *aiStruct) actForBuilding(b *battlefield, bld *building) {
	if bld.currentOrder.code == ORDER_WAIT_FOR_BUILDING_PLACEMENT {
		ai.placeBuilding(b, bld, bld.currentOrder.targetActor.(*building))
		return
	}
	if bld.currentAction.code != ACTION_WAIT {
		return
	}
	if bld.getHitpointsPercentage() < 50 || (!ai.isPoor() && bld.getHitpointsPercentage() < 100 && rnd.OneChanceFrom(10)) {
		bld.isRepairingSelf = true
	}
	if bld.getStaticData().Builds != nil && ai.checkIfShouldBuildNow() {
		bld.currentOrder.code = ORDER_BUILD
		bld.currentOrder.targetActorCode = int(ai.selectWhatToBuild(bld))
		ai.alreadyOrderedBuildThisTick = true
		return
	}

	if bld.getStaticData().Produces != nil && (!ai.areAllTaskForcesFull() || rnd.OneChanceFrom(20)) &&
		(!ai.isPoor() || rnd.OneChanceFrom(40)) {

		bld.currentOrder.code = ORDER_PRODUCE
		bld.currentOrder.targetActorCode = ai.selectWhatToProduce(bld)
		return
	}
}
