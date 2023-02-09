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
	experience             int
}

func (b *building) isAlive() bool {
	return b.currentHitpoints > 0
}

func (b *building) getHitpoints() int {
	return b.currentHitpoints
}

func (b *building) setHitpoints(hp int) {
	b.currentHitpoints = hp
}

func (b *building) getMaxHitpoints() int {
	return modifyMaxHpByExpLevel(b.getStaticData().MaxHitpoints, b.getExperienceLevel())
}

func (b *building) getHitpointsPercentage() int {
	return getPercentInt(b.currentHitpoints, b.getMaxHitpoints())
}

func createBuilding(code buildingCode, topLeftX, topLeftY int, fact *faction) *building {
	var turr *turret
	if sTableBuildings[code].TurretData != nil {
		turr = &turret{staticData: sTableBuildings[code].TurretData, rotationDegree: 270}
	}
	return &building{
		code:             code,
		currentHitpoints: sTableBuildings[code].MaxHitpoints,
		topLeftX:         topLeftX,
		topLeftY:         topLeftY,
		faction:          fact,
		turret:           turr,
		rallyTileX:       -1,
		rallytileY:       -1,
	}
}

func (b *building) receiveExperienceAmount(amnt int) {
	b.experience += amnt
}

func (b *building) receiveHealing(amnt int) {
	b.currentHitpoints += amnt
	if b.currentHitpoints > b.getMaxHitpoints() {
		b.currentHitpoints = b.getMaxHitpoints()
	}
}

func (b *building) getRegenAmount() int {
	return getVeterancyBasedRegen(b.getExperienceLevel())
}

func (b *building) getExperience() int {
	return b.experience
}

func (b *building) getExperienceLevel() int {
	return getExperienceLevelByAmountAndCost(b.experience, b.getStaticData().Cost)
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
	h := b.getStaticData().H
	// prevent closing bottom side for producing buildings
	if b.getStaticData().NeedsEmptyRowBelowWhenConstructing {
		h++
	}
	return b.topLeftX, b.topLeftY, b.getStaticData().W, h
}

func (b *building) getUnitPlacementAbsoluteCoords() (int, int) {
	return b.topLeftX + b.getStaticData().UnitPlacementX, b.topLeftY + b.getStaticData().UnitPlacementY
}

func (b *building) getName() string {
	return b.getStaticData().DisplayedName
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
	return float64(b.topLeftX) + float64(b.getStaticData().W)/2, float64(b.topLeftY) + float64(b.getStaticData().H)/2
}

func (b *building) isPresentAt(tileX, tileY int) bool {
	w, h := b.getStaticData().W, b.getStaticData().H
	return geometry.AreCoordsInTileRect(tileX, tileY, b.topLeftX, b.topLeftY, w, h)
}

func (b *building) getStaticData() *buildingStatic {
	return sTableBuildings[b.code]
}

func (b *building) isInAir() bool {
	return false
}
