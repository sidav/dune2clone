package main

import rl "github.com/gen2brain/raylib-go/raylib"

type building struct {
	currentAction      action
	topLeftX, topLeftY int // tile coords
	code               int
	faction            *faction
	isSelected         bool
}

func (b *building) markSelected(s bool) {
	b.isSelected = s
}

func (b *building) getName() string {
	return b.getStaticData().displayedName
}

func (b *building) getCurrentAction() *action {
	return &b.currentAction
}

func (b *building) getFaction() *faction {
	return b.faction
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
	BLD_POWERPLANT
	BLD_FACTORY
)

type buildingStatic struct {
	w, h          int
	displayedName string
	cost          int
	buildTime     int   // seconds
	builds        []int // buildings
	produces      []int // units

	// ui-only things:
	hotkeyToBuild string
}

var sTableBuildings = map[int]*buildingStatic{
	BLD_BASE: {
		w:             2,
		h:             2,
		displayedName: "Construction Yard",
		cost:          0,
		buildTime:     100,
		builds:        []int{BLD_POWERPLANT, BLD_FACTORY},
		produces:      nil,
	},
	BLD_POWERPLANT: {
		w:             2,
		h:             2,
		displayedName: "Power Plant",
		cost:          500,
		buildTime:     5,
		builds:        nil,
		produces:      nil,
		hotkeyToBuild: "P",
	},
	BLD_FACTORY: {
		w:             3,
		h:             2,
		displayedName: "Factory",
		cost:          1000,
		buildTime:     10,
		builds:        nil,
		produces:      []int{UNT_TANK},
		hotkeyToBuild: "F",
	},
}
