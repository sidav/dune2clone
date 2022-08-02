package main

import "dune2clone/geometry"

func (b *battlefield) actForProjectile(p *projectile) {
	// move forward
	vx, vy := geometry.DegreeToUnitVector(p.rotationDegree)
	spd := p.getStaticData().speed
	p.centerX += spd * vx
	p.centerY += spd * vy
	p.fuel -= spd
}
