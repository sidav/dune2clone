package main

import "dune2clone/geometry"

type aiDecisionWeight struct {
	weightCode string
	weight     int
}

func (ai *aiStruct) selectWhatToBuild(builder *building) buildingCode {
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
		decisionWeights = append(decisionWeights, aiDecisionWeight{"defense", 3})
	}
	// production
	if ai.current.production == 0 {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"production", 10})
	} else if ai.current.production < ai.desired.production {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"production", 5})
	}

	var code buildingCode = -1
	decidedIndex := -1
	for code == -1 {
		decidedIndex = rnd.SelectRandomIndexFromWeighted(len(decisionWeights), func(i int) int { return decisionWeights[i].weight })
		code = ai.selectRandomBuildableCodeByFunction(availableCodes, decisionWeights[decidedIndex].weightCode)
	}
	ai.debugWritef("decided to build %s from weights %v\n", decisionWeights[decidedIndex].weightCode, decisionWeights)
	return code
}

func (ai *aiStruct) selectRandomBuildableCodeByFunction(availableCodes []buildingCode, function string) buildingCode {
	candidates := make([]buildingCode, 0)
	for i := range availableCodes {
		fncts := ai.deduceBuildingFunctions(availableCodes[i])
		for j := range fncts {
			if (function == "any" || fncts[j] == function) && ai.canUseBuilding(availableCodes[i]) {
				candidates = append(candidates, availableCodes[i])
			}
		}
	}
	if len(candidates) == 0 {
		//panic("No such function: " + function)
		return -1
	}

	// assign weight for random selection according to AI current money
	index := -1
	for index == -1 {
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

func (ai *aiStruct) deduceBuildingFunctions(bldCode buildingCode) []string {
	codes := make([]string, 0)
	bsd := sTableBuildings[bldCode]
	if bsd.receivesResources {
		codes = append(codes, "eco")
	}
	if bsd.givesEnergy > 0 {
		codes = append(codes, "energy")
	}
	if bsd.storageAmount > 0 {
		codes = append(codes, "silo")
	}
	if len(bsd.builds) > 0 {
		codes = append(codes, "builder")
	}
	if len(bsd.produces) > 0 {
		codes = append(codes, "production")
	}
	if bsd.turretCode != TRT_NONE {
		codes = append(codes, "defense")
	}
	return codes
}

func (ai *aiStruct) canUseBuilding(bldCode buildingCode) bool {
	switch bldCode {
	default:
		return ai.controlsFaction.isTechAvailableForBuildingOfCode(bldCode)
	}
}

func (ai *aiStruct) placeBuilding(b *battlefield, builder, whatIsBuilt *building) {
	distToSearchFromBuilder := 3 + (2 * ai.current.nonDefenseBuildings / 3)
	startX, startY := geometry.TrueCoordsToTileCoords(builder.getPhysicalCenterCoords())
	sx, sy := -1, -1
	// defenses placement logic
	if whatIsBuilt.turret != nil {
		w, h := b.getSize()
		sx, sy = geometry.SpiralSearchForLowestScoreFrom(
			func(x, y int) int {
				return geometry.GetApproxDistFromTo(x, y, w/2, h/2)/10 -
					geometry.GetApproxDistFromTo(startX, startY, x, y) +
					+ rnd.Rand(distToSearchFromBuilder)
			},
			func(x, y int) bool {
				return geometry.GetSqDistFromCoordsToRectangleBorder(x, y, 0, 0, w, h) > 10 &&
					b.canBuildingBePlacedAt(whatIsBuilt, x, y, 1, false)
			},
			startX, startY, distToSearchFromBuilder, rnd.Rand(4),
		)

	// non-defenses placement logic
	} else {
		sx, sy = geometry.SpiralSearchForClosestConditionFrom(
			func(x, y int) bool {
				return b.canBuildingBePlacedAt(whatIsBuilt, x, y, 1, false)
			},
			startX, startY, distToSearchFromBuilder, rnd.Rand(4))
		if sx == -1 || sy == -1 {
			sx, sy = geometry.SpiralSearchForClosestConditionFrom(
				func(x, y int) bool {
					return b.canBuildingBePlacedAt(whatIsBuilt, x, y, 0, false)
				},
				startX, startY, distToSearchFromBuilder, rnd.Rand(4))
		}
	}
	if sx != -1 && sy != -1 {
		builder.currentOrder.targetTileX = sx
		builder.currentOrder.targetTileY = sy
	}
}
