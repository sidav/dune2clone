package main

import rl "github.com/gen2brain/raylib-go/raylib"

type unit struct {
	code             int
	centerX, centerY float64
	currentAction    *action
}

func (u *unit) getPartsSprites() []rl.Texture2D {
	chassisSprite := unitChassisAtlaces[sTableUnits[u.code].chassisSpriteCode].atlas[0][0]
	cannonSprite := unitCannonsAtlaces[sTableUnits[u.code].cannonSpriteCode].atlas[0][0]
	return []rl.Texture2D{
		chassisSprite,
		cannonSprite,
	}
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

	speed float64
}

var sTableUnits = map[int]*unitStatic{
	UNT_TANK: {
		cannonSpriteCode:  "tank",
		chassisSpriteCode: "tank",
		speed: 0.1,
	},
}
