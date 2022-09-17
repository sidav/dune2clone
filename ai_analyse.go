package main

type aiAnalytics struct {
	// buildings
	nonDefenseBuildings int
	defenses            int

	builders   int
	eco        int
	production int

	// units
	combatUnits    int
	nonCombatUnits int

	harvesters int
	transports int
}

func (aa *aiAnalytics) reset() {
	*aa = aiAnalytics{}
}

func (aa *aiAnalytics) increaseCountersForBuilding(bld *building) {
	if bld.getStaticData().receivesResources {
		aa.eco++
	}
	if bld.getStaticData().produces != nil {
		aa.production++
	}
	if bld.getStaticData().builds != nil {
		aa.builders++
	}
	if bld.turret != nil {
		aa.defenses++
	} else {
		aa.nonDefenseBuildings++
	}
	if bld.currentAction.code == ACTION_BUILD {
		// count not yet built structures too
		if underConstruction, ok := bld.currentAction.targetActor.(*building); ok {
			aa.increaseCountersForBuilding(underConstruction)
		}
	}
}

func (aa *aiAnalytics) increaseCountersForUnit(ai *aiStruct, u *unit) {
	switch ai.deduceUnitFunction(u.code) {
	case "harvester":
		ai.current.harvesters++
		ai.current.nonCombatUnits++
	case "transport":
		ai.current.transports++
		ai.current.nonCombatUnits++
	case "combat":
		ai.current.combatUnits++
	}
}

func (ai *aiStruct) aiAnalyze(b *battlefield) {
	// debugWritef("AI %s ANALYZE: It is tick %d\n", ai.name, b.currentTick)
	// debugWritef("AI %s ANALYZE: I have %.f money\n", ai.name, ai.controlsFaction.getMoney())

	ai.current.reset()

	for i := b.buildings.Front(); i != nil; i = i.Next() {
		if bld, ok := i.Value.(*building); ok {
			if bld.getFaction() == ai.controlsFaction {
				ai.current.increaseCountersForBuilding(bld)
			}
		}
	}

	for i := b.units.Front(); i != nil; i = i.Next() {
		// debugWritef("req: %d,%d; act: %f, %f -> %d, %d \n", x, y, b.units[i].centerX, b.units[i].centerY, tx, ty)
		if i.Value.(*unit).getFaction() == ai.controlsFaction {
			ai.current.increaseCountersForUnit(ai, i.Value.(*unit))
		}
	}
	// debugWritef("AI: analyze shows that %+v\n", ai.current)
}
