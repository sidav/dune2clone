package main

import "fmt"

func (ai *aiStruct) selectWhatToProduce(producer *building) int {
	availableCodes := producer.getStaticData().produces
	decisionWeights := []aiDecisionWeight{{"any", 1}}

	if ai.current.harvesters < ai.desired.harvesters {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"harvester", 5})
	} else {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"harvester", 1})
	}

	if ai.current.combatUnits < ai.desired.combatUnits {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"combat", 10})
	} else {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"combat", 2})
	}
	if ai.current.transports < ai.desired.transports {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"transport", 2})
	}

	if ai.current.builders < ai.desired.builders {
		decisionWeights = append(decisionWeights, aiDecisionWeight{"deployable", 1})
	}

	code := -1
	decidedIndex := -1
	for code == -1 {
		decidedIndex = rnd.SelectRandomIndexFromWeighted(len(decisionWeights), func(i int) int { return decisionWeights[i].weight })
		code = ai.selectRandomProducableCodeByFunction(availableCodes, decisionWeights[decidedIndex].weightCode)
	}

	ai.debugWritef("decided to produce %s from weights %v\n", decisionWeights[decidedIndex].weightCode, decisionWeights)
	return code
}

func (ai *aiStruct) selectRandomProducableCodeByFunction(availableCodes []int, function string) int {
	candidates := make([]int, 0)
	if function == "any" {
		candidates = availableCodes
	} else {
		for i := range availableCodes {
			if ai.deduceUnitFunction(availableCodes[i]) == function {
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
				if int(ai.controlsFaction.getMoney()) > sTableUnits[consideredCode].cost {
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

func (ai *aiStruct) deduceUnitFunction(untCode int) string {
	usd := sTableUnits[untCode]
	if usd.maxCargoAmount > 0 {
		return "harvester"
	}
	if usd.canBeDeployed {
		return "deployable"
	}
	if usd.isTransport {
		return "transport"
	}
	if len(usd.turretsData) > 0 {
		return "combat"
	}
	panic(fmt.Sprintf("%d: wat is it?!", untCode))
}
