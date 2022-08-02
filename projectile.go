package main

type projectile struct {
	faction          *faction
	code             int
	centerX, centerY float64
	rotationDegree   int
	fuel             float64 // how many 'speeds' it spends until it is destroyed
}

func (p *projectile) getStaticData() *projectileStatic {
	return sTableProjectiles[p.code]
}

const (
	PRJ_CANNON = iota
)

type projectileStatic struct {
	spriteCode string
	size       float64
	speed      float64
}

var sTableProjectiles = map[int]*projectileStatic{
	PRJ_CANNON: {
		spriteCode: "cannon",
		size:       0.3,
		speed:      0.7,
	},
}
