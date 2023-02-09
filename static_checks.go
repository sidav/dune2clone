package main

import "fmt"

// check units or buildings datas for possible human errors
func performAllDataSanityChecks() {
	debugWritef("===== %-30s =====\n", "RUNNING STATICS SANITY CHECKS")
	for _, v := range sTableUnits {
		if len(v.TurretsData) > 0 {
			for i, t := range v.TurretsData {
				performTurretSanity(fmt.Sprintf("%s turret %d", v.DisplayedName, i+1), t)
			}
		}
		if v.MaxHitpoints == 0 {
			debugWritef("! It looks like %s has 0 HP!\n", v.DisplayedName)
		}
		if v.ArmorType == ARMORTYPE_FORGOTTEN_TO_BE_SET {
			debugWritef("! It looks like %s has no armor type!\n", v.DisplayedName)
		}
	}
	for _, v := range sTableBuildings {
		if v.TurretData != nil {
			performTurretSanity(v.DisplayedName, v.TurretData)
		}
		if v.MaxHitpoints == 0 {
			debugWritef(" It looks like %s has 0 HP!\n", v.DisplayedName)
		}
	}
	debugWritef("===== %-30s =====\n", "STATICS SANITY CHECKS ENDED")
}

func performTurretSanity(source string, t *TurretStatic) {
	if t.AttacksAir == false && t.AttacksLand == false {
		debugWritef("! It looks like %s attacks nothing!\n", source)
	}
	if t.FiredProjectileData == nil {
		debugWritef(" It looks like %s has no projectile set!\n", source)
	} else {
		performProjectileSanity(source+" projectile", t.FiredProjectileData)
	}
}

func performProjectileSanity(source string, proj *projectileStatic) {
	if proj.HitDamage == 0 && proj.SplashDamage == 0 {
		debugWritef(" It looks like %s has zero damage.\n", source)
	}
	if proj.HitDamage > 0 && proj.SplashDamage > 0 {
		debugWritef(" It looks like %s has both splash and hit damage.\n", source)
	}
	if proj.DamageType == DAMAGETYPE_FORGOTTEN_TO_BE_SET {
		debugWritef(" It looks like %s turret has no damage type set.\n", source)
	}
	if proj.SplashRadius > 0 && proj.SplashDamage == 0 {
		debugWritef(" It looks like %s has zero-damage splash.\n", source)
	}
	if proj.SplashRadius == 0 && proj.SplashDamage > 0 {
		debugWritef(" It looks like %s has plash without radius.\n", source)
	}
}
