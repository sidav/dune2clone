package main

import "dune2clone/geometry"

type turret struct {
	code           int
	rotationDegree int
	nextTickToAct  int

	targetActor              actor
	targetTileX, targetTileY int
}

func (t *turret) canRotate() bool {
	return t.getStaticData().rotateSpeed > 0
}

func (t *turret) getStaticData() *turretStatic {
	return sTableTurrets[t.code]
}

func (t *turret) normalizeDegrees() {
	t.rotationDegree = geometry.NormalizeDegree(t.rotationDegree)
}

const (
	TRT_NONE = iota
	TRT_TANK
	TRT_QUAD
	TRT_CANNON_BUILDING
)

type turretStatic struct {
	spriteCode                string
	rotateSpeed               int
	fireRange, attackCooldown int

	fireSpreadDegrees int
	shotRangeSpread   float64
}

var sTableTurrets = map[int]*turretStatic{
	TRT_TANK: {
		spriteCode:        "tank",
		rotateSpeed:       7,
		fireRange:         5,
		fireSpreadDegrees: 7,
		shotRangeSpread:   0.5,
		attackCooldown:    25,
	},
	TRT_QUAD: {
		spriteCode:        "",
		rotateSpeed:       0,
		fireRange:         4,
		fireSpreadDegrees: 7,
		shotRangeSpread:   0.3,
		attackCooldown:    25,
	},
	TRT_CANNON_BUILDING: {
		spriteCode:        "tank",
		rotateSpeed:       3,
		fireRange:         7,
		fireSpreadDegrees: 7,
		shotRangeSpread:   0.5,
		attackCooldown:    50,
	},
}
