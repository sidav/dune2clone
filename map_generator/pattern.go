package map_generator

type generationPattern struct {
	Name                string
	MinWidth, MinHeight int
	generationFunc      func(*GeneratedMap)
}

var allPatterns = []*generationPattern{
	{
		Name:           "2 players duel",
		MinWidth:       64,
		MinHeight:      64,
		generationFunc: generateByTwoPlayersPattern,
	},
	{
		Name:           "3 players FFA",
		MinWidth:       64,
		MinHeight:      64,
		generationFunc: generateByThreePlayersPattern,
	},
	{
		Name:           "4 players FFA",
		MinWidth:       96,
		MinHeight:      96,
		generationFunc: generateByFourPlayersPattern,
	},
}

func GetPatternByIndex(ind int) *generationPattern {
	if ind < 0 {
		ind = -ind
	}
	return allPatterns[ind%len(allPatterns)]
}
