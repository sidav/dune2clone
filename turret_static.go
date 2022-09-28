package main

type turretStatic struct {
	spriteCode  string // empty means invisible turret
	rotateSpeed int

	turretCenterX, turretCenterY float64 // relative to unit's center

	fireRange, attackCooldown int
	fireSpreadDegrees         int
	shotRangeSpread           float64

	attacksLand, attacksAir bool

	firedProjectileData *projectileStatic
}
