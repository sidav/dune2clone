package main

import (
	"dune2clone/geometry"
	"math"
)

type armorCode int

const (
	ARMORTYPE_FORGOTTEN_TO_BE_SET armorCode = iota
	ARMORTYPE_INFANTRY
	ARMORTYPE_HEAVY
	ARMORTYPE_BUILDING
)

type damageCode int

const (
	DAMAGETYPE_FORGOTTEN_TO_BE_SET damageCode = iota
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
		unt.recalculateSquadSize()
	}
}

func (b *battlefield) dealSplashDamage(centerX, centerY, radius float64, damage int, damageType damageCode) {
	// TODO: air units too
	searchRadius := int(math.Ceil(radius))
	ctx, cty := geometry.TrueCoordsToTileCoords(centerX, centerY)
	for x := ctx - searchRadius; x <= ctx+searchRadius; x++ {
		for y := cty - searchRadius; y <= cty+searchRadius; y++ {
			if b.areTileCoordsValid(x, y) && b.tiles[x][y].isOccupiedByActor != nil {

				if u, ok := b.tiles[x][y].isOccupiedByActor.(*unit); ok {
					if geometry.GetApproxDistFloat64(u.centerX, u.centerY, centerX, centerY) <= radius+0.5 {
						b.dealDamageToActor(damage, damageType, u)
					}
				}

				if bld, ok := b.tiles[x][y].isOccupiedByActor.(*building); ok {
					cx, cy := geometry.TileCoordsToTrueCoords(x, y)
					if geometry.GetApproxDistFloat64(cx, cy, centerX, centerY) <= radius+0.25 {
						b.dealDamageToActor(damage, damageType, bld)
					}
				}

			}
		}
	}
}

func calculateDamageOnArmor(dmg int, dmgType damageCode, armType armorCode) int {
	if dmgType == DAMAGETYPE_FORGOTTEN_TO_BE_SET {
		panic("Oh, damage is nothing")
	}
	if armType == ARMORTYPE_FORGOTTEN_TO_BE_SET {
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
		case ARMORTYPE_BUILDING:
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

	if dmg > 0 {
		dmg = int(math.Round(float64(dmg*percent) / 100.0)) // round up when neccessary
		if dmg == 0 {
			dmg = 1
		}
	}
	return dmg
}
