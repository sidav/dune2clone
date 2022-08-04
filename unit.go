package main

import (
	"dune2clone/geometry"
)

type unit struct {
	code             int
	centerX, centerY float64
	faction          *faction
	turret           *turret // maybe turrets array?..
	currentAction    action
	currentOrder     order
	currentHitpoints int
	chassisDegree    int

	isSelected bool // for rendering selection thingy
}

func createUnit(code, tx, ty int, fact *faction) *unit {
	cx, cy := geometry.TileCoordsToPhysicalCoords(tx, ty)
	return &unit{
		code:             code,
		centerX:          cx,
		centerY:          cy,
		currentHitpoints: sTableUnits[code].maxHitpoints,
		faction:          fact,
		turret:           &turret{code: sTableUnits[code].turretCode, rotationDegree: 270},
		chassisDegree:    270,
	}
}

func (u *unit) markSelected(b bool) {
	u.isSelected = b
}

func (u *unit) getPhysicalCenterCoords() (float64, float64) {
	return u.centerX, u.centerY
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
	u.turret.normalizeDegrees()
}

func (u *unit) rotateChassisTowardsVector(vx, vy float64) {
	if geometry.IsVectorDegreeEqualTo(vx, vy, u.chassisDegree) {
		return
	}
	degs := geometry.GetDegreeOfFloatVector(vx, vy)
	rotateSpeed := geometry.GetDiffForRotationStep(u.chassisDegree, degs, u.getStaticData().chassisRotationSpeed)
	//debugWritef("targetdegs %d, unitdegs %d, diff %d, rotateSpeed %d\n", degs, u.chassisDegree)
	u.chassisDegree += rotateSpeed
	u.turret.rotationDegree += rotateSpeed
	u.normalizeDegrees()
}

func (u *unit) getStaticData() *unitStatic {
	return sTableUnits[u.code]
}

const (
	UNT_TANK = iota
	UNT_QUAD
	UNT_MSLTANK
)

type unitStatic struct {
	displayedName     string
	chassisSpriteCode string

	turretCode int

	maxHitpoints int

	movementSpeed        float64
	chassisRotationSpeed int

	cost          int
	buildTime     int // seconds
	hotkeyToBuild string
}

var sTableUnits = map[int]*unitStatic{
	UNT_QUAD: {
		displayedName:        "Quad",
		chassisSpriteCode:    "quad",
		maxHitpoints:         75,
		movementSpeed:        0.25,
		turretCode:           TRT_QUAD,
		chassisRotationSpeed: 7,
		cost:                 350,
		buildTime:            3,
		hotkeyToBuild:        "Q",
	},
	UNT_TANK: {
		displayedName:        "Super duper tank",
		chassisSpriteCode:    "tank",
		movementSpeed:        0.1,
		maxHitpoints:         120,
		turretCode:           TRT_TANK,
		chassisRotationSpeed: 5,
		cost:                 450,
		buildTime:            7,
		hotkeyToBuild:        "T",
	},
	UNT_MSLTANK: {
		displayedName:        "Missile tank",
		chassisSpriteCode:    "quad",
		movementSpeed:        0.05,
		maxHitpoints:         50,
		turretCode:           TRT_MSLTANK,
		chassisRotationSpeed: 8,
		cost:                 1150,
		buildTime:            12,
		hotkeyToBuild:        "M",
	},
}
