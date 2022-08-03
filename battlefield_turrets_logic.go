package main

import (
	"dune2clone/geometry"
	"math"
)

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
		shooterTileX, shooterTileY = geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
		shooterX, shooterY = u.centerX, u.centerY
	}
	if bld, ok := a.(*building); ok {
		shooterTileX, shooterTileY = bld.topLeftX, bld.topLeftY
		shooterX, shooterY = bld.getPhysicalCenterCoords()
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
	if t.rotationDegree == rotateTo {
		// debugWritef("tick %d: PEWPEW\n", b.currentTick) // TODO
		projX, projY := geometry.TileCoordsToPhysicalCoords(shooterTileX, shooterTileY)
		degreeSpread := rnd.RandInRange(-t.getStaticData().fireSpreadDegrees, t.getStaticData().fireSpreadDegrees)
		rangeSpread := t.getStaticData().shotRangeSpread * (float64(rnd.RandInRange(-100, 100))/100)
		b.addProjectile(&projectile{
			faction:        a.getFaction(),
			code:           PRJ_CANNON,
			centerX:        projX,
			centerY:        projY,
			rotationDegree: t.rotationDegree + degreeSpread,
			fuel:           math.Sqrt((targetCenterX-shooterX)*(targetCenterX-shooterX)+(targetCenterY-shooterY)*(targetCenterY-shooterY)) + rangeSpread*rangeSpread,
		})
		t.nextTickToAct = b.currentTick + t.getStaticData().attackCooldown
	} else if t.canRotate() {
		t.rotationDegree += geometry.GetDiffForRotationStep(t.rotationDegree, rotateTo, t.getStaticData().rotateSpeed)
		t.normalizeDegrees()
	}
}
