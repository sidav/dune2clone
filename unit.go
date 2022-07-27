package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

type unit struct {
	code             int
	centerX, centerY float64
	faction          *faction
	currentAction    action
	currentOrder     order
	chassisDegree    int
	cannonDegree     int

	isSelected bool // for rendering selection thingy
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
	if u.getStaticData().cannonRotationSpeed > 0 {
		return []rl.Texture2D{
			unitChassisAtlaces[sTableUnits[u.code].chassisSpriteCode].atlas[degreeToRotationFrameNumber(u.chassisDegree, 8)][0],
			unitCannonsAtlaces[sTableUnits[u.code].cannonSpriteCode].atlas[degreeToRotationFrameNumber(u.cannonDegree, 8)][0],
		}
	}
	return []rl.Texture2D{
		unitChassisAtlaces[sTableUnits[u.code].chassisSpriteCode].atlas[degreeToRotationFrameNumber(u.chassisDegree, 8)][0],
	}
}

func (u *unit) normalizeDegrees() {
	if u.cannonDegree < 0 {
		u.cannonDegree += 360
	}
	if u.cannonDegree > 360 {
		u.cannonDegree -= 360
	}
	if u.chassisDegree < 0 {
		u.chassisDegree += 360
	}
	if u.chassisDegree > 360 {
		u.chassisDegree -= 360
	}
}

func (u *unit) rotateChassisTowardsVector(vx, vy float64) {
	if isVectorDegreeEqualTo(vx, vy, u.chassisDegree) {
		return
	}
	degs := int(180 * math.Atan2(vy, vx) / 3.14159265358)
	if degs < 0 {
		degs += 360
	}
	diff := u.chassisDegree - degs
	for diff < 0 {
		diff += 360
	}
	rotateSpeed := u.getStaticData().chassisRotationSpeed
	if rotateSpeed > diff {
		rotateSpeed = diff
	}
	if diff <= 180 {
		rotateSpeed = -rotateSpeed
	}

	// debugWritef("targetdegs %d, unitdegs %d, diff %d, rotateSpeed %d\n", degs, u.chassisDegree, u.cannonDegree, )
	u.chassisDegree += rotateSpeed
	u.cannonDegree += rotateSpeed
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
	displayedName string

	cannonSpriteCode  string
	chassisSpriteCode string

	movementSpeed        float64
	chassisRotationSpeed int
	cannonRotationSpeed  int // 0 means that the unit doesn't have separate cannon

	cost          int
	buildTime     int // seconds
	hotkeyToBuild string
}

var sTableUnits = map[int]*unitStatic{
	UNT_QUAD: {
		displayedName:    "Quad",
		// cannonSpriteCode: "quad",
		chassisSpriteCode: "quad",
		movementSpeed:        0.25,
		chassisRotationSpeed: 7,
		cannonRotationSpeed:  0,
		cost:                 350,
		buildTime:            3,
		hotkeyToBuild:        "Q",
	},
	UNT_TANK: {
		displayedName:        "Super duper tank",
		cannonSpriteCode:     "tank",
		chassisSpriteCode:    "tank",
		movementSpeed:        0.1,
		chassisRotationSpeed: 5,
		cannonRotationSpeed:  7,
		cost:                 450,
		buildTime:            7,
		hotkeyToBuild:        "T",
	},
}
