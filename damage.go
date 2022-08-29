package main

import "math"

func (b *battlefield) dealDamageToActor(dmg int, act actor) {
	if bld, ok := act.(*building); ok {
		bld.currentHitpoints -= dmg
	}
	if unt, ok := act.(*unit); ok {
		unt.currentHitpoints -= dmg
		if unt.getStaticData().maxSquadSize > 1 {
			unt.squadSize = int(
				math.Ceil(float64(unt.getStaticData().maxSquadSize)*
					float64(unt.currentHitpoints)/float64(unt.getStaticData().maxHitpoints)),
				)
		}
	}
}
