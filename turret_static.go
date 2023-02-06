package main

type TurretStatic struct {
	SpriteCode  string `json:"sprite_code,omitempty"` // empty means invisible turret
	RotateSpeed int    `json:"rotate_speed,omitempty"`

	TurretCenterX float64 `json:"turret_center_x"`
	TurretCenterY float64 `json:"turret_center_y"` // relative to unit's center

	FireRange         int     `json:"fire_range"`
	AttackCooldown    int     `json:"attack_cooldown"`
	FireSpreadDegrees int     `json:"fire_spread_degrees,omitempty"`
	ShotRangeSpread   float64 `json:"shot_range_spread,omitempty"`

	AttacksLand bool `json:"attacks_land"`
	AttacksAir  bool `json:"attacks_air"`

	FiredProjectileData *projectileStatic `json:"fired_projectile_data,omitempty"`
}
