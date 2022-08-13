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

func (ai *aiStruct) aiAnalyze(b *battlefield) {
	debugWritef("AI %s ANALYZE: It is tick %d\n", ai.name, b.currentTick)
	debugWritef("AI %s ANALYZE: I have %.f money\n", ai.name, ai.controlsFaction.getMoney())

	ai.current.reset()

	for i := b.buildings.Front(); i != nil; i = i.Next() {
		if bld, ok := i.Value.(*building); ok {
			if bld.getFaction() == ai.controlsFaction {
				ai.current.buildings++
				if bld.getStaticData().receivesResources {
					ai.current.eco++
				}
				if bld.getStaticData().produces != nil {
					ai.current.production++
				}
				if bld.getStaticData().builds != nil {
					ai.current.builders++
				}
				if bld.turret != nil {
					ai.current.defenses++
				}
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
