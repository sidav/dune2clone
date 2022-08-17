package main

type effect struct {
	centerX, centerY float64
	code             effectCode
	creationTick     int
}

func (e *effect) getStaticData() *effectStatic {
	return sTableEffects[e.code]
}

func (e *effect) getExpirationPercent(currentTick int) int {
	return 100 * (currentTick - e.creationTick) / e.getStaticData().defaultLifeTime
}

type effectCode int

const (
	EFFECT_NONE = iota
	EFFECT_SMALL_EXPLOSION
)

type effectStatic struct {
	spriteCode      string
	defaultLifeTime int
}

var sTableEffects = map[effectCode]*effectStatic{
	EFFECT_SMALL_EXPLOSION: {
		spriteCode:      "smallexplosion",
		defaultLifeTime: 15,
	},
}
