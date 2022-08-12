package main

func createAi(f *faction) *aiStruct {
	ai := aiStruct{
		controlsFaction: f,
		current:         aiAnalytics{},
		desired: aiAnalytics{
			builders:   2,
			eco:        1,
			production: 1,
			defenses:   5,
			units:      10,
		},
		max: aiAnalytics{
			buildings:  20,
			builders:   2,
			eco:        5,
			production: 5,
			defenses:   10,
			units:      20,
		},
	}
	return &ai
}
