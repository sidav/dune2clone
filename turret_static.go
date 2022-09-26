package main

type turretStatic struct {
	spriteCode  string // empty means invisible turret
	rotateSpeed int

	fireRange, attackCooldown int
	fireSpreadDegrees         int
	shotRangeSpread           float64

	attacksLand, attacksAir bool

	firesProjectileOfCode int
	projectileDamage      int
	projectileDamageType  damageCode
}

const (
	TRT_NONE = iota
	TRT_TANK
	TRT_TANK2
	TRT_DEVASTATOR
	TRT_INFANTRY
	TRT_ROCKETINFANTRY
	TRT_HEAVYINFANTRY
	TRT_MSLTANK
	TRT_HARVESTER
	TRT_AATANK
	TRT_QUAD
	TRT_AIR_GUNSHIP
	TRT_CANNON_BUILDING
	TRT_MINIGUN_BUILDING
	TRT_BUILDING_FORTRESS
	TRT_BUILDING_AA
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
		projectileDamageType:  DAMAGETYPE_HEAVY,
	},
	TRT_TANK2: {
		spriteCode:            "tank2",
		firesProjectileOfCode: PRJ_SHELL,
		attacksLand:           true,
		rotateSpeed:           7,
		fireRange:             5,
		fireSpreadDegrees:     7,
		shotRangeSpread:       0.7,
		attackCooldown:        45,
		projectileDamage:      30,
		projectileDamageType:  DAMAGETYPE_HEAVY,
	},
	TRT_DEVASTATOR: {
		spriteCode:            "devastator",
		firesProjectileOfCode: PRJ_SHELL,
		attacksLand:           true,
		rotateSpeed:           0,
		fireRange:             6,
		fireSpreadDegrees:     7,
		shotRangeSpread:       0.5,
		attackCooldown:        75,
		projectileDamage:      47,
		projectileDamageType:  DAMAGETYPE_HEAVY,
	},
	TRT_INFANTRY: {
		spriteCode:            "",
		firesProjectileOfCode: PRJ_BULLETS,
		attacksLand:           true,
		rotateSpeed:           0,
		fireRange:             4,
		fireSpreadDegrees:     7,
		shotRangeSpread:       0.5,
		attackCooldown:        45,
		projectileDamage:      4,
		projectileDamageType:  DAMAGETYPE_ANTI_INFANTRY,
	},
	TRT_ROCKETINFANTRY: {
		spriteCode:            "",
		firesProjectileOfCode: PRJ_INFANTRY_MISSILE,
		attacksLand:           true,
		attacksAir:            true,
		rotateSpeed:           0,
		fireRange:             5,
		fireSpreadDegrees:     45,
		shotRangeSpread:       0.5,
		attackCooldown:        105,
		projectileDamage:      8,
		projectileDamageType:  DAMAGETYPE_HEAVY,
	},
	TRT_HEAVYINFANTRY: {
		spriteCode:            "",
		firesProjectileOfCode: PRJ_OMNI,
		attacksLand:           true,
		rotateSpeed:           0,
		fireRange:             5,
		fireSpreadDegrees:     8,
		shotRangeSpread:       0.45,
		attackCooldown:        40,
		projectileDamage:      6,
		projectileDamageType:  DAMAGETYPE_OMNI,
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
		projectileDamage:      50,
		projectileDamageType:  DAMAGETYPE_ANTI_BUILDING,
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
		projectileDamageType:  DAMAGETYPE_HEAVY,
	},
	TRT_HARVESTER: {
		spriteCode:            "",
		firesProjectileOfCode: PRJ_BULLETS,
		attacksLand:           true,
		rotateSpeed:           90,
		fireRange:             5,
		fireSpreadDegrees:     11,
		shotRangeSpread:       0.4,
		attackCooldown:        15,
		projectileDamage:      3,
		projectileDamageType:  DAMAGETYPE_ANTI_INFANTRY,
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
		projectileDamage:      4,
		projectileDamageType:  DAMAGETYPE_ANTI_INFANTRY,
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
		projectileDamage:      12,
		projectileDamageType:  DAMAGETYPE_HEAVY,
	},
	TRT_MINIGUN_BUILDING: {
		spriteCode:            "bld_turret_minigun",
		firesProjectileOfCode: PRJ_BULLETS,
		attacksLand:           true,
		attacksAir:            true,
		rotateSpeed:           17,
		fireRange:             6,
		fireSpreadDegrees:     7,
		shotRangeSpread:       0.7,
		attackCooldown:        5,
		projectileDamage:      4,
		projectileDamageType:  DAMAGETYPE_ANTI_INFANTRY,
	},
	TRT_CANNON_BUILDING: {
		spriteCode:            "bld_turret_cannon",
		firesProjectileOfCode: PRJ_SHELL,
		attacksLand:           true,
		rotateSpeed:           15,
		fireRange:             6,
		fireSpreadDegrees:     7,
		shotRangeSpread:       0.7,
		attackCooldown:        50,
		projectileDamage:      15,
		projectileDamageType:  DAMAGETYPE_HEAVY,
	},
	TRT_BUILDING_FORTRESS: {
		spriteCode:            "bld_fortress_cannon",
		firesProjectileOfCode: PRJ_SHELL,
		attacksLand:           true,
		rotateSpeed:           5,
		fireRange:             15,
		fireSpreadDegrees:     5,
		shotRangeSpread:       0.3,
		attackCooldown:        80,
		projectileDamage:      25,
		projectileDamageType:  DAMAGETYPE_HEAVY,
	},
	TRT_BUILDING_AA: {
		spriteCode:            "",
		firesProjectileOfCode: PRJ_AA_MISSILE,
		attacksLand:           false,
		attacksAir:            true,
		rotateSpeed:           180,
		fireRange:             8,
		fireSpreadDegrees:     30,
		shotRangeSpread:       0.3,
		attackCooldown:        100,
		projectileDamage:      25,
		projectileDamageType:  DAMAGETYPE_HEAVY,
	},
}
