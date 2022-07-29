package main

func (b *battlefield) actorForActorsTurret(a actor) {
	if a.getCurrentAction().code != ACTION_ROTATE {
		if u, ok := a.(*unit); ok {
			b.actTurret(a, u.turret)
		}
	}
}

func (b *battlefield) actTurret(a actor, t *turret) {
	if t.nextTickToAct > b.currentTick {
		return
	}
	tx, ty := 0, 0
	if u, ok := a.(*unit); ok {
		tx, ty = trueCoordsToTileCoords(u.centerX, u.centerY)
	}
	if _, ok := a.(*building); ok {
		panic("Not implemented")
	}
	// if targetActor not set...
	actorsInRange := b.getListOfActorsInRangeFrom(tx, ty, t.getStaticData().fireRange)
	for l := actorsInRange.Front(); l != nil; l = l.Next() {
		targetCandidate := l.Value
		if targetCandidate.(actor).getFaction() == a.getFaction() {
			continue
		}
		if tc, ok := targetCandidate.(*unit); ok {
			targX, targY := trueCoordsToTileCoords(tc.centerX, tc.centerY)
			rotateTo := getDegreeOfIntVector(targX-tx, targY-ty)
			if t.rotationDegree == rotateTo {
				debugWritef("tick %d: PEWPEW\n", b.currentTick) // TODO
				t.nextTickToAct = b.currentTick + t.getStaticData().attackCooldown
			} else if t.canRotate() {
				t.rotationDegree += getDiffForRotationStep(t.rotationDegree, rotateTo, t.getStaticData().rotateSpeed)
			}
		}
	}
}
