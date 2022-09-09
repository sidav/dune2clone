package main

import (
	"dune2clone/geometry"
)

type building struct {
	currentAction          action
	currentOrder           order
	currentHitpoints       int
	topLeftX, topLeftY     int // tile coords
	code                   int
	faction                *faction
	isSelected             bool
	turret                 *turret
	unitPlacedInside       *unit
	rallyTileX, rallytileY int
}

func (b *building) isAlive() bool {
	return b.currentHitpoints > 0
}

func createBuilding(code, topLeftX, topLeftY int, fact *faction) *building {
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

func (b *building) getVisionRange() int {
	if b.faction.lastAvailableEnergy < 0 {
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

func (b *building) getUnitPlacementCoords() (int, int) {
	return b.topLeftX + b.getStaticData().unitPlacementX, b.topLeftY + b.getStaticData().unitPlacementY
}

func (b *building) getName() string {
	return b.getStaticData().displayedName
}

func (b *building) getCurrentAction() *action {
	return &b.currentAction
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
