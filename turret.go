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
	TRT_MSLTANK
	TRT_QUAD
	TRT_CANNON_BUILDING
)

type turretStatic struct {
	spriteCode                string
	firesProjectileOfCode     int
	rotateSpeed               int
	fireRange, attackCooldown int

	fireSpreadDegrees int
	shotRangeSpread   float64
}

var sTableTurrets = map[int]*turretStatic{
	TRT_TANK: {
		spriteCode:            "tank",
		firesProjectileOfCode: PRJ_SHELL,
		rotateSpeed:           7,
		fireRange:             5,
		fireSpreadDegrees:     7,
		shotRangeSpread:       0.7,
		attackCooldown:        45,
	},
	TRT_MSLTANK: {
		spriteCode:            "tank",
		firesProjectileOfCode: PRJ_MISSILE,
		rotateSpeed:           15,
		fireRange:             10,
		fireSpreadDegrees:     50,
		shotRangeSpread:       0.7,
		attackCooldown:        75,
	},
	TRT_QUAD: {
		spriteCode:            "",
		firesProjectileOfCode: PRJ_SHELL,
		rotateSpeed:           0,
		fireRange:             4,
		fireSpreadDegrees:     7,
		shotRangeSpread:       0.3,
		attackCooldown:        25,
	},
	TRT_CANNON_BUILDING: {
		spriteCode:            "cannon_turret",
		firesProjectileOfCode: PRJ_SHELL,
		rotateSpeed:           3,
		fireRange:             7,
		fireSpreadDegrees:     7,
		shotRangeSpread:       0.7,
		attackCooldown:        50,
	},
}
