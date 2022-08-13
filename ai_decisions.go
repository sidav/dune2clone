package main

import "dune2clone/geometry"

func (ai *aiStruct) selectWhatToBuild(builder *building) int {
	availableCodes := builder.getStaticData().builds
	// make the list of weights
	decisionWeights := []struct {
		key    string
		weight int
	}{
		{"eco", 0},
		{"energy", 0},
		{"silo", 0},
		{"builder", 0},
		{"production", 0},
		{"any", 0},
	}
	// create weights according to the needs
	for i := range decisionWeights {
		switch decisionWeights[i].key {
		case "eco":
			if ai.current.eco == 0 {
				decisionWeights[i].weight = 200
			} else if ai.current.eco < ai.desired.eco {
				decisionWeights[i].weight = 3
			}
		case "energy":
			if ai.controlsFaction.getAvailableEnergy() <= 0 {
				decisionWeights[i].weight = 100
			} else if ai.controlsFaction.getAvailableEnergy() <= 5 {
				decisionWeights[i].weight = 3
			}
		case "silo":
			if ai.controlsFaction.getStorageRemaining() < 500 {
				decisionWeights[i].weight = 10
			}
		case "builder":
			if ai.current.builders < ai.desired.builders {
				decisionWeights[i].weight = 1
			} else {
				decisionWeights[i].weight = 0
			}
		case "production":
			if ai.current.production == 0 {
				decisionWeights[i].weight = 10
			}
			if ai.current.production < ai.desired.production {
				decisionWeights[i].weight = 3
			} else {
				decisionWeights[i].weight = 1
			}
		case "any":
			decisionWeights[i].weight = 1
		default:
			panic("No such function: " + decisionWeights[i].key)
		}
	}

	decidedIndex := rnd.SelectRandomIndexFromWeighted(len(decisionWeights), func(i int) int { return decisionWeights[i].weight })
	return ai.selectRandomBuildableCodeByFunction(availableCodes, decisionWeights[decidedIndex].key)
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
			if sTableBuildings[code].givesEnergy > 0 {
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
