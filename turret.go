package main

import "dune2clone/geometry"

type turret struct {
	staticData     *turretStatic
	rotationDegree int
	nextTickToAct  int

	targetActor              actor
	targetTileX, targetTileY int
}

func (t *turret) canRotate() bool {
	return t.getStaticData().rotateSpeed > 0
}

func (t *turret) getStaticData() *turretStatic {
	return t.staticData
}

func (t *turret) normalizeDegrees() {
	t.rotationDegree = geometry.NormalizeDegree(t.rotationDegree)
}
