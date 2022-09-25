package main

import (
	"dune2clone/geometry"
)

type building struct {
	currentAction          action
	currentOrder           order
	currentHitpoints       int
	topLeftX, topLeftY     int // tile coords
	code                   buildingCode
	faction                *faction
	isSelected             bool
	turret                 *turret
	unitPlacedInside       *unit
	rallyTileX, rallytileY int
	isRepairingSelf        bool
}

func (b *building) isAlive() bool {
	return b.currentHitpoints > 0
}

func (b *building) getHitpoints() int {
	return b.currentHitpoints
}

func (b *building) getMaxHitpoints() int {
	return b.getStaticData().maxHitpoints
}

func (b *building) getHitpointsPercentage() int {
	return getPercentInt(b.currentHitpoints, b.getStaticData().maxHitpoints)
}

func createBuilding(code buildingCode, topLeftX, topLeftY int, fact *faction) *building {
	var turr *turret
	if sTableBuildings[code].turretCode != TRT_NONE {
		turr = &turret{code: sTableBuildings[code].turretCode, rotationDegree: 270}
	}
	return &building{
		code:             code,
		currentHitpoints: sTableBuildings[code].maxHitpoints,
		topLeftX:         topLeftX,
		topLeftY:         topLeftY,
		faction:          fact,
		turret:           turr,
		rallyTileX:       -1,
		rallytileY:       -1,
	}
}

func (b *building) markSelected(s bool) {
	b.isSelected = s
}

func (b *building) isUnderConstruction() bool {
	return b.currentAction.code == ACTION_BEING_BUILT
}

func (b *building) getVisionRange() int {
	if b.faction.lastAvailableEnergy < 0 || b.isUnderConstruction() {
		return 1
	}
	return 4
}

func (b *building) getDimensionsForConstructon() (int, int, int, int) {
	h := b.getStaticData().h
	// prevent closing bottom side for producing buildings
	if b.getStaticData().needsEmptyRowBelowWhenConstructing {
		h++
	}
	return b.topLeftX, b.topLeftY, b.getStaticData().w, h
}

func (b *building) getUnitPlacementAbsoluteCoords() (int, int) {
	return b.topLeftX + b.getStaticData().unitPlacementX, b.topLeftY + b.getStaticData().unitPlacementY
}

func (b *building) getName() string {
	return b.getStaticData().displayedName
}

func (b *building) getCurrentAction() *action {
	return &b.currentAction
}

func (b *building) getCurrentOrder() *order {
	return &b.currentOrder
}

func (b *building) getFaction() *faction {
	return b.faction
}

func (b *building) getPhysicalCenterCoords() (float64, float64) {
	return float64(b.topLeftX) + float64(b.getStaticData().w)/2, float64(b.topLeftY) + float64(b.getStaticData().h)/2
}

func (b *building) isPresentAt(tileX, tileY int) bool {
	w, h := b.getStaticData().w, b.getStaticData().h
	return geometry.AreCoordsInTileRect(tileX, tileY, b.topLeftX, b.topLeftY, w, h)
}

func (b *building) getStaticData() *buildingStatic {
	return sTableBuildings[b.code]
}

func (b *building) isInAir() bool {
	return false
}
