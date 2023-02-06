package main

type projectile struct {
	faction          *faction
	staticData       *projectileStatic
	centerX, centerY float64
	rotationDegree   int
	fuel             float64 // how many 'speeds' it spends until it is destroyed
	whoShot          actor
	targetActor      actor // for homing projectiles
	setToRemove      bool
}

func (p *projectile) getStaticData() *projectileStatic {
	return p.staticData
}

func (p *projectile) isHoming() bool {
	return p.getStaticData().RotationSpeed > 0
}

type projectileStatic struct {
	SpriteCode                string     `json:"sprite_code,omitempty"`
	Size                      float64    `json:"size,omitempty"`
	Speed                     float64    `json:"speed,omitempty"`
	SplashRadius              float64    `json:"splash_radius,omitempty"`
	CreatesEffectOnImpact     bool       `json:"creates_effect_on_impact,omitempty"`
	EffectCreatedOnImpactCode effectCode `json:"effect_created_on_impact_code,omitempty"`
	RotationSpeed             int        `json:"rotation_speed,omitempty"`

	HitDamage    int        `json:"hit_damage,omitempty"`
	SplashDamage int        `json:"splash_damage,omitempty"`
	DamageType   damageCode `json:"damage_type,omitempty"`
}
