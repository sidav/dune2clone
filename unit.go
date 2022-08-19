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

	currentCargoAmount int // for harvesters. TODO: separate struct?

	carriedUnit *unit // for transports

	isSelected bool // for rendering selection thingy
}

func createUnit(code, tx, ty int, fact *faction) *unit {
	cx, cy := geometry.TileCoordsToPhysicalCoords(tx, ty)
	var turr *turret
	if sTableUnits[code].turretCode != TRT_NONE {
		turr = &turret{code: sTableUnits[code].turretCode, rotationDegree: 270}
	}
	u := &unit{
		code:             code,
		centerX:          cx,
		centerY:          cy,
		currentHitpoints: sTableUnits[code].maxHitpoints,
		faction:          fact,
		turret:           turr,
		chassisDegree:    270,
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
	if u.turret != nil {
		u.turret.normalizeDegrees()
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
	if u.turret != nil {
		u.turret.rotationDegree += rotateSpeed
	}
	u.normalizeDegrees()
	if u.carriedUnit != nil {
		u.carriedUnit.chassisDegree = u.chassisDegree
	}
}

func (u *unit) getStaticData() *unitStatic {
	return sTableUnits[u.code]
}

const (
	UNT_TANK = iota
	UNT_QUAD
	UNT_MSLTANK
	UNT_HARVESTER

	// aircrafts
	AIR_TRANSPORT
	AIR_GUNSHIP
)

type unitStatic struct {
	displayedName     string
	chassisSpriteCode string

	turretCode int

	maxHitpoints int

	movementSpeed        float64
	chassisRotationSpeed int

	maxCargoAmount int // for harvesters

	defaultOrderOnCreation orderCode

	isAircraft  bool
	isTransport bool

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
	UNT_HARVESTER: {
		displayedName:          "Harvester",
		chassisSpriteCode:      "harvester",
		defaultOrderOnCreation: ORDER_HARVEST,
		maxCargoAmount:         700,
		movementSpeed:          0.07,
		maxHitpoints:           250,
		turretCode:             TRT_NONE,
		chassisRotationSpeed:   4,
		cost:                   1600,
		buildTime:              12,
		hotkeyToBuild:          "H",
	},
	// aircrafts
	AIR_TRANSPORT: {
		displayedName:        "Carrier aircraft",
		chassisSpriteCode:    "air_transport",
		maxHitpoints:         100,
		movementSpeed:        0.2,
		turretCode:           TRT_NONE,
		chassisRotationSpeed: 5,
		cost:                 500,
		buildTime:            1,
		hotkeyToBuild:        "C",
		isAircraft:           true,
		isTransport:          true,
	},
	AIR_GUNSHIP: {
		displayedName:        "Gunship",
		chassisSpriteCode:    "air_gunship",
		maxHitpoints:         50,
		movementSpeed:        0.25,
		turretCode:           TRT_AIR_GUNSHIP,
		chassisRotationSpeed: 3,
		cost:                 500,
		buildTime:            1,
		hotkeyToBuild:        "G",
		isAircraft:           true,
	},
}
