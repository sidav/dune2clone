package main

import "dune2clone/geometry"

type aiController struct {
	controlsFaction *faction
}

func (ai *aiController) aiControl(b *battlefield) {
	debugWritef("AI: It is tick %d\n", b.currentTick)
	debugWritef("AI: I have %.f money\n", ai.controlsFaction.money)
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

func (ai *aiController) actForUnit(b *battlefield, u *unit) {
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

func (ai *aiController) actForBuilding(b *battlefield, bld *building) {
	if bld.currentAction.code == ACTION_BUILD && bld.currentAction.getCompletionPercent() >= 100 {
		ai.placeBuilding(b, bld, bld.currentAction.targetActor.(*building))
	}
	if bld.currentAction.code != ACTION_WAIT {
		return
	}
	if bld.getStaticData().builds != nil {
		code := bld.getStaticData().builds[rnd.Rand(len(bld.getStaticData().builds))]
		bld.currentAction.code = ACTION_BUILD
		bld.currentAction.targetActor = createBuilding(code, 0, 0, bld.faction)
	}
	if bld.getStaticData().produces != nil {
		code := bld.getStaticData().produces[rnd.Rand(len(bld.getStaticData().produces))]
		bld.currentAction.code = ACTION_BUILD
		bld.currentAction.targetActor = createUnit(code, 0, 0, bld.faction)
	}
}

func (ai *aiController) placeBuilding(b *battlefield, builder, whatIsBuilt *building) {
	startX, startY := geometry.TrueCoordsToTileCoords(builder.getPhysicalCenterCoords())
	sx, sy := geometry.SpiralSearchForConditionFrom(
		func(x, y int) bool {
			return b.canBuildingBePlacedAt(whatIsBuilt, x, y, false)
		},
		startX, startY, 5, 0)
	if sx != -1 && sy != -1 {
		whatIsBuilt.topLeftX, whatIsBuilt.topLeftY = sx, sy
		b.addActor(whatIsBuilt)
		builder.currentAction.reset()
	}
}
