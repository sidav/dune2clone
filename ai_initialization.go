package main

func createAi(f *faction, n string) *aiStruct {
	ai := aiStruct{
		name:            n,
		controlsFaction: f,
		current:         aiAnalytics{},
		moneyPoorMax:    5500,
		moneyRichMin:    15000,
		desired: aiAnalytics{
			defenses:            5,
			builders:            2,
			eco:                 2,
			production:          3,
			combatUnits:         20,
			nonCombatUnits:      5,
			harvesters:          5,
			transports:          3,
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
	return &ai
}
