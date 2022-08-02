package main

func (b *battlefield) actForProjectile(p *projectile) {
	// move forward
	vx, vy := degreeToUnitVector(p.rotationDegree)
	spd := p.getStaticData().speed
	p.centerX += spd * vx
	p.centerY += spd * vy
	p.fuel -= spd
}
