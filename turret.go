package main

import "dune2clone/geometry"

type turret struct {
	staticData     *TurretStatic
	rotationDegree int
	nextTickToAct  int

	targetActor              actor
	targetTileX, targetTileY int
}

func (t *turret) canRotate() bool {
	return t.getStaticData().RotateSpeed > 0
}

func (t *turret) getStaticData() *TurretStatic {
	return t.staticData
}

func (t *turret) normalizeDegrees() {
	t.rotationDegree = geometry.NormalizeDegree(t.rotationDegree)
}
