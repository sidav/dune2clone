package main

import (
	"dune2clone/geometry"
	"math"
)

type unit struct {
	code             int
	centerX, centerY float64
	faction          *faction
	turrets          []*turret // maybe turrets array?..
	currentAction    action
	currentOrder     order
	currentHitpoints int
	chassisDegree    int

	squadSize int // > 0 if unit is a squad

	currentCargoAmount int // for harvesters. TODO: separate struct?

	carriedUnit *unit // for transports

	isSelected bool // for rendering selection thingy
	experience int
}

func createUnit(code, tx, ty int, fact *faction) *unit {
	cx, cy := geometry.TileCoordsToTrueCoords(tx, ty)
	u := &unit{
		code:             code,
		centerX:          cx,
		centerY:          cy,
		currentHitpoints: sTableUnits[code].MaxHitpoints,
		squadSize:        sTableUnits[code].MaxSquadSize,
		faction:          fact,
		chassisDegree:    270,
	}
	if sTableUnits[code].TurretsData != nil {
		for i := range sTableUnits[code].TurretsData {
			u.turrets = append(u.turrets, &turret{staticData: sTableUnits[code].TurretsData[i], rotationDegree: 270})
		}
	}

	u.currentOrder.code = u.getStaticData().DefaultOrderOnCreation
	u.currentOrder.targetTileX = -1
	u.currentOrder.targetTileY = -1
	u.currentAction.resetAction()
	return u
}

func (u *unit) receiveExperienceAmount(amnt int) {
	u.experience += amnt
	if u.getStaticData().HasEliteVersion && u.getExperienceLevel() == MAX_VETERANCY_LEVEL {
		u.convertToElite()
	}
}

func (u *unit) getRegenAmount() int {
	return u.getStaticData().HpRegen + getVeterancyBasedRegen(u.getExperienceLevel())
}

func (u *unit) setHitpoints(hp int) {
	u.currentHitpoints = hp
}

func (u *unit) convertToElite() {
	u.code = u.getStaticData().EliteVersionCode
	u.turrets = []*turret{}
	if u.getStaticData().TurretsData != nil {
		for i := range u.getStaticData().TurretsData {
			u.turrets = append(u.turrets, &turret{staticData: u.getStaticData().TurretsData[i], rotationDegree: 270})
		}
	}
}

func (u *unit) getExperienceLevel() int {
	return getExperienceLevelByAmountAndCost(u.experience, u.getStaticData().Cost)
}

func (u *unit) isAlive() bool {
	return u.currentHitpoints > 0
}

func (u *unit) receiveHealing(amount int) {
	u.currentHitpoints += amount
	if u.currentHitpoints > u.getMaxHitpoints() {
		u.currentHitpoints = u.getMaxHitpoints()
	}
}

func (u *unit) getExperience() int {
	return u.experience
}

func (u *unit) recalculateSquadSize() {
	if u.getStaticData().MaxSquadSize > 1 {
		u.squadSize = int(
			math.Ceil(float64(u.getStaticData().MaxSquadSize) *
				float64(u.currentHitpoints) / float64(u.getMaxHitpoints())),
		)
	}
}

func (u *unit) getHitpoints() int {
	return u.currentHitpoints
}

func (u *unit) getMaxHitpoints() int {
	return modifyMaxHpByExpLevel(u.getStaticData().MaxHitpoints, u.getExperienceLevel())
}

func (u *unit) getMovementSpeed() float64 {
	return modifyUnitSpeedByLevel(u.getStaticData().MovementSpeed, u.getExperienceLevel())
}

func (u *unit) getHitpointsPercentage() int {
	return getPercentInt(u.currentHitpoints, u.getMaxHitpoints())
}

func (u *unit) markSelected(b bool) {
	u.isSelected = b
}

func (u *unit) getPhysicalCenterCoords() (float64, float64) {
	return u.centerX, u.centerY
}

func (u *unit) getTileCoords() (int, int) {
	return geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
}

func (u *unit) getVisionRange() int {
	vr := u.getStaticData().VisionRange
	if vr <= 0 {
		vr = 1
	}
	return vr
}

func (u *unit) setPhysicalCenterCoords(x, y float64) {
	u.centerX = x
	u.centerY = y
	if u.carriedUnit != nil {
		u.carriedUnit.centerX = x
		u.carriedUnit.centerY = y
	}
}

func (u *unit) isPresentAt(tileX, tileY int) bool {
	tx, ty := geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
	return tx == tileX && ty == tileY
}

func (u *unit) getName() string {
	return u.getStaticData().DisplayedName
}

func (u *unit) getCurrentAction() *action {
	return &u.currentAction
}

func (u *unit) isInTileCenter() bool {
	_, fracX := math.Modf(u.centerX)
	_, fracY := math.Modf(u.centerY)
	return areFloatsRoughlyEqual(fracX, 0.5) && areFloatsRoughlyEqual(fracY, 0.5)
}

func (u *unit) getCurrentOrder() *order {
	return &u.currentOrder
}

func (u *unit) getFaction() *faction {
	return u.faction
}

func (u *unit) normalizeDegrees() {
	u.chassisDegree = geometry.NormalizeDegree(u.chassisDegree)
	for i := range u.turrets {
		u.turrets[i].normalizeDegrees()
	}
}

func (u *unit) rotateChassisTowardsVector(vx, vy float64) {
	degs := geometry.GetDegreeOfFloatVector(vx, vy)
	u.rotateChassisTowardsDegree(degs)
}

func (u *unit) rotateChassisTowardsDegree(deg int) {
	if u.chassisDegree == deg {
		return
	}
	rotateSpeed := geometry.GetDiffForRotationStep(u.chassisDegree, deg, u.getStaticData().ChassisRotationSpeed)
	u.chassisDegree += rotateSpeed
	for i := range u.turrets {
		u.turrets[i].rotationDegree += rotateSpeed
	}
	u.normalizeDegrees()
	if u.carriedUnit != nil {
		u.carriedUnit.chassisDegree = u.chassisDegree
		for _, t := range u.carriedUnit.turrets {
			t.rotationDegree = u.chassisDegree
		}
	}
}

func (u *unit) getStaticData() *unitStatic {
	return sTableUnits[u.code]
}

func (u *unit) getMainTurretRange() int {
	if len(u.turrets) == 0 {
		return 0
	}
	return u.turrets[0].getStaticData().FireRange
}

func (u *unit) isInAir() bool {
	return u.getStaticData().IsAircraft
}
