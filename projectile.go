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
	return p.getStaticData().rotationSpeed > 0
}

type projectileStatic struct {
	spriteCode                string
	size                      float64
	speed                     float64
	splashRadius              float64
	createsEffectOnImpact     bool
	effectCreatedOnImpactCode effectCode
	rotationSpeed             int

	hitDamage    int
	splashDamage int
	damageType   damageCode
}
