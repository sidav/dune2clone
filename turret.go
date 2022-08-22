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

type turretStatic struct {
	spriteCode  string // empty means invisible turret
	rotateSpeed int

	fireRange, attackCooldown int
	fireSpreadDegrees         int
	shotRangeSpread           float64

	attacksLand, attacksAir bool

	firesProjectileOfCode int
	projectileDamage      int
}

const (
	TRT_NONE = iota
	TRT_TANK
	TRT_MSLTANK
	TRT_AATANK
	TRT_QUAD
	TRT_AIR_GUNSHIP
	TRT_CANNON_BUILDING
	TRT_MINIGUN_BUILDING
	TRT_BUILDING_FORTRESS
)

var sTableTurrets = map[int]*turretStatic{
	TRT_TANK: {
		spriteCode:            "tank",
		firesProjectileOfCode: PRJ_SHELL,
		attacksLand:           true,
		rotateSpeed:           7,
		fireRange:             5,
		fireSpreadDegrees:     7,
		shotRangeSpread:       0.7,
		attackCooldown:        45,
		projectileDamage:      30,
	},
	TRT_MSLTANK: {
		spriteCode:            "msltank",
		firesProjectileOfCode: PRJ_MISSILE,
		attacksLand:           true,
		rotateSpeed:           15,
		fireRange:             10,
		fireSpreadDegrees:     35,
		shotRangeSpread:       0.7,
		attackCooldown:        150,
		projectileDamage:      45,
	},
	TRT_AATANK: {
		spriteCode:            "aamsltank",
		firesProjectileOfCode: PRJ_AA_MISSILE,
		attacksAir:            true,
		rotateSpeed:           15,
		fireRange:             10,
		fireSpreadDegrees:     35,
		shotRangeSpread:       0.7,
		attackCooldown:        75,
		projectileDamage:      45,
	},
	TRT_QUAD: {
		spriteCode:            "",
		firesProjectileOfCode: PRJ_BULLETS,
		attacksLand:           true,
		rotateSpeed:           0,
		fireRange:             4,
		fireSpreadDegrees:     6,
		shotRangeSpread:       0.3,
		attackCooldown:        5,
		projectileDamage:      1,
	},
	TRT_AIR_GUNSHIP: {
		spriteCode:            "",
		firesProjectileOfCode: PRJ_SHELL,
		attacksLand:           true,
		rotateSpeed:           180,
		fireRange:             6,
		fireSpreadDegrees:     15,
		shotRangeSpread:       2.0,
		attackCooldown:        15,
		projectileDamage:      10,
	},
	TRT_MINIGUN_BUILDING: {
		spriteCode:            "bld_turret_minigun",
		firesProjectileOfCode: PRJ_BULLETS,
		attacksLand:           true,
		attacksAir:            true,
		rotateSpeed:           15,
		fireRange:             6,
		fireSpreadDegrees:     7,
		shotRangeSpread:       0.7,
		attackCooldown:        5,
		projectileDamage:      2,
	},
	TRT_CANNON_BUILDING: {
		spriteCode:            "bld_turret_cannon",
		firesProjectileOfCode: PRJ_SHELL,
		attacksLand:           true,
		rotateSpeed:           5,
		fireRange:             6,
		fireSpreadDegrees:     7,
		shotRangeSpread:       0.7,
		attackCooldown:        50,
		projectileDamage:      15,
	},
	TRT_BUILDING_FORTRESS: {
		spriteCode:            "bld_fortress_cannon",
		firesProjectileOfCode: PRJ_SHELL,
		attacksLand:           true,
		rotateSpeed:           5,
		fireRange:             15,
		fireSpreadDegrees:     5,
		shotRangeSpread:       0.3,
		attackCooldown:        50,
		projectileDamage:      25,
	},
}
