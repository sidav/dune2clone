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

	var hitTarget actor
	if p.targetActor != nil && p.getStaticData().rotationSpeed > 0 {
		targX, targY := p.targetActor.getPhysicalCenterCoords()
		rotateTo := geometry.GetDegreeOfFloatVector(targX-p.centerX, targY-p.centerY)
		p.rotationDegree += geometry.GetDiffForRotationStep(p.rotationDegree, rotateTo, p.getStaticData().rotationSpeed)
		p.rotationDegree = geometry.NormalizeDegree(p.rotationDegree)
		if geometry.GetApproxDistFloat64(targX, targY, p.centerX, p.centerY) < 1 {
			hitTarget = p.targetActor
			p.setToRemove = true
		}
	}
	if p.fuel <= 0 && hitTarget == nil {
		tilex, tiley := geometry.TrueCoordsToTileCoords(p.centerX, p.centerY)
		hitTarget = b.getActorAtTileCoordinates(tilex, tiley)
		if hitTarget != nil && hitTarget.isInAir() != p.targetActor.isInAir() {
			hitTarget = nil
		}
		p.setToRemove = true
	}
	if p.setToRemove {
		if hitTarget != nil {
			b.dealDamageToActor(p.damage, p.damageType, p.targetActor)
		}
		b.dealSplashDamage(p.centerX, p.centerY, p.getStaticData().splashRadius, p.splashDamage, p.damageType)
		if p.getStaticData().createsEffectOnImpact {
			b.addEffect(&effect{
				centerX:            p.centerX,
				centerY:            p.centerY,
				splashCircleRadius: p.getStaticData().splashRadius,
				code:               p.getStaticData().effectCreatedOnImpactCode,
				creationTick:       b.currentTick,
			})
		}
	}
	// debugWritef("%+v spd: %f\n", p, spd)
}
