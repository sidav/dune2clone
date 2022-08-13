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
	if ai.current.eco == 0 {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"eco", 200})
	} else if ai.current.eco < ai.desired.eco {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"eco", 3})
	}
	// energy
	if ai.controlsFaction.getAvailableEnergy() <= 0 {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"energy", 100})
	} else if ai.controlsFaction.getAvailableEnergy() <= 5 {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"energy", 100})
	}
	// silos
	if ai.controlsFaction.getStorageRemaining() < 500 {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"silo", 10})
	}
	// builders
	if ai.current.builders < ai.desired.builders {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"builder", 1})
	}
	// production
	if ai.current.production == 0 {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"production", 10})
	} else if ai.current.production < ai.desired.production {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"production", 3})
	}

	decidedIndex := rnd.SelectRandomIndexFromWeighted(len(decisionWeights), func(i int) int { return decisionWeights[i].weight })
	return ai.selectRandomBuildableCodeByFunction(availableCodes, decisionWeights[decidedIndex].weightCode)
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
			if sTableBuildings[code].givesEnergy > 0 {
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
	case "any":
		return availableCodes[rnd.Rand(len(availableCodes))]
	default:
		panic("No such function: " + function)
	}
	return candidates[rnd.Rand(len(candidates))]
}

func (ai *aiStruct) placeBuilding(b *battlefield, builder, whatIsBuilt *building) {
	startX, startY := geometry.TrueCoordsToTileCoords(builder.getPhysicalCenterCoords())
	sx, sy := geometry.SpiralSearchForConditionFrom(
		func(x, y int) bool {
			return b.canBuildingBePlacedAt(whatIsBuilt, x, y, 1, false)
		},
		startX, startY, 16, 0)
	if sx == -1 || sy == -1 {
		sx, sy = geometry.SpiralSearchForConditionFrom(
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
