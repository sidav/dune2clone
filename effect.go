package main

type effect struct {
	centerX, centerY   float64
	code               effectCode
	creationTick       int
	splashCircleRadius float64
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
	EFFECT_REGULAR_EXPLOSION
	EFFECT_BIGGER_EXPLOSION
)

type effectStatic struct {
	spriteCode      string
	defaultLifeTime int
}

var sTableEffects = map[effectCode]*effectStatic{
	EFFECT_SMALL_EXPLOSION: {
		spriteCode:      "smallexplosion",
		defaultLifeTime: 16,
	},
	EFFECT_REGULAR_EXPLOSION: {
		spriteCode:      "regularexplosion",
		defaultLifeTime: 15,
	},
	EFFECT_BIGGER_EXPLOSION: {
		spriteCode:      "biggerexplosion",
		defaultLifeTime: 36,
	},
}
