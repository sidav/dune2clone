package main

import "math"

func (b *battlefield) actorForActorsTurret(a actor) {
	if a.getCurrentAction().code != ACTION_ROTATE {
		if u, ok := a.(*unit); ok {
			b.actTurret(a, u.turret)
		}
		if bld, ok := a.(*building); ok {
			b.actTurret(a, bld.turret)
		}
	}
}

func (b *battlefield) actTurret(a actor, t *turret) {
	if t.nextTickToAct > b.currentTick {
		return
	}
	shooterTileX, shooterTileY := 0, 0
	shooterX, shooterY := 0.0, 0.0
	if u, ok := a.(*unit); ok {
		shooterTileX, shooterTileY = trueCoordsToTileCoords(u.centerX, u.centerY)
		shooterX, shooterY = u.centerX, u.centerY
	}
	if bld, ok := a.(*building); ok {
		shooterTileX, shooterTileY = bld.topLeftX, bld.topLeftY
		shooterX, shooterY = float64(bld.topLeftX) + float64(bld.getStaticData().w)/2, float64(bld.topLeftY) + float64(bld.getStaticData().h)/2
	}
	t.targetActor = nil
	// if targetActor not set...
	actorsInRange := b.getListOfActorsInRangeFrom(shooterTileX, shooterTileY, t.getStaticData().fireRange)
	for l := actorsInRange.Front(); l != nil; l = l.Next() {
		targetCandidate := l.Value.(actor)
		if targetCandidate.getFaction() != a.getFaction() {
			t.targetActor = targetCandidate
			break
		}
	}
	if tc, ok := t.targetActor.(*unit); ok {
		targX, targY := trueCoordsToTileCoords(tc.centerX, tc.centerY)
		rotateTo := getDegreeOfIntVector(targX-shooterTileX, targY-shooterTileY)
		if a.getFaction() == b.factions[0] {
			debugWritef("TARGET ACQUIRED, rot from %d to %d\n", t.rotationDegree, rotateTo)
		}
		if t.rotationDegree == rotateTo {
			// debugWritef("tick %d: PEWPEW\n", b.currentTick) // TODO
			projX, projY := tileCoordsToPhysicalCoords(shooterTileX, shooterTileY)
			degreeSpread := rnd.RandInRange(-t.getStaticData().fireSpreadDegrees, t.getStaticData().fireSpreadDegrees)
			rangeSpread := 0.0 // t.getStaticData().shotRangeSpread * 100 / float64(rnd.RandInRange(-100, 100))
			b.addProjectile(&projectile{
				faction:        a.getFaction(),
				code:           PRJ_CANNON,
				centerX:        projX,
				centerY:        projY,
				rotationDegree: t.rotationDegree + degreeSpread,
				fuel:           math.Sqrt((tc.centerX-shooterX)*(tc.centerX-shooterX)+(tc.centerY-shooterY)*(tc.centerY-shooterY)) + rangeSpread,
			})
			t.nextTickToAct = b.currentTick + t.getStaticData().attackCooldown
		} else if t.canRotate() {
			t.rotationDegree += getDiffForRotationStep(t.rotationDegree, rotateTo, t.getStaticData().rotateSpeed)
			t.normalizeDegrees()
		}
	}
}
