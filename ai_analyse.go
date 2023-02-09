package main

import "dune2clone/geometry"

type aiAnalytics struct {
	// buildings
	nonDefenseBuildings int
	defenses            int

	builders     int
	eco          int
	production   int
	repairDepots int

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
	if bld.getStaticData().ReceivesResources {
		aa.eco++
	}
	if bld.getStaticData().Produces != nil {
		aa.production++
	}
	if bld.getStaticData().Builds != nil {
		aa.builders++
	}
	if bld.getStaticData().RepairsUnits {
		aa.repairDepots++
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
	currTrueBaseCenterX, currTrueBaseCenterY := 0.0, 0.0
	totalBlds := 0
	for i := b.buildings.Front(); i != nil; i = i.Next() {
		if bld, ok := i.Value.(*building); ok {
			if bld.getFaction() == ai.controlsFaction {
				ai.current.increaseCountersForBuilding(bld)
				// calculating current base center
				bcx, bcy := bld.getPhysicalCenterCoords()
				currTrueBaseCenterX += bcx
				currTrueBaseCenterY += bcy
				totalBlds++
			}
		}
	}
	if totalBlds > 0 {
		ai.currBaseCenterX, ai.currBaseCenterY =
			geometry.TrueCoordsToTileCoords(currTrueBaseCenterX/float64(totalBlds), currTrueBaseCenterY/float64(totalBlds))
	}

	for i := b.units.Front(); i != nil; i = i.Next() {
		// debugWritef("req: %d,%d; act: %f, %f -> %d, %d \n", x, y, b.units[i].centerX, b.units[i].centerY, tx, ty)
		if i.Value.(*unit).getFaction() == ai.controlsFaction {
			ai.current.increaseCountersForUnit(ai, i.Value.(*unit))
		}
	}
	// ai.debugWritef("analyze shows that %+v\n", ai)
}
