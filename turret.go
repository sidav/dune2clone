package main

import "dune2clone/geometry"

type turret struct {
	code           int
	rotationDegree int
	nextTickToAct  int

	targetActor              actor
	targetTileX, targetTileY int
}

func (t *turret) canRotate() bool {
	return t.getStaticData().rotateSpeed > 0
}

func (t *turret) getStaticData() *turretStatic {
	return sTableTurrets[t.code]
}

func (t *turret) normalizeDegrees() {
	t.rotationDegree = geometry.NormalizeDegree(t.rotationDegree)
}
