package main

import "fmt"

// check units or buildings datas for possible human errors
func performAllDataSanityChecks() {
	debugWritef("===== %-30s =====\n", "RUNNING STATICS SANITY CHECKS")
	for _, v := range sTableUnits {
		if len(v.turretsData) > 0 {
			for i, t := range v.turretsData {
				performTurretSanity(fmt.Sprintf("%s turret %d", v.displayedName, i+1), t)
			}
		}
		if v.maxHitpoints == 0 {
			debugWritef("! It looks like %s has 0 HP!\n", v.displayedName)
		}
		if v.armorType == ARMORTYPE_FORGOTTEN_TO_BE_SET {
			debugWritef("! It looks like %s has no armor type!\n", v.displayedName)
		}
	}
	for _, v := range sTableBuildings {
		if v.turretData != nil {
			performTurretSanity(v.displayedName, v.turretData)
		}
		if v.maxHitpoints == 0 {
			debugWritef(" It looks like %s has 0 HP!\n", v.displayedName)
		}
	}
	debugWritef("===== %-30s =====\n","STATICS SANITY CHECKS ENDED")
}

func performTurretSanity(source string, t *turretStatic) {
	if t.attacksAir == false && t.attacksLand == false {
		debugWritef("! It looks like %s attacks nothing!\n", source)
	}
	if t.firedProjectileData == nil {
		debugWritef(" It looks like %s has no projectile set!\n", source)
	} else {
		performProjectileSanity(source + " projectile", t.firedProjectileData)
	}
}

func performProjectileSanity(source string, proj *projectileStatic) {
	if proj.hitDamage == 0 && proj.splashDamage == 0 {
		debugWritef(" It looks like %s has zero damage.\n", source)
	}
	if proj.hitDamage > 0 && proj.splashDamage > 0 {
		debugWritef(" It looks like %s has both splash and hit damage.\n", source)
	}
	if proj.damageType == DAMAGETYPE_FORGOTTEN_TO_BE_SET {
		debugWritef(" It looks like %s turret has no damage type set.\n", source)
	}
	if proj.splashRadius > 0 && proj.splashDamage == 0 {
		debugWritef(" It looks like %s has zero-damage splash.\n", source)
	}
	if proj.splashRadius == 0 && proj.splashDamage > 0 {
		debugWritef(" It looks like %s has plash without radius.\n", source)
	}
}
