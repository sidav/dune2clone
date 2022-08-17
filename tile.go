package main

type tile struct {
	code               int
	spriteVariantIndex int
	resourcesAmount    int
}

func (t *tile) getStaticData() *tileStaticData {
	return sTableTiles[t.code]
}

const (
	TILE_SAND = iota
	TILE_BUILDABLE
	TILE_BUILDABLE_DAMAGED
	TILE_ROCK
	TILE_CONCRETE
	TILE_RESOURCE_VEIN
)

type tileStaticData struct {
	spriteCodes      []string
	canBuildHere     bool
	canHaveResources bool // for resource growth, for example
	growsResources   bool
}

var sTableTiles = map[int]*tileStaticData{
	TILE_SAND: {
		spriteCodes:      []string{"sand1", "sand2", "sand3"},
		canBuildHere:     false,
		canHaveResources: true,
	},
	TILE_BUILDABLE: {
		spriteCodes:  []string{"buildable1"},
		canBuildHere: true,
	},
	TILE_BUILDABLE_DAMAGED: {
		spriteCodes:  []string{"buildabledamaged"},
		canBuildHere: true,
	},
	TILE_RESOURCE_VEIN: {
		spriteCodes:    []string{"resourcevein"},
		canBuildHere:   false,
		growsResources: true,
	},
}
