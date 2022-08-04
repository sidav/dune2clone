package main

import "dune2clone/geometry"

func (b *battlefield) actForProjectile(p *projectile) {
	// move forward
	vx, vy := geometry.DegreeToUnitVector(p.rotationDegree)
	spd := p.getStaticData().speed
	p.centerX += spd * vx
	p.centerY += spd * vy
	p.fuel -= spd
	if p.targetActor != nil && p.getStaticData().rotationSpeed > 0 {
		targX, targY := p.targetActor.getPhysicalCenterCoords()
		rotateTo := geometry.GetDegreeOfFloatVector(targX-p.centerX, targY-p.centerY)
		p.rotationDegree += geometry.GetDiffForRotationStep(p.rotationDegree, rotateTo, p.getStaticData().rotationSpeed)
		p.rotationDegree = geometry.NormalizeDegree(p.rotationDegree)
	}
	if p.fuel <= 0 {
		tilex, tiley := geometry.TrueCoordsToTileCoords(p.centerX, p.centerY)
		targ := b.getActorAtTileCoordinates(tilex, tiley)
		if targ != nil {
			b.dealDamageToActor(5, targ)
		}
	}
}
