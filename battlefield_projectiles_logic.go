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
	spd := math.Min(p.getStaticData().Speed, p.fuel)
	p.centerX += spd * vx
	p.centerY += spd * vy
	p.fuel -= spd

	var hitTarget actor
	if p.targetActor != nil && p.isHoming() {
		targX, targY := p.targetActor.getPhysicalCenterCoords()
		rotateTo := geometry.GetDegreeOfFloatVector(targX-p.centerX, targY-p.centerY)
		p.rotationDegree += geometry.GetDiffForRotationStep(p.rotationDegree, rotateTo, p.getStaticData().RotationSpeed)
		p.rotationDegree = geometry.NormalizeDegree(p.rotationDegree)
		if geometry.GetApproxDistFloat64(targX, targY, p.centerX, p.centerY) < 1 {
			hitTarget = p.targetActor
			p.setToRemove = true
		}
	}
	if p.fuel <= 0 && hitTarget == nil {
		tilex, tiley := geometry.TrueCoordsToTileCoords(p.centerX, p.centerY)
		hitTarget = b.getActorAtTileCoordinates(tilex, tiley)
		if hitTarget != nil && !(hitTarget.isAlive() && hitTarget.isInAir() == p.targetActor.isInAir()) {
			hitTarget = nil
		}
		p.setToRemove = true
	}
	if p.setToRemove {
		b.dealSplashDamage(p.centerX, p.centerY, p.getStaticData().SplashRadius,
			modifyDamageByUnitExpLevel(p.getStaticData().SplashDamage, p.whoShot.getExperienceLevel()),
			p.getStaticData().DamageType)
		if hitTarget != nil {
			b.dealDamageToActor(modifyDamageByUnitExpLevel(p.getStaticData().HitDamage, p.whoShot.getExperienceLevel()),
				p.getStaticData().DamageType, p.targetActor)
			// add experience
			expAmount := 1
			if !hitTarget.isAlive() {
				p.whoShot.getFaction().gameStatistics.totalDestroyed++
				if u, ok := hitTarget.(*unit); ok {
					expAmount = u.getStaticData().Cost
				}
				if b, ok := hitTarget.(*building); ok {
					expAmount = b.getStaticData().Cost
				}
			}
			p.whoShot.receiveExperienceAmount(int(p.whoShot.getFaction().experienceMultiplier * float64(expAmount)))
		}
		if p.getStaticData().CreatesEffectOnImpact {
			b.addEffect(&effect{
				centerX:            p.centerX,
				centerY:            p.centerY,
				splashCircleRadius: p.getStaticData().SplashRadius,
				code:               p.getStaticData().EffectCreatedOnImpactCode,
				creationTick:       b.currentTick,
			})
		}
	}
	// debugWritef("%+v spd: %f\n", p, spd)
}
