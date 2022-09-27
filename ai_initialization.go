package main

func createAi(f *faction, name, personality string) *aiStruct {
	ai := aiStruct{
		name:                        name,
		controlsFaction:             f,
		current:                     aiAnalytics{},
		alreadyOrderedBuildThisTick: false,
		desired: aiAnalytics{
			defenses:       5,
			builders:       2,
			eco:            2,
			production:     3,
			combatUnits:    20,
			nonCombatUnits: 5,
			harvesters:     5,
			transports:     3,
		},
		max: aiAnalytics{
			nonDefenseBuildings: 30,
			builders:            2,
			eco:                 5,
			production:          5,
			defenses:            10,
			combatUnits:         25,
			nonCombatUnits:      15,
			harvesters:          10,
			transports:          5,
		},
	}

	var persSetter func(*aiStruct)
	if personality == "random" {
		// select random value from a map
		selectedIndex := rnd.Rand(len(aiPersonalitySetters))
		currIndex := 0
		for k, _ := range aiPersonalitySetters {
			if selectedIndex == currIndex {
				persSetter = aiPersonalitySetters[k]
			}
			currIndex++
		}
	} else {
		if _, ok := aiPersonalitySetters[personality]; !ok {
			panic("No such AI personality: " + personality)
		}
		persSetter = aiPersonalitySetters[personality]
	}
	persSetter(&ai)
	return &ai
}

var aiPersonalitySetters = map[string]func(*aiStruct) {
	"balanced": func(ai *aiStruct) {
		ai.personalityName = "Balanced"
		ai.moneyPoorMax = 2500
		ai.moneyRichMin = 10000
		ai.taskForces = []*aiTaskForce{
			{
				mission:     AITF_MISSION_RECON,
				desiredSize: 1,
				units:       make([]*unit, 0),
			},
			{
				mission:                      AITF_MISSION_ATTACK,
				desiredSize:                  10,
				maxFullnessPercentForRetreat: 10,
				units:                        make([]*unit, 0),
			},
			{
				mission:     AITF_MISSION_DEFEND,
				desiredSize: 5,
				units:       make([]*unit, 0),
			},
		}
	},
	"rush": func(ai *aiStruct){
		ai.personalityName = "Rush"
		ai.moneyPoorMax = 1000
		ai.moneyRichMin = 5000
		ai.taskForces = []*aiTaskForce{
			{
				mission:     AITF_MISSION_RECON,
				desiredSize: 1,
				units:       make([]*unit, 0),
			},
			{
				mission:                      AITF_MISSION_ATTACK,
				desiredSize:                  2,
				maxFullnessPercentForRetreat: 0,
				units:                        make([]*unit, 0),
			},
			{
				mission:                      AITF_MISSION_ATTACK,
				desiredSize:                  2,
				maxFullnessPercentForRetreat: 0,
				units:                        make([]*unit, 0),
			},
			{
				mission:     AITF_MISSION_DEFEND,
				desiredSize: 2,
				units:       make([]*unit, 0),
			},
		}
	},
	"turtle": func(ai *aiStruct){
		ai.personalityName = "Turtle"
		ai.moneyPoorMax = 2500
		ai.moneyRichMin = 5000
		ai.taskForces = []*aiTaskForce{
			{
				mission:     AITF_MISSION_RECON,
				desiredSize: 1,
				units:       make([]*unit, 0),
			},
			{
				mission:                      AITF_MISSION_ATTACK,
				desiredSize:                  25,
				maxFullnessPercentForRetreat: 10,
				units:                        make([]*unit, 0),
			},
			{
				mission:     AITF_MISSION_DEFEND,
				desiredSize: 4,
				units:       make([]*unit, 0),
			},
			{
				mission:     AITF_MISSION_DEFEND,
				desiredSize: 4,
				units:       make([]*unit, 0),
			},
			{
				mission:     AITF_MISSION_DEFEND,
				desiredSize: 4,
				units:       make([]*unit, 0),
			},
		}
	},
	"pressure": func(ai *aiStruct){
		ai.personalityName = "Pressure"
		ai.moneyPoorMax = 2500
		ai.moneyRichMin = 8000
		ai.taskForces = []*aiTaskForce{
			{
				mission:     AITF_MISSION_RECON,
				desiredSize: 1,
				units:       make([]*unit, 0),
			},
			{
				mission:                      AITF_MISSION_ATTACK,
				desiredSize:                  3,
				maxFullnessPercentForRetreat: 0,
				units:                        make([]*unit, 0),
			},
			{
				mission:                      AITF_MISSION_ATTACK,
				desiredSize:                  15,
				maxFullnessPercentForRetreat: 0,
				units:                        make([]*unit, 0),
			},
			{
				mission:     AITF_MISSION_DEFEND,
				desiredSize: 4,
				units:       make([]*unit, 0),
			},
			{
				mission:     AITF_MISSION_DEFEND,
				desiredSize: 4,
				units:       make([]*unit, 0),
			},
		}
	},
}
