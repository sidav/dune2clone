package main

import (
	"dune2clone/geometry"
)

func (b *battlefield) actorForActorsTurret(a actor) {
	if a.getCurrentAction().code != ACTION_ROTATE {
		if u, ok := a.(*unit); ok {
			for i := range u.turrets {
				if u.turrets[i].nextTickToAct > b.currentTick {
					continue
				}
				shots := u.squadSize
				if u.squadSize == 0 {
					u.squadSize = 1
				}
				for j := 0; j < shots; j++ {
					b.actTurret(a, u.turrets[i])
				}
			}
		}
		if bld, ok := a.(*building); ok {
			// buildings' turrets won't shoot without energy
			if bld.turret.nextTickToAct <= b.currentTick && bld.faction.getAvailableEnergy() >= 0 {
				b.actTurret(a, bld.turret)
			}
		}
	}
}

func (b *battlefield) actTurret(shooter actor, t *turret) {
	shooterTileX, shooterTileY := 0, 0
	shooterX, shooterY := 0.0, 0.0
	if u, ok := shooter.(*unit); ok {
		shooterTileX, shooterTileY = geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
		shooterX, shooterY = u.centerX, u.centerY
	}
	if bld, ok := shooter.(*building); ok {
		shooterTileX, shooterTileY = bld.topLeftX, bld.topLeftY
		shooterX, shooterY = bld.getPhysicalCenterCoords()
	}
	t.targetActor = nil
	// if targetActor not set...
	actorsInRange := b.getListOfActorsInRangeFrom(shooterTileX, shooterTileY, t.getStaticData().fireRange)
	for l := actorsInRange.Front(); l != nil; l = l.Next() {
		targetCandidate := l.Value.(actor)
		if targetCandidate.getFaction() != shooter.getFaction() && b.canTurretAttackActor(t, targetCandidate) {
			t.targetActor = targetCandidate
			break
		}
	}
	if t.targetActor == nil {
		return
	}
	var targetCenterX, targetCenterY float64
	rotateTo := 0
	if tc, ok := t.targetActor.(*unit); ok {
		targetCenterX, targetCenterY = tc.centerX, tc.centerY
		rotateTo = geometry.GetDegreeOfFloatVector(targetCenterX-shooterX, targetCenterY-shooterY)
	} else if tc, ok := t.targetActor.(*building); ok {
		targetCenterX, targetCenterY = tc.getPhysicalCenterCoords()
		rotateTo = geometry.GetDegreeOfFloatVector(targetCenterX-shooterX, targetCenterY-shooterY)
	}
	if t.canRotate() {
		t.rotationDegree += geometry.GetDiffForRotationStep(t.rotationDegree, rotateTo, t.getStaticData().rotateSpeed)
		t.normalizeDegrees()
	}

	if abs(t.rotationDegree-rotateTo) <= t.getStaticData().fireSpreadDegrees/2 {
		// debugWritef("tick %d: PEWPEW\n", b.currentTick) // TODO
		projX, projY := shooterX, shooterY
		degreeSpread := rnd.RandInRange(-t.getStaticData().fireSpreadDegrees, t.getStaticData().fireSpreadDegrees)
		rangeSpread := t.getStaticData().shotRangeSpread * (float64(rnd.RandInRange(-100, 100)) / 100)
		b.addProjectile(&projectile{
			faction:        shooter.getFaction(),
			code:           t.getStaticData().firesProjectileOfCode,
			centerX:        projX,
			centerY:        projY,
			rotationDegree: t.rotationDegree + degreeSpread,
			fuel:           geometry.GetApproxDistFloat64(targetCenterX, targetCenterY, shooterX, shooterY) + rangeSpread,
			targetActor:    t.targetActor,
			damage:         t.getStaticData().projectileDamage,
		})
		t.nextTickToAct = b.currentTick + t.getStaticData().attackCooldown
	}
}

func (b *battlefield) canTurretAttackActor(t *turret, a actor) bool {
	if a.isInAir() {
		return t.getStaticData().attacksAir
	}
	return t.getStaticData().attacksLand
}
