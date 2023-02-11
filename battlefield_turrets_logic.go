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
			if bld.turret.nextTickToAct <= b.currentTick && bld.faction.getAvailableEnergy() >= 0 && !bld.isUnderConstruction() {
				b.actTurret(a, bld.turret)
			}
		}
	}
}

func (b *battlefield) actTurret(shooter actor, t *turret) {
	// shooterTileX, shooterTileY := 0, 0
	shooterX, shooterY := 0.0, 0.0
	turretRange := modifyTurretRangeByUnitExpLevel(t.getStaticData().FireRange, shooter.getExperienceLevel())
	if u, ok := shooter.(*unit); ok {
		// shooterTileX, shooterTileY = geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
		shooterX, shooterY = u.centerX, u.centerY
	}
	if bld, ok := shooter.(*building); ok {
		// shooterTileX, shooterTileY = bld.topLeftX, bld.topLeftY
		shooterX, shooterY = bld.getPhysicalCenterCoords()
	}
	orderedTarget := shooter.getCurrentOrder().targetActor
	if orderedTarget != nil && orderedTarget.getFaction() != shooter.getFaction() {
		t.targetActor = orderedTarget
	}
	if t.targetActor != nil && (!b.canFactionSeeActor(shooter.getFaction(), t.targetActor) ||
		!t.targetActor.isAlive() || !b.areActorsInRangeFromEachOther(shooter, t.targetActor, turretRange) ||
		!b.canTurretAttackActor(t, t.targetActor)) {

		t.targetActor = nil
	}
	if t.targetActor == nil {
		// if targetActor not set...
		actorsInRange := b.getListOfActorsInRangeFromActor(shooter, turretRange)
		for l := actorsInRange.Front(); l != nil; l = l.Next() {
			targetCandidate := l.Value.(actor)
			if targetCandidate.getFaction() != shooter.getFaction() &&
				b.canTurretAttackActor(t, targetCandidate) && b.canFactionSeeActor(shooter.getFaction(), targetCandidate) {
				// try to select a target which pose max danger to the shooter
				if t.targetActor != nil {
					if b.canActorAttackActor(t.targetActor, shooter) {
						break
					} else {
						t.targetActor = targetCandidate
					}
				} else {
					t.targetActor = targetCandidate
				}
			}
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
		t.rotationDegree += geometry.GetDiffForRotationStep(t.rotationDegree, rotateTo, t.getStaticData().RotateSpeed)
		t.normalizeDegrees()
	}

	if abs(t.rotationDegree-rotateTo) <= t.getStaticData().FireSpreadDegrees/2 {
		b.shootAsTurretAtTarget(shooter, t)
	}
}

func (b *battlefield) shootAsTurretAtTarget(shooter actor, t *turret) {
	shooterX, shooterY := shooter.getPhysicalCenterCoords()
	targetCenterX, targetCenterY := t.targetActor.getPhysicalCenterCoords()
	vectX, vectY := geometry.DegreeToUnitVector(t.rotationDegree)
	projX, projY := shooterX+vectX/2, shooterY+vectY/2 // TODO: turret displacement
	degreeSpread := rnd.RandInRange(-t.getStaticData().FireSpreadDegrees, t.getStaticData().FireSpreadDegrees)
	rangeSpread := t.getStaticData().ShotRangeSpread * float64(rnd.RandInRange(-100, 100)) / 100
	proj := &projectile{
		faction:        shooter.getFaction(),
		staticData:     t.getStaticData().FiredProjectileData,
		centerX:        projX,
		centerY:        projY,
		rotationDegree: t.rotationDegree + degreeSpread,
		// next 0.5 is for initial projectile displacement
		fuel:        geometry.GetPreciseDistFloat64(targetCenterX, targetCenterY, shooterX, shooterY) + rangeSpread - 0.5,
		whoShot:     shooter,
		targetActor: t.targetActor,
	}
	if proj.isHoming() {
		proj.fuel *= 1.5
	}
	b.addProjectile(proj)

	if t.getStaticData().MaxShotsInVolley > 1 {
		t.shotsInCurrentVolley++
		if t.shotsInCurrentVolley < t.getStaticData().MaxShotsInVolley {
			t.nextTickToAct = b.currentTick + t.getStaticData().CooldownPerShot
			return
		} else {
			t.shotsInCurrentVolley = 0
		}
	}
	t.nextTickToAct = b.currentTick + modifyTurretCooldownByExpLevel(t.getStaticData().CooldownAfterVolley, shooter.getExperienceLevel())
}

func (b *battlefield) canTurretAttackActor(t *turret, a actor) bool {
	if a.isInAir() {
		return t.getStaticData().AttacksAir
	}
	return t.getStaticData().AttacksLand
}
