package main

import rl "github.com/gen2brain/raylib-go/raylib"

type building struct {
	topLeftX, topLeftY float64
	code               int
}

func (b *building) getSprite() rl.Texture2D {
	return buildingsAtlaces[b.code].atlas[0][0]
}

const (
	BLD_BASE = iota
)

type buildingStatic struct {
	w, h int
}

var sTableBuildings = map[int]buildingStatic {
	BLD_BASE: {
		w: 2,
		h: 2,
	},
}
