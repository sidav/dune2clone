package main

import "dune2clone/geometry"

type aiDecisionWeight struct {
	weightCode string
	weight     int
}

func (ai *aiStruct) selectWhatToBuild(builder *building) int {
	availableCodes := builder.getStaticData().builds
	// make the list of weights
	decisionWeights := []aiDecisionWeight{{"any", 1}}
	// create weights according to the needs
	// eco
	if ai.current.eco == 0 && ai.isPoor() {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"eco", 200})
	} else if ai.isPoor() && ai.current.eco < ai.desired.eco {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"eco", 15})
	} else if ai.current.eco < ai.desired.eco {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"eco", 5})
	}
	// energy
	if ai.controlsFaction.getAvailableEnergy() <= 0 {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"energy", 100})
	} else if ai.controlsFaction.getAvailableEnergy() <= 5 {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"energy", 5})
	}
	// silos
	if ai.controlsFaction.getStorageRemaining() < 500 && ai.controlsFaction.resourceStorage > 0 {
		if ai.isPoor() {
			decisionWeights = append(decisionWeights, aiDecisionWeight{"silo", 50})
		} else if ai.isRich() {
			decisionWeights = append(decisionWeights, aiDecisionWeight{"silo", 10})
		} else {
			decisionWeights = append(decisionWeights, aiDecisionWeight{"silo", 25})
		}
	}
	// builders
	if ai.current.builders < ai.desired.builders && ai.current.builders < ai.max.builders {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"builder", 1})
	}
	// defenses
	if ai.current.defenses < ai.desired.defenses && ai.current.defenses < ai.max.defenses {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"defense", 4})
	}
	// production
	if ai.current.production == 0 {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"production", 10})
	} else if ai.current.production < ai.desired.production {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"production", 5})
	}

	code := -1
	decidedIndex := -1
	for code == -1 {
		decidedIndex = rnd.SelectRandomIndexFromWeighted(len(decisionWeights), func(i int) int { return decisionWeights[i].weight })
		code = ai.selectRandomBuildableCodeByFunction(availableCodes, decisionWeights[decidedIndex].weightCode)
	}
	debugWritef("AI %s decided to build %s from weights %v\n", ai.name, decisionWeights[decidedIndex].weightCode, decisionWeights)
	return code
}

func (ai *aiStruct) selectRandomBuildableCodeByFunction(availableCodes []int, function string) int {
	candidates := make([]int, 0)
	switch function {
	case "eco":
		for _, code := range availableCodes {
			if sTableBuildings[code].receivesResources {
				candidates = append(candidates, code)
			}
		}
	case "energy":
		for _, code := range availableCodes {
			if sTableBuildings[code].givesEnergy > 0 && sTableBuildings[code].builds == nil {
				candidates = append(candidates, code)
			}
		}
	case "silo":
		for _, code := range availableCodes {
			if sTableBuildings[code].storageAmount > 0 { // && sTableBuildings[code].receivesResources == false {
				candidates = append(candidates, code)
			}
		}
	case "builder":
		for _, code := range availableCodes {
			if sTableBuildings[code].builds != nil {
				candidates = append(candidates, code)
			}
		}
	case "production":
		for _, code := range availableCodes {
			if sTableBuildings[code].produces != nil {
				candidates = append(candidates, code)
			}
		}
	case "defense":
		for _, code := range availableCodes {
			if sTableBuildings[code].turretCode != TRT_NONE {
				candidates = append(candidates, code)
			}
		}
	case "any":
		candidates = availableCodes
	default:
		panic("No such function: " + function)
	}
	if len(candidates) == 0 {
		return -1
	}

	// assign weight for random selection according to AI current money
	index := -1
	for index == -1 || !ai.canUseBuilding(candidates[index]) {
		index = rnd.SelectRandomIndexFromWeighted(len(candidates),
			func(x int) int {
				consideredCode := candidates[x]
				if int(ai.controlsFaction.getMoney()) > sTableBuildings[consideredCode].cost {
					return 5
				} else if !ai.isPoor() {
					return 3
				}
				return 1
			},
		)
	}
	return candidates[index]
}

func (ai *aiStruct) canUseBuilding(bldCode int) bool {
	switch bldCode {
	case BLD_REPAIR_DEPOT: return false
	case BLD_FUSION: return false
	default: return true
	}
}

func (ai *aiStruct) placeBuilding(b *battlefield, builder, whatIsBuilt *building) {
	startX, startY := geometry.TrueCoordsToTileCoords(builder.getPhysicalCenterCoords())
	placementSearchFunc := geometry.SpiralSearchForClosestConditionFrom
	if whatIsBuilt.turret != nil {
		placementSearchFunc = geometry.SpiralSearchForFarthestConditionFrom
	}
	sx, sy := placementSearchFunc(
		func(x, y int) bool {
			return b.canBuildingBePlacedAt(whatIsBuilt, x, y, 1, false)
		},
		startX, startY, 16, 0)
	if sx == -1 || sy == -1 {
		sx, sy = placementSearchFunc(
			func(x, y int) bool {
				return b.canBuildingBePlacedAt(whatIsBuilt, x, y, 0, false)
			},
			startX, startY, 16, 0)
	}
	if sx != -1 && sy != -1 {
		builder.currentOrder.targetTileX = sx
		builder.currentOrder.targetTileY = sy
	}
}

func (ai *aiStruct) deployDeployableUnitSomewhere(b *battlefield, u *unit) {
	bld := createBuilding(u.getStaticData().deploysInto, 0, 0, u.faction)
	tx, ty := u.getTileCoords()
	if !b.canUnitBeDeployedAt(u, tx, ty) {
		depX, depY := geometry.SpiralSearchForFarthestConditionFrom(
			func(x, y int) bool {
				return b.canBuildingBePlacedAt(bld, x, y, 0, true) && rnd.OneChanceFrom(32)
			},
			tx, ty, 16, rnd.Rand(4),
		)
		u.currentOrder.code = ORDER_MOVE
		u.currentOrder.setTargetTileCoords(depX, depY)
	} else {
		u.currentOrder.code = ORDER_DEPLOY
	}
}
