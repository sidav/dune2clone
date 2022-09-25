package main

import "dune2clone/geometry"

func (ai *aiStruct) deployDeployableUnitSomewhere(b *battlefield, u *unit) {
	bld := createBuilding(u.getStaticData().deploysInto, 0, 0, u.faction)
	tx, ty := u.getTileCoords()
	if !b.canUnitBeDeployedAt(u, tx, ty) {
		depX, depY := -1, -1
		if ai.current.builders > 0 {
			depX, depY = geometry.SpiralSearchForFarthestConditionFrom(
				func(x, y int) bool {
					return b.canBuildingBePlacedAt(bld, x, y, 0, true) && rnd.OneChanceFrom(32)
				},
				tx, ty, 16, rnd.Rand(4),
			)
		} else {
			depX, depY = geometry.SpiralSearchForClosestConditionFrom(
				func(x, y int) bool {
					return b.canBuildingBePlacedAt(bld, x, y, 0, true) && rnd.OneChanceFrom(32)
				},
				tx, ty, 16, rnd.Rand(4),
			)
		}
		u.currentOrder.code = ORDER_MOVE
		u.currentOrder.setTargetTileCoords(depX, depY)
	} else {
		u.currentOrder.code = ORDER_DEPLOY
	}
}

func (ai *aiStruct) sendUnitForRepairs(u *unit) {
	u.currentOrder.resetOrder()
	u.currentOrder.code = ORDER_MOVE_TO_REPAIR
}

func (ai *aiStruct) shouldUnitBeSentForRepairs(u *unit) bool {
	return ai.current.repairDepots > 0 && u.getHitpointsPercentage() < 33
}
