package main

import rl "github.com/gen2brain/raylib-go/raylib"

type building struct {
	currentAction      action
	topLeftX, topLeftY int // tile coords
	code               int
	faction            *faction
	isSelected         bool
	turret             *turret
}

func createBuilding(code, topLeftX, topLeftY int, fact *faction) *building {
	var turr *turret
	if sTableBuildings[code].turretCode != TRT_NONE {
		turr = &turret{code: sTableBuildings[code].turretCode, rotationDegree: 270}
	}
	return &building{
		code:     code,
		topLeftX: topLeftX,
		topLeftY: topLeftY,
		faction:  fact,
		turret:   turr,
	}
}

func (b *building) getPartsSprites() []rl.Texture2D {
	if b.turret != nil {
		return []rl.Texture2D{
			buildingsAtlaces[b.code].atlas[0][0],
			turretsAtlaces[b.turret.getStaticData().spriteCode].atlas[degreeToRotationFrameNumber(b.turret.rotationDegree, 8)][0],
		}
	}
	return []rl.Texture2D{
		buildingsAtlaces[b.code].atlas[0][0],
	}
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
	BLD_TURRET
)

type buildingStatic struct {
	w, h          int
	displayedName string
	cost          int
	buildTime     int   // seconds
	builds        []int // buildings
	produces      []int // units

	turretCode int

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
		builds:        []int{BLD_POWERPLANT, BLD_FACTORY, BLD_TURRET},
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
		buildTime:     7,
		builds:        nil,
		produces:      []int{UNT_TANK, UNT_QUAD},
		hotkeyToBuild: "F",
	},
	BLD_TURRET: {
		w:             1,
		h:             1,
		displayedName: "Defense tower",
		cost:          250,
		buildTime:     1,
		builds:        nil,
		produces:      nil,
		turretCode:    TRT_TANK,
		hotkeyToBuild: "T",
	},
}
