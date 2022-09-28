package main

func createAi(f *faction, name, personality string) *aiStruct {
	ai := aiStruct{
		name:                        name,
		controlsFaction:             f,
		current:                     aiAnalytics{},
		alreadyOrderedBuildThisTick: false,
		desired: aiAnalytics{
			defenses:       4,
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
				break
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

var aiPersonalitySetters = map[string]func(*aiStruct){
	"balanced": func(ai *aiStruct) {
		ai.personalityName = "Balanced"
		ai.moneyPoorMax = 2500
		ai.moneyRichMin = 10000
		addTaskForceToAi(ai, AITF_MISSION_RECON, 1, 0)
		addTaskForceToAi(ai, AITF_MISSION_ATTACK, 10, 0)
		addTaskForceToAi(ai, AITF_MISSION_DEFEND, 5, 0)
		ai.desired = aiAnalytics{
			defenses:       4,
			builders:       2,
			eco:            2,
			production:     3,
			combatUnits:    20,
			nonCombatUnits: 5,
			harvesters:     5,
			transports:     3,
		}
	},
	"rush": func(ai *aiStruct) {
		ai.personalityName = "Rush"
		ai.moneyPoorMax = 1000
		ai.moneyRichMin = 5000
		addTaskForceToAi(ai, AITF_MISSION_RECON, 1, 0)
		addTaskForceToAi(ai, AITF_MISSION_ATTACK, 2, 0)
		addTaskForceToAi(ai, AITF_MISSION_ATTACK, 2, 0)
		addTaskForceToAi(ai, AITF_MISSION_ATTACK, 2, 0)
		ai.desired = aiAnalytics{
			defenses:       0,
			builders:       2,
			eco:            1,
			production:     2,
			combatUnits:    20,
			nonCombatUnits: 5,
			harvesters:     5,
			transports:     3,
		}
	},
	"turtle": func(ai *aiStruct) {
		ai.personalityName = "Turtle"
		ai.moneyPoorMax = 2500
		ai.moneyRichMin = 5000
		addTaskForceToAi(ai, AITF_MISSION_RECON, 1, 0)
		addTaskForceToAi(ai, AITF_MISSION_ATTACK, 10, 0)
		addTaskForceToAi(ai, AITF_MISSION_DEFEND, 5, 0)
		addTaskForceToAi(ai, AITF_MISSION_DEFEND, 3, 0)
		ai.desired = aiAnalytics{
			defenses:       7,
			builders:       2,
			eco:            2,
			production:     3,
			combatUnits:    20,
			nonCombatUnits: 5,
			harvesters:     5,
			transports:     3,
		}
	},
	"pressure": func(ai *aiStruct) {
		ai.personalityName = "Pressure"
		ai.moneyPoorMax = 2500
		ai.moneyRichMin = 8000
		addTaskForceToAi(ai, AITF_MISSION_RECON, 1, 0)
		addTaskForceToAi(ai, AITF_MISSION_ATTACK, 3, 0)
		addTaskForceToAi(ai, AITF_MISSION_ATTACK, 10, 0)
		addTaskForceToAi(ai, AITF_MISSION_DEFEND, 5, 0)
		ai.desired = aiAnalytics{
			defenses:       3,
			builders:       2,
			eco:            2,
			production:     3,
			combatUnits:    25,
			nonCombatUnits: 5,
			harvesters:     5,
			transports:     3,
		}
	},
}

func addTaskForceToAi(ai *aiStruct, mission taskForceMission, desiredSize, retreatPercent int) {
	ai.taskForces = append(ai.taskForces, &aiTaskForce{
		mission:                      mission,
		desiredSize:                  desiredSize,
		maxFullnessPercentForRetreat: retreatPercent,
		units:                        make([]*unit, 0),
	})
}
