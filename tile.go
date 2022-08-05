package main

type tile struct {
	code               int
	spriteVariantIndex int
	resourcesAmount    int
}

const (
	TILE_SAND = iota
	TILE_ROCK
	TILE_CONCRETE
)

type tileStaticData struct {
	spriteCodes []string
}

var sTableTiles = map[int]tileStaticData{
	TILE_SAND: {spriteCodes: []string{"sand1", "sand2", "sand3"}},
}
