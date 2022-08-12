package main

import "dune2clone/geometry"

type aiStruct struct {
	controlsFaction *faction
	current         aiAnalytics
	desired         aiAnalytics
	max             aiAnalytics
}

func (ai *aiStruct) aiControl(b *battlefield) {
	debugWritef("AI ACT: It is tick %d\n", b.currentTick)
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
	u.currentOrder.code = ORDER_MOVE
	for !b.isTileClearToBeMovedInto(u.currentOrder.targetTileX, u.currentOrder.targetTileY, u) {
		u.currentOrder.targetTileX = rnd.Rand(len(b.tiles))
		u.currentOrder.targetTileY = rnd.Rand(len(b.tiles[0]))
	}
}

func (ai *aiStruct) actForBuilding(b *battlefield, bld *building) {
	if bld.currentOrder.code == ORDER_WAIT_FOR_BUILDING_PLACEMENT {
		ai.placeBuilding(b, bld, bld.currentAction.targetActor.(*building))
	}
	if bld.currentAction.code != ACTION_WAIT {
		return
	}
	if bld.getStaticData().builds != nil && ai.current.buildings < ai.max.buildings {
		bld.currentOrder.code = ORDER_BUILD
		bld.currentOrder.targetActorCode = ai.selectWhatToBuild(bld)
	}
	if bld.getStaticData().produces != nil && ai.current.units < ai.max.units {
		code := bld.getStaticData().produces[rnd.Rand(len(bld.getStaticData().produces))]
		bld.currentOrder.code = ORDER_PRODUCE
		bld.currentOrder.targetActorCode = code
	}
}

func (ai *aiStruct) selectWhatToBuild(builder *building) int {
	availableCodes := builder.getStaticData().builds
	// first of all, AI needs eco
	if ai.current.eco < 1 || ai.current.eco < ai.desired.eco && rnd.OneChanceFrom(25) {
		for _, code := range availableCodes {
			if sTableBuildings[code].receivesResources {
				return code
			}
		}
	}
	if ai.current.production < 1 || ai.current.production < ai.desired.production && rnd.OneChanceFrom(25) {
		for _, code := range availableCodes {
			if sTableBuildings[code].produces != nil {
				return code
			}
		}
	}
	if ai.current.builders < ai.desired.builders && rnd.OneChanceFrom(25) {
		for _, code := range availableCodes {
			if sTableBuildings[code].builds != nil {
				return code
			}
		}
	}
	return availableCodes[rnd.Rand(len(availableCodes))]
}

func (ai *aiStruct) placeBuilding(b *battlefield, builder, whatIsBuilt *building) {
	startX, startY := geometry.TrueCoordsToTileCoords(builder.getPhysicalCenterCoords())
	sx, sy := geometry.SpiralSearchForConditionFrom(
		func(x, y int) bool {
			return b.canBuildingBePlacedAt(whatIsBuilt, x, y, 1,false)
		},
		startX, startY, 16, 0)
	if sx != -1 && sy != -1 {
		builder.currentOrder.targetTileX = sx
		builder.currentOrder.targetTileY = sy
	}
}
