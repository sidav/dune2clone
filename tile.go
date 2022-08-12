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
	TILE_ROCK
	TILE_CONCRETE
)

type tileStaticData struct {
	spriteCodes  []string
	canBuildHere bool
}

var sTableTiles = map[int]*tileStaticData{
	TILE_SAND: {
		spriteCodes:  []string{"sand1", "sand2", "sand3"},
		canBuildHere: false,
	},
	TILE_BUILDABLE: {
		spriteCodes:  []string{"buildable1"},
		canBuildHere: true,
	},
}
