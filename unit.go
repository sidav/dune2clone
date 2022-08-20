package main

import (
	"dune2clone/geometry"
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

	currentCargoAmount int // for harvesters. TODO: separate struct?

	carriedUnit *unit // for transports

	isSelected bool // for rendering selection thingy
}

func createUnit(code, tx, ty int, fact *faction) *unit {
	cx, cy := geometry.TileCoordsToPhysicalCoords(tx, ty)
	u := &unit{
		code:             code,
		centerX:          cx,
		centerY:          cy,
		currentHitpoints: sTableUnits[code].maxHitpoints,
		faction:          fact,
		chassisDegree:    270,
	}
	if sTableUnits[code].turretsData != nil {
		for i := range sTableUnits[code].turretsData {
			u.turrets = append(u.turrets, &turret{code: sTableUnits[code].turretsData[i].turretCode, rotationDegree: 270})
		}
	}

	u.currentOrder.code = u.getStaticData().defaultOrderOnCreation
	u.currentOrder.targetTileX = -1
	u.currentOrder.targetTileY = -1
	u.currentAction.reset()
	return u
}

func (u *unit) markSelected(b bool) {
	u.isSelected = b
}

func (u *unit) getPhysicalCenterCoords() (float64, float64) {
	return u.centerX, u.centerY
}

func (u *unit) getVisionRange() int {
	vr := u.getStaticData().visionRange
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
	return u.getStaticData().displayedName
}

func (u *unit) getCurrentAction() *action {
	return &u.currentAction
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
	rotateSpeed := geometry.GetDiffForRotationStep(u.chassisDegree, deg, u.getStaticData().chassisRotationSpeed)
	u.chassisDegree += rotateSpeed
	for i := range u.turrets {
		u.turrets[i].rotationDegree += rotateSpeed
	}
	u.normalizeDegrees()
	if u.carriedUnit != nil {
		u.carriedUnit.chassisDegree = u.chassisDegree
	}
}

func (u *unit) getStaticData() *unitStatic {
	return sTableUnits[u.code]
}
