package main

func (b *battlefield) dealDamageToActor(dmg int, act actor) {
	if bld, ok := act.(*building); ok {
		bld.currentHitpoints -= dmg
	}
	if unt, ok := act.(*unit); ok {
		unt.currentHitpoints -= dmg
	}
}
