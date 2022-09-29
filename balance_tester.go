package main

func testCombatBalance() {
	debugWritef("===== BALANCE CHECK =====\n")
	for k1, v1 := range sTableUnits {
		if len(v1.turretsData) > 0 {
			for k2, v2 := range sTableUnits {
				if len(v2.turretsData) > 0 && k1 != k2 {
					autobattleTwoUnits(createUnit(k1, 0, 0, nil), createUnit(k2, 0, 0, nil))
				}
			}
		}
	}
}

func autobattleTwoUnits(u1, u2 *unit) {
	round := 0
	damagePossible := false
	for {
		dmg1 := getAutobattledDamageOfAllTurretsOnUnit(u1, u2, round)
		if dmg1 > 0 {
			u2.currentHitpoints -= dmg1
			u2.recalculateSquadSize()
			damagePossible = true
		}
		dmg2 := getAutobattledDamageOfAllTurretsOnUnit(u2, u1, round)
		if dmg2 > 0 {
			damagePossible = true
			u1.currentHitpoints -= dmg2
			u1.recalculateSquadSize()
		}
		if !damagePossible && round >= 250 {
			return
		}
		round++
		if !u1.isAlive() && !u2.isAlive() {
			debugWritef("%-10s DRAWS with %20s in %4d rounds\n", u1.getName(), u2.getName(), round)
			return
		}
		if !u1.isAlive() {
			debugWritef("%-10s LOSES to %22s in %4d rounds, enemy has %d%% HP\n", u1.getName(), u2.getName(), round, u2.getHitpointsPercentage())
			return
		}
		if !u2.isAlive() {
			debugWritef("%-10s WINS on %23s in %4d rounds, left with %d%% HP\n", u1.getName(), u2.getName(), round,
				u1.getHitpointsPercentage())
			return
		}
	}
}

func getAutobattledDamageOfAllTurretsOnUnit(shooter, target *unit, round int) int {
	dmg := 0
	sqSize := 1;
	if shooter.squadSize > 1 {
		sqSize = shooter.squadSize
	}
	for i := 0; i < sqSize; i++ {
		for _, t := range shooter.getStaticData().turretsData {
			if t.attacksLand && !target.isInAir() || t.attacksAir && target.isInAir() {
				if round%t.attackCooldown == 0 {
					dmg += calculateDamageOnArmor(t.firedProjectileData.hitDamage,
						t.firedProjectileData.damageType, target.getStaticData().armorType)
					dmg += calculateDamageOnArmor(t.firedProjectileData.splashDamage,
						t.firedProjectileData.damageType, target.getStaticData().armorType)
				}
			}
		}
	}
	return dmg
}
