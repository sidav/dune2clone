package main

type aiStruct struct {
	name            string
	controlsFaction *faction
	current         aiAnalytics
	desired         aiAnalytics
	max             aiAnalytics

	alreadyOrderedBuildThisTick bool
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
	if len(u.turrets) > 0 {
		for i := b.units.Front(); i != nil; i = i.Next() {
			var selectedTarget actor = nil
			if i.Value.(actor).getFaction() != ai.controlsFaction {
				if selectedTarget == nil || rnd.OneChanceFrom(4) {
					u.currentOrder.code = ORDER_ATTACK
					u.currentOrder.targetActor = i.Value.(actor)
				}
			}
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
	if bld.getStaticData().builds != nil && ai.current.buildings < ai.max.buildings && !ai.alreadyOrderedBuildThisTick {
		bld.currentOrder.code = ORDER_BUILD
		bld.currentOrder.targetActorCode = ai.selectWhatToBuild(bld)
		ai.alreadyOrderedBuildThisTick = true
		return
	}
	if bld.getStaticData().produces != nil && ai.current.units < ai.max.units {
		code := bld.getStaticData().produces[rnd.Rand(len(bld.getStaticData().produces))]
		bld.currentOrder.code = ORDER_PRODUCE
		bld.currentOrder.targetActorCode = code
	}
}
