package main

import "sort"

type balanceTester struct {
}

func (bt *balanceTester) testCombatBalance() {
	debugWritef("===== BALANCE CHECK =====\n")
	unitCodesToCompare := make([]int, 0)
	for k1, v1 := range sTableUnits {
		if len(v1.turretsData) > 0 {
			unitCodesToCompare = append(unitCodesToCompare, k1)
		}
	}
	bldCodesToCompare := make([]buildingCode, 0)
	for k1, v1 := range sTableBuildings {
		if v1.turretData != nil {
			bldCodesToCompare = append(bldCodesToCompare, k1)
		}
	}
	sort.Slice(unitCodesToCompare, func(i, j int) bool {
		return sTableUnits[unitCodesToCompare[i]].displayedName[0] < sTableUnits[unitCodesToCompare[j]].displayedName[0]
	})
	for i := range unitCodesToCompare {
		debugWritef("\n")
		totalWon := 0.0
		totalBattles := 0.0
		for j := range unitCodesToCompare {
			if i != j {
				won := bt.autobattleTwoActors(
					createUnit(unitCodesToCompare[i], 0, 0, nil),
					createUnit(unitCodesToCompare[j], 0, 0, nil))
				if won > 0 {
					totalWon++
				}
				if won != 0 {
					totalBattles++
				}
			}
		}
		for j := range bldCodesToCompare {
			won := bt.autobattleTwoActors(
				createUnit(unitCodesToCompare[i], 0, 0, nil),
				createBuilding(bldCodesToCompare[j], 0, 0, nil))
			if won > 0 {
				totalWon++
			}
			if won != 0 {
				totalBattles++
			}
		}
		debugWritef("%s win rate is %.1f%%\n", sTableUnits[unitCodesToCompare[i]].displayedName,
			100*totalWon/totalBattles)
	}
}

// returns -1 if lost, 1 if won and 0 if incomparable
func (bt *balanceTester) autobattleTwoActors(u1, u2 actor) int {
	round := 0
	damagePossible := false
	for {
		dmg1 := bt.getAutobattledDamageOfAllTurretsOnUnit(u1, u2, round)
		dmg2 := bt.getAutobattledDamageOfAllTurretsOnUnit(u2, u1, round)
		if dmg1 > 0 {
			bt.dealDamageToActor(u2, dmg1)
			damagePossible = true
			debugWritef("%s hits %s for %d damage!\n", u1.getName(), u2.getName(), dmg1)
		}
		if dmg2 > 0 {
			bt.dealDamageToActor(u1, dmg2)
			damagePossible = true
			debugWritef("%s hits %s for %d damage!\n", u2.getName(), u1.getName(), dmg2)
		}
		if !damagePossible && round >= 500 {
			return 0
		}
		round++
		if !u1.isAlive() && !u2.isAlive() {
			debugWritef("%-10s DRAWS with %20s in %4d rounds\n", u1.getName(), u2.getName(), round)
			return -1
		}
		if !u1.isAlive() {
			debugWritef("%-10s LOSES to %22s in %4d rounds, enemy has %d%% HP\n", u1.getName(), u2.getName(), round, u2.getHitpointsPercentage())
			return -1
		}
		if !u2.isAlive() {
			debugWritef("%-10s WINS on %23s in %4d rounds, left with %d%% HP\n", u1.getName(), u2.getName(), round,
				u1.getHitpointsPercentage())
			return 1
		}
	}
}

func (bt *balanceTester) getAutobattledDamageOfAllTurretsOnUnit(shooter, target actor, round int) int {
	dmg := 0
	sqSize := 1
	// first, calculate damage
	if u, ok := shooter.(*unit); ok {
		if u.squadSize > 1 {
			sqSize = u.squadSize
		}
		for i := 0; i < sqSize; i++ {
			for _, t := range u.getStaticData().turretsData {
				dmg += bt.getTotalDamageOfSingleTurretOnActorInRound(t, target, round)
			}
		}
	}
	if bld, ok := shooter.(*building); ok {
		dmg += bt.getTotalDamageOfSingleTurretOnActorInRound(bld.turret.getStaticData(), target, round)
	}
	return dmg
}

func (bt *balanceTester) dealDamageToActor(target actor, damage int) {
	if bld, ok := target.(*building); ok {
		bld.currentHitpoints -= damage
	}
	if unt, ok := target.(*unit); ok {
		unt.currentHitpoints -= damage
		unt.recalculateSquadSize()
	}
}

func (bt *balanceTester) getTotalDamageOfSingleTurretOnActorInRound(t *turretStatic, target actor, round int) int {
	dmg := 0
	if t.attacksLand && !target.isInAir() || t.attacksAir && target.isInAir() {
		armorType := ARMORTYPE_BUILDING
		if u, ok := target.(*unit); ok {
			armorType = u.getStaticData().armorType
		}
		if round%t.attackCooldown == 0 {
			dmg += calculateDamageOnArmor(t.firedProjectileData.hitDamage,
				t.firedProjectileData.damageType, armorType)
			dmg += calculateDamageOnArmor(t.firedProjectileData.splashDamage,
				t.firedProjectileData.damageType, armorType)
		}
	}
	return dmg
}
