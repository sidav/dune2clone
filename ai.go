package main

type aiStruct struct {
	name            string
	controlsFaction *faction
	moneyPoorMax    float64
	moneyRichMin    float64
	current         aiAnalytics
	desired         aiAnalytics
	max             aiAnalytics

	alreadyOrderedBuildThisTick bool
}

func (ai *aiStruct) isPoor() bool {
	return ai.controlsFaction.getMoney() <= ai.moneyPoorMax
}

func (ai *aiStruct) isRich() bool {
	return ai.controlsFaction.getMoney() > ai.moneyRichMin
}

func (ai *aiStruct) aiControl(b *battlefield) {
	// debugWritef("AI ACT: It is tick %d\n", b.currentTick)
	ai.alreadyOrderedBuildThisTick = false
	for i := b.buildings.Front(); i != nil; i = i.Next() {
		if bld, ok := i.Value.(*building); ok {
			if bld.getFaction() == ai.controlsFaction {
				ai.actForBuilding(b, bld)
			}
		}
	}
	for i := b.units.Front(); i != nil; i = i.Next() {
		if unt, ok := i.Value.(*unit); ok {
			if unt.getFaction() == ai.controlsFaction {
				ai.actForUnit(b, unt)
			}
		}
	}
}

func (ai *aiStruct) actForUnit(b *battlefield, u *unit) {
	if u.currentAction.code != ACTION_WAIT || u.getStaticData().maxCargoAmount > 0 {
		return
	}
	if u.currentOrder.code != ORDER_NONE {
		return
	}
	u.currentOrder.resetOrder()
	if u.getStaticData().canBeDeployed {
		ai.deployDeployableUnitSomewhere(b, u)
		return
	}
	if len(u.turrets) > 0 {
		// temporary solution. TODO: task-force based decisions
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
		if selectedTarget != nil && b.canActorAttackActor(u, selectedTarget) {
			u.currentOrder.code = ORDER_ATTACK
			u.currentOrder.targetActor = selectedTarget
		}
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
	if bld.getStaticData().builds != nil && ai.current.nonDefenseBuildings < ai.max.nonDefenseBuildings && !ai.alreadyOrderedBuildThisTick {
		bld.currentOrder.code = ORDER_BUILD
		bld.currentOrder.targetActorCode = int(ai.selectWhatToBuild(bld))
		ai.alreadyOrderedBuildThisTick = true
		return
	}
	if bld.getStaticData().produces != nil && (!ai.isPoor() || rnd.OneChanceFrom(50)) {
		bld.currentOrder.code = ORDER_PRODUCE
		bld.currentOrder.targetActorCode = ai.selectWhatToProduce(bld)
		return
	}
}
