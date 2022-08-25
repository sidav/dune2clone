package main

import (
	"dune2clone/geometry"
	"math"
)

func (b *battlefield) actForProjectile(p *projectile) {
	if p.setToRemove {
		return // workaround for emptying the list
	}
	// move forward
	vx, vy := geometry.DegreeToUnitVector(p.rotationDegree)
	spd := math.Min(p.getStaticData().speed, p.fuel)
	p.centerX += spd * vx
	p.centerY += spd * vy
	p.fuel -= spd
	if p.targetActor != nil && p.getStaticData().rotationSpeed > 0 {
		targX, targY := p.targetActor.getPhysicalCenterCoords()
		rotateTo := geometry.GetDegreeOfFloatVector(targX-p.centerX, targY-p.centerY)
		p.rotationDegree += geometry.GetDiffForRotationStep(p.rotationDegree, rotateTo, p.getStaticData().rotationSpeed)
		p.rotationDegree = geometry.NormalizeDegree(p.rotationDegree)
		if geometry.GetApproxDistFloat64(targX, targY, p.centerX, p.centerY) < 1 {
			b.dealDamageToActor(p.damage, p.targetActor)
			p.setToRemove = true
			if p.getStaticData().createsEffectOnImpact {
				b.addEffect(&effect{
					centerX:      p.centerX,
					centerY:      p.centerY,
					code:         p.getStaticData().effectCreatedOnImpactCode,
					creationTick: b.currentTick,
				})
			}
			return
		}
	}
	if p.fuel <= 0 {
		tilex, tiley := geometry.TrueCoordsToTileCoords(p.centerX, p.centerY)
		targ := b.getActorAtTileCoordinates(tilex, tiley)
		if targ != nil && (p.targetActor == nil || targ.isInAir() == p.targetActor.isInAir()) {
			b.dealDamageToActor(p.damage, targ)
		}
		if p.getStaticData().createsEffectOnImpact {
			b.addEffect(&effect{
				centerX:      p.centerX,
				centerY:      p.centerY,
				code:         p.getStaticData().effectCreatedOnImpactCode,
				creationTick: b.currentTick,
			})
		}
		p.setToRemove = true
	}
	// debugWritef("%+v spd: %f\n", p, spd)
}
