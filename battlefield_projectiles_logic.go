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
		b.dealSplashDamage(p.centerX, p.centerY, p.getStaticData().splashRadius, p.getStaticData().splashDamage, p.getStaticData().damageType)
		if hitTarget != nil {
			b.dealDamageToActor(p.getStaticData().hitDamage, p.getStaticData().damageType, p.targetActor)
			// add experience
			if !hitTarget.isAlive() {
				if u, ok := hitTarget.(*unit); ok {
					p.whoShot.addExperienceAmount(u.getStaticData().cost)
				}
				if b, ok := hitTarget.(*building); ok {
					p.whoShot.addExperienceAmount(b.getStaticData().cost)
				}
			} else {
				p.whoShot.addExperienceAmount(1)
			}
		}
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
