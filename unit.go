package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

type unit struct {
	code             int
	centerX, centerY float64
	faction          *faction
	turret           turret // maybe turrets array?..
	currentAction    action
	currentOrder     order
	chassisDegree    int

	isSelected bool // for rendering selection thingy
}

func createUnit(code, tx, ty int, fact *faction) *unit {
	cx, cy := tileCoordsToPhysicalCoords(tx, ty)
	return &unit{
		code:          code,
		centerX:       cx,
		centerY:       cy,
		faction:       fact,
		turret:        turret{code: sTableUnits[code].turretCode, rotationDegree: 270},
		chassisDegree: 270,
	}
}

func (u *unit) markSelected(b bool) {
	u.isSelected = b
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

func (u *unit) getPartsSprites() []rl.Texture2D {
	if u.turret.canRotate() {
		return []rl.Texture2D{
			unitChassisAtlaces[sTableUnits[u.code].chassisSpriteCode].atlas[degreeToRotationFrameNumber(u.chassisDegree, 8)][0],
			unitCannonsAtlaces[u.turret.getStaticData().spriteCode].atlas[degreeToRotationFrameNumber(u.turret.rotationDegree, 8)][0],
		}
	}
	return []rl.Texture2D{
		unitChassisAtlaces[sTableUnits[u.code].chassisSpriteCode].atlas[degreeToRotationFrameNumber(u.chassisDegree, 8)][0],
	}
}

func (u *unit) normalizeDegrees() {
	if u.chassisDegree < 0 {
		u.chassisDegree += 360
	}
	if u.chassisDegree >= 360 {
		u.chassisDegree -= 360
	}
}

func (u *unit) rotateChassisTowardsVector(vx, vy float64) {
	if isVectorDegreeEqualTo(vx, vy, u.chassisDegree) {
		return
	}
	degs := int(180 * math.Atan2(vy, vx) / 3.14159265358)
	rotateSpeed := getDiffForRotationStep(u.chassisDegree, degs, u.getStaticData().chassisRotationSpeed)

	// debugWritef("targetdegs %d, unitdegs %d, diff %d, rotateSpeed %d\n", degs, u.chassisDegree, u.cannonDegree, )
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
)

type unitStatic struct {
	displayedName     string
	chassisSpriteCode string

	turretCode int

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
		turretCode:           TRT_TANK,
		chassisRotationSpeed: 5,
		cost:                 450,
		buildTime:            7,
		hotkeyToBuild:        "T",
	},
}
