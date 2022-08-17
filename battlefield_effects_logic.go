package main

func (b *battlefield) RandomlyAddEffectInTileRect(code effectCode, startTickSpread, tx, ty, w, h, count int) {
	for i := 0; i < count; i++ {
		b.addEffect(&effect{
			centerX:      float64(rnd.RandInRange(tx*10, (tx+w)*10)) / 10,
			centerY:      float64(rnd.RandInRange(ty*10, (ty+h)*10)) / 10,
			code:         code,
			creationTick: b.currentTick + rnd.Rand(startTickSpread),
		},
		)
	}
}

func (b *battlefield) actForEffect(e *effect) {

}
