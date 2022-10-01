package main

func getExperienceLevelByAmountAndCost(totalExp, cost int) int {
	switch {
	case totalExp < cost:
		return 0
	case totalExp < 2*cost:
		return 1
	case totalExp < 3*cost:
		return 2
	case totalExp < 4*cost:
		return 3
	default:
		return 4
	}
}

func getExpLevelName(expLevel int) string {
	if expLevel == MAX_VETERANCY_LEVEL {
		return "ELITE"
	}
	switch expLevel {
	case 0:
		return "ROOKIE"
	default:
		return "VETERAN"
	}
}

func modifyDamageByUnitExpLevel(damage, level int) int {
	return (100 + (10 * level)) * damage / 100
}

func modifyTurretRangeByUnitExpLevel(tRange, level int) int {
	if level < MAX_VETERANCY_LEVEL-1 {
		return tRange
	}
	return tRange + 1
}

func modifyUnitMaxHpByExpLevel(maxHp, level int) int {
	return (100 + level*5) * maxHp / 100
}

func modifyUnitSpeedByExpLevel(spd float64, level int) float64 {
	return spd * (1.0 + float64(level)*0.05)
}

func modifyTurretCooldownByUnitExpLevel(tCd, level int) int {
	return (100 - 4*level) * tCd / 100
}

func modifyVisionRangeByUnitExpLevel(radius, level int) int {
	if level <= 1 {
		return radius
	}
	return radius + level - 1
}

func getVeterancyBasedRegen(level int) int {
	if level <= 1 {
		return 0
	}
	return level - 1
}
