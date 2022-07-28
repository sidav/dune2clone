package main

type turret struct {
	code           int
	rotationDegree int
	nextTickToAct  int
}

func (t *turret) canRotate() bool {
	return t.getStaticData().rotateSpeed > 0
}

func (t *turret) getStaticData() *turretStatic {
	return sTableTurrets[t.code]
}

func (t *turret) normalizeDegrees() {
	if t.rotationDegree < 0 {
		t.rotationDegree += 360
	}
	if t.rotationDegree >= 360 {
		t.rotationDegree -= 360
	}
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
}
