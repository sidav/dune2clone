package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

type unit struct {
	code             int
	centerX, centerY float64
	currentAction    action
	chassisDegree    int
	cannonDegree     int

	isSelected bool // for rendering selection thingy
}

func (u *unit) markSelected(b bool) {
	u.isSelected = b
}

func (u *unit) getPartsSprites() []rl.Texture2D {
	chassisSprite := unitChassisAtlaces[sTableUnits[u.code].chassisSpriteCode].atlas[degreeToRotationFrameNumber(u.chassisDegree, 8)][0]
	cannonSprite := unitCannonsAtlaces[sTableUnits[u.code].cannonSpriteCode].atlas[degreeToRotationFrameNumber(u.cannonDegree, 8)][0]
	return []rl.Texture2D{
		chassisSprite,
		cannonSprite,
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

func (u *unit) rotateChassisTowardsVector(vx, vy float64) bool {
	degs := int(180 * math.Atan2(vy, vx) / 3.14159265358)
	if degs < 0 {
		degs += 360
	}
	if u.chassisDegree == degs {
		return true
	}
	// debugWritef("targetdegs %d, unitdegs %d, cannondegs %d\n", degs, u.chassisDegree, u.cannonDegree)
	diff := u.chassisDegree-degs
	if diff < 0 {
		diff += 360
	}
	rotateSpeed := u.getStaticData().rotationSpeed
	if abs(diff-rotateSpeed) < 0 {
		rotateSpeed = diff
	}
	if diff >= 180 {
		u.chassisDegree += rotateSpeed
		u.cannonDegree += rotateSpeed
	} else {
		u.chassisDegree -= rotateSpeed
		u.cannonDegree -= rotateSpeed
	}
	u.normalizeDegrees()
	return false
}

func (u *unit) getStaticData() *unitStatic {
	return sTableUnits[u.code]
}

const (
	UNT_TANK = iota
)

type unitStatic struct {
	cannonSpriteCode  string
	chassisSpriteCode string

	movementSpeed float64
	rotationSpeed int
}

var sTableUnits = map[int]*unitStatic{
	UNT_TANK: {
		cannonSpriteCode:  "tank",
		chassisSpriteCode: "tank",
		movementSpeed:     0.1,
		rotationSpeed:     5,
	},
}
