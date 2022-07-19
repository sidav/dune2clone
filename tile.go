package main

import rl "github.com/gen2brain/raylib-go/raylib"

type tile struct {
	code int
}

func (t *tile) getSpritesAtlas() *spriteAtlas {
	return tilesAtlaces[tableTileStatic[t.code].stringCode]
}

func (t *tile) getSprite() rl.Texture2D {
	return tilesAtlaces[tableTileStatic[t.code].stringCode].atlas[0][0]
}

const (
	TILE_SAND = iota
	TILE_ROCK
	TILE_CONCRETE
)

type tileStaticData struct {
	stringCode string
}

var tableTileStatic = map[int]tileStaticData{
	TILE_SAND: {stringCode: "sand"},
}
