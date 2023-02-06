package main

type TurretStatic struct {
	SpriteCode  string `json:"sprite_code,omitempty"` // empty means invisible turret
	RotateSpeed int    `json:"rotate_speed,omitempty"`

	TurretCenterX, TurretCenterY float64 // relative to unit's center

	FireRange, AttackCooldown int
	FireSpreadDegrees         int     `json:"fire_spread_degrees,omitempty"`
	ShotRangeSpread           float64 `json:"shot_range_spread,omitempty"`

	AttacksLand, AttacksAir bool

	FiredProjectileData *projectileStatic `json:"fired_projectile_data,omitempty"`
}
