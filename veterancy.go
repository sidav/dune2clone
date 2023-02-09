package main

func getExperienceLevelByAmountAndCost(totalExp, cost int) int {
	multipliedCost := float64(cost) * config.Gameplay.CostForLevelUpMultiplier
	exponent := config.Gameplay.CostForLevelUpExponent
	floatTotalExp := float64(totalExp)
	switch {
	case floatTotalExp < multipliedCost:
		return 0
	case floatTotalExp < 2*multipliedCost*exponent:
		return 1
	case floatTotalExp < 3*multipliedCost*exponent*exponent:
		return 2
	case floatTotalExp < 4*multipliedCost*exponent*exponent*exponent:
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
	return (100 + (level * config.Gameplay.VeterancyDamageBonusForLevelPercent)) * damage / 100
}

func modifyTurretRangeByUnitExpLevel(tRange, level int) int {
	if level < MAX_VETERANCY_LEVEL-1 {
		return tRange
	}
	return tRange + 1
}

func modifyMaxHpByExpLevel(maxHp, level int) int {
	return (100 + level*config.Gameplay.VeterancyHpBonusForLevelPercent) * maxHp / 100
}

func modifyUnitSpeedByLevel(spd float64, level int) float64 {
	return spd * (1.0 + float64(level)*config.Gameplay.VeterancySpeedBonusForLevel)
}

func modifyTurretCooldownByExpLevel(tCd, level int) int {
	return (100 - level*config.Gameplay.VeterancyFireCooldownReductionForLevelPercent) * tCd / 100
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
