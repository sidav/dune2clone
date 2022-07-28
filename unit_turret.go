package main

type turret struct {
	code           int
	rotationDegree int
	nextTickToAct  int

	targetActor              *actor
	targetTileX, targetTileY int
}

func (t *turret) canRotate() bool {
	return t.getStaticData().rotateSpeed > 0
}

func (t *turret) getStaticData() *turretStatic {
	return sTableTurrets[t.code]
}

func (t *turret) normalizeDegrees() {
	t.rotationDegree = normalizeDegree(t.rotationDegree)
}

const (
	TRT_TANK = iota
	TRT_QUAD
)

type turretStatic struct {
	spriteCode  string
	rotateSpeed int
}

var sTableTurrets = map[int]*turretStatic{
	TRT_TANK: {
		spriteCode:  "tank",
		rotateSpeed: 7,
	},
	TRT_QUAD: {
		spriteCode:  "",
		rotateSpeed: 0,
	},
}
