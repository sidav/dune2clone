package main

func createAi(f *faction, n string) *aiStruct {
	ai := aiStruct{
		name:            n,
		controlsFaction: f,
		current:         aiAnalytics{},
		moneyPoorMax:    5500,
		moneyRichMin:    15000,
		desired: aiAnalytics{
			builders:   2,
			eco:        2,
			production: 3,
			defenses:   5,
			units:      25,
		},
		max: aiAnalytics{
			nonDefenseBuildings: 30,
			builders:            2,
			eco:                 5,
			production:          5,
			defenses:            10,
			units:               50,
		},
	}
	return &ai
}
