package main

import "math"

type armorCode int

const (
	ARMORTYPE_CRASH_IF_THIS_SET armorCode = iota
	ARMORTYPE_INFANTRY
	ARMORTYPE_HEAVY
	ARMORTYPE_BUILDING
)

type damageCode int

const (
	DAMAGETYPE_CRASH_IF_THIS_SET damageCode = iota
	DAMAGETYPE_ANTI_INFANTRY
	DAMAGETYPE_ANTI_BUILDING
	DAMAGETYPE_OMNI
	DAMAGETYPE_HEAVY
)

func (b *battlefield) dealDamageToActor(dmg int, dmgType damageCode, act actor) {
	if bld, ok := act.(*building); ok {
		bld.currentHitpoints -= calculateDamageOnArmor(dmg, dmgType, ARMORTYPE_BUILDING)
	}
	if unt, ok := act.(*unit); ok {
		unt.currentHitpoints -= calculateDamageOnArmor(dmg, dmgType, unt.getStaticData().armorType)
		if unt.getStaticData().maxSquadSize > 1 {
			unt.squadSize = int(
				math.Ceil(float64(unt.getStaticData().maxSquadSize) *
					float64(unt.currentHitpoints) / float64(unt.getStaticData().maxHitpoints)),
			)
		}
	}
}

func calculateDamageOnArmor(dmg int, dmgType damageCode, armType armorCode) int {
	if dmgType == DAMAGETYPE_CRASH_IF_THIS_SET {
		panic("Oh, damage is nothing")
	}
	if armType == ARMORTYPE_CRASH_IF_THIS_SET {
		panic("Oh, armor is nothing")
	}
	percent := 100
	switch dmgType {
	case DAMAGETYPE_OMNI:
		percent = 100
	case DAMAGETYPE_ANTI_INFANTRY:
		switch armType {
		case ARMORTYPE_HEAVY:
			percent = 25
		}
	case DAMAGETYPE_HEAVY:
		switch armType {
		case ARMORTYPE_INFANTRY:
			percent = 25
		}
	case DAMAGETYPE_ANTI_BUILDING:
		if armType != ARMORTYPE_BUILDING {
			percent = 25
		}
	}

	dmg = int(math.Round(float64(dmg*percent) / 100.0)) // round up when neccessary
	if dmg == 0 {
		dmg = 1
	}
	return dmg
}
