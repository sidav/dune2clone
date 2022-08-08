package main

type projectile struct {
	faction          *faction
	code             int
	centerX, centerY float64
	rotationDegree   int
	fuel             float64 // how many 'speeds' it spends until it is destroyed
	targetActor      actor   // for homing projectiles
	damage           int     // set by turret, not proj static data
}

func (p *projectile) getStaticData() *projectileStatic {
	return sTableProjectiles[p.code]
}

const (
	PRJ_SHELL = iota
	PRJ_BULLETS
	PRJ_MISSILE
)

type projectileStatic struct {
	spriteCode    string
	size          float64
	speed         float64
	rotationSpeed int
}

var sTableProjectiles = map[int]*projectileStatic{
	PRJ_SHELL: {
		spriteCode: "shell",
		size:       0.3,
		speed:      0.7,
	},
	PRJ_BULLETS: {
		spriteCode: "bullets",
		size:       0.2,
		speed:      0.7,
	},
	PRJ_MISSILE: {
		spriteCode:    "missile",
		size:          0.3,
		speed:         0.2,
		rotationSpeed: 1,
	},
}
