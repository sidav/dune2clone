package main

import "dune2clone/geometry"

func (ai *aiStruct) selectWhatToBuild(builder *building) int {
	availableCodes := builder.getStaticData().builds
	// first of all, AI needs eco
	if ai.current.eco < 1 || ai.current.eco < ai.desired.eco && rnd.OneChanceFrom(5) {
		for _, code := range availableCodes {
			if sTableBuildings[code].receivesResources {
				return code
			}
		}
	}
	if ai.current.production < 1 || ai.current.production < ai.desired.production && rnd.OneChanceFrom(4) {
		for _, code := range availableCodes {
			if sTableBuildings[code].produces != nil {
				return code
			}
		}
	}
	if ai.controlsFaction.maxResources - ai.controlsFaction.resources < 500 && rnd.OneChanceFrom(3) {
		for _, code := range availableCodes {
			if sTableBuildings[code].storageAmount > 0 && !sTableBuildings[code].receivesResources {
				return code
			}
		}
	}
	if ai.current.builders < ai.desired.builders && rnd.OneChanceFrom(10) {
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
