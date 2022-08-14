package main

import (
	"dune2clone/geometry"
)

type aircraft struct {
	code             int
	centerX, centerY float64
	faction          *faction
	turret           *turret // maybe turrets array?..
	currentAction    action
	currentOrder     order
	currentHitpoints int
	chassisDegree    int

	currentCargoAmount int // for harvesters. TODO: separate struct?

	isSelected bool // for rendering selection thingy
}

func createAircraft(code, tx, ty int, fact *faction) *aircraft {
	cx, cy := geometry.TileCoordsToPhysicalCoords(tx, ty)
	var turr *turret
	if sTableUnits[code].turretCode != TRT_NONE {
		turr = &turret{code: sTableUnits[code].turretCode, rotationDegree: 270}
	}
	a := &aircraft{
		code:             code,
		centerX:          cx,
		centerY:          cy,
		currentHitpoints: sTableUnits[code].maxHitpoints,
		faction:          fact,
		turret:           turr,
		chassisDegree:    270,
	}
	a.currentOrder.code = a.getStaticData().defaultOrderOnCreation
	a.currentOrder.targetTileX = -1
	a.currentOrder.targetTileY = -1
	return a
}

func (a *aircraft ) markSelected(b bool) {
	a.isSelected = b
}

func (a *aircraft ) getPhysicalCenterCoords() (float64, float64) {
	return a.centerX, a.centerY
}

//func (a *aircraft ) isPresentAt(tileX, tileY int) bool {
//	tx, ty := geometry.TrueCoordsToTileCoords(a.centerX, a.centerY)
//	return tx == tileX && ty == tileY
//}

func (a *aircraft) getName() string {
	return a.getStaticData().displayedName
}

func (a *aircraft ) getCurrentAction() *action {
	return &a.currentAction
}

func (a *aircraft ) getFaction() *faction {
	return a.faction
}

func (a *aircraft ) normalizeDegrees() {
	a.chassisDegree = geometry.NormalizeDegree(a.chassisDegree)
	if a.turret != nil {
		a.turret.normalizeDegrees()
	}
}

func (a *aircraft ) rotateChassisTowardsVector(vx, vy float64) {
	if geometry.IsVectorDegreeEqualTo(vx, vy, a.chassisDegree) {
		return
	}
	degs := geometry.GetDegreeOfFloatVector(vx, vy)
	rotateSpeed := geometry.GetDiffForRotationStep(a.chassisDegree, degs, a.getStaticData().chassisRotationSpeed)
	//debugWritef("targetdegs %d, unitdegs %d, diff %d, rotateSpeed %d\n", degs, a.chassisDegree)
	a.chassisDegree += rotateSpeed
	if a.turret != nil {
		a.turret.rotationDegree += rotateSpeed
	}
	a.normalizeDegrees()
}

func (a *aircraft ) getStaticData() *aircraftStatic {
	return sTableAircrafts[a.code]
}

const (
	AIR_TRANSPORT = iota
	AIR_COMBAT
)

type aircraftStatic struct {
	displayedName     string
	chassisSpriteCode string

	turretCode int

	maxHitpoints int

	movementSpeed        float64
	chassisRotationSpeed int

	maxCargoAmount int // for harvesters

	defaultOrderOnCreation orderCode

	cost          int
	buildTime     int // seconds
	hotkeyToBuild string
}

var sTableAircrafts = map[int]*aircraftStatic{
	AIR_TRANSPORT: {
		displayedName:        "Carrier aircraft",
		chassisSpriteCode:    "transport",
		maxHitpoints:         100,
		movementSpeed:        0.5,
		turretCode:           TRT_NONE,
		chassisRotationSpeed: 30,
		cost:                 500,
		buildTime:            10,
		hotkeyToBuild:        "C",
	},
	//AIR_COMBAT: {
	//},
}
