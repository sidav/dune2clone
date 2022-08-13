package main

type aiAnalytics struct {
	buildings  int
	builders   int
	eco        int
	production int
	defenses   int

	units int
}

func (aa *aiAnalytics) reset() {
	*aa = aiAnalytics{}
}

func (aa *aiAnalytics) increaseCountersForBuilding(bld *building) {
	aa.buildings++
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
	}
	if bld.currentAction.code == ACTION_BUILD {
		// count not yet built structures too
		if underConstruction, ok := bld.currentAction.targetActor.(*building); ok {
			aa.increaseCountersForBuilding(underConstruction)
		}
	}
}

func (ai *aiStruct) aiAnalyze(b *battlefield) {
	debugWritef("AI %s ANALYZE: It is tick %d\n", ai.name, b.currentTick)
	debugWritef("AI %s ANALYZE: I have %.f money\n", ai.name, ai.controlsFaction.getMoney())

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
			ai.current.units++
		}
	}
	debugWritef("AI: analyze shows that %+v\n", ai.current)
}
