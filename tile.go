package main

type tile struct {
	code               int
	spriteVariantIndex int
	resourcesAmount    int

	// for faster collision detection
	isOccupiedByActor actor

	hasResourceVein bool
}

func (t *tile) getStaticData() *tileStaticData {
	return sTableTiles[t.code]
}

func (t *tile) canHaveResources() bool {
	return t.getStaticData().canHaveResources && (!t.hasResourceVein)
}

const (
	TILE_SAND = iota
	TILE_BUILDABLE
	TILE_BUILDABLE_DAMAGED
	TILE_ROCK
	TILE_CONCRETE
)

type tileStaticData struct {
	spriteCodes      []string
	canBuildHere     bool
	canHaveResources bool // for resource growth, for example
	canBeWalkedOn    bool
}

var sTableTiles = map[int]*tileStaticData{
	TILE_SAND: {
		spriteCodes:      []string{"sand1", "sand2", "sand3"},
		canBuildHere:     false,
		canHaveResources: true,
		canBeWalkedOn:    true,
	},
	TILE_BUILDABLE: {
		spriteCodes:   []string{"buildable1"},
		canBuildHere:  true,
		canBeWalkedOn: true,
	},
	TILE_BUILDABLE_DAMAGED: {
		spriteCodes:   []string{"buildabledamaged"},
		canBuildHere:  true,
		canBeWalkedOn: true,
	},
	TILE_ROCK: {
		spriteCodes:  []string{"rock1"},
		canBuildHere: false,
	},
}
