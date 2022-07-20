package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

type unit struct {
	code             int
	centerX, centerY float64
	currentAction    *action
	chassisDegree    int
	cannonDegree     int
}

func (u *unit) getPartsSprites() []rl.Texture2D {
	chassisSprite := unitChassisAtlaces[sTableUnits[u.code].chassisSpriteCode].atlas[degreeToRotationFrameNumber(u.chassisDegree)][0]
	cannonSprite := unitCannonsAtlaces[sTableUnits[u.code].cannonSpriteCode].atlas[degreeToRotationFrameNumber(u.cannonDegree)][0]
	return []rl.Texture2D{
		chassisSprite,
		cannonSprite,
	}
}

func (u *unit) rotateChassisTowardsVector(vx, vy float64) bool {
	degs := int(180 * math.Atan2(vy, vx) / 3.14159265358)
	// debugWritef("degs %d, unitdegs %d\n", degs, u.chassisDegree)
	if u.chassisDegree > degs {
		u.chassisDegree -= u.getStaticData().rotationSpeed
		u.cannonDegree -= u.getStaticData().rotationSpeed
		if u.chassisDegree < degs {
			u.chassisDegree = degs
		}
	} else if u.chassisDegree < degs {
		u.chassisDegree += u.getStaticData().rotationSpeed
		u.cannonDegree += u.getStaticData().rotationSpeed
		if u.chassisDegree > degs {
			u.chassisDegree = degs
		}
	}
	if u.chassisDegree == degs {
		return true
	}
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
