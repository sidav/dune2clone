package main

import rl "github.com/gen2brain/raylib-go/raylib"

type building struct {
	topLeftX, topLeftY int // tile coords
	code               int

	isSelected bool
}

func (b *building) markSelected(s bool) {
	b.isSelected = s
}

func (b *building) getSprite() rl.Texture2D {
	return buildingsAtlaces[b.code].atlas[0][0]
}

func (b *building) isPresentAt(tileX, tileY int) bool {
	w, h := b.getStaticData().w, b.getStaticData().h
	return areCoordsInTileRect(tileX, tileY, b.topLeftX, b.topLeftY, w, h)
}

func (b *building) getStaticData() *buildingStatic {
	return sTableBuildings[b.code]
}

//////////////////////////////////////

const (
	BLD_BASE = iota
)

type buildingStatic struct {
	w, h int
}

var sTableBuildings = map[int]*buildingStatic{
	BLD_BASE: {
		w: 2,
		h: 2,
	},
}
