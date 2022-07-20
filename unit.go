package main

import rl "github.com/gen2brain/raylib-go/raylib"

type unit struct {
	code int
	x, y float64
}

func (u *unit) getPartsSprites() []rl.Texture2D {
	chassisSprite := unitChassisAtlaces[sTableUnits[u.code].chassisSpriteCode].atlas[0][0]
	cannonSprite := unitCannonsAtlaces[sTableUnits[u.code].cannonSpriteCode].atlas[0][0]
	return []rl.Texture2D{
		chassisSprite,
		cannonSprite,
	}
}

const (
	UNT_TANK = iota
)

type unitStatic struct {
	cannonSpriteCode  string
	chassisSpriteCode string
}

var sTableUnits = map[int]unitStatic{
	UNT_TANK: {
		cannonSpriteCode: "tank",
		chassisSpriteCode: "tank",
	},
}
