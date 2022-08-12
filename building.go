package main

import (
	"dune2clone/geometry"
)

type building struct {
	currentAction      action
	currentOrder       order
	currentHitpoints   int
	topLeftX, topLeftY int // tile coords
	code               int
	faction            *faction
	isSelected         bool
	turret             *turret

	unitPlacedInside *unit
}

func createBuilding(code, topLeftX, topLeftY int, fact *faction) *building {
	var turr *turret
	if sTableBuildings[code].turretCode != TRT_NONE {
		turr = &turret{code: sTableBuildings[code].turretCode, rotationDegree: 270}
	}
	return &building{
		code:             code,
		currentHitpoints: sTableBuildings[code].maxHitpoints,
		topLeftX:         topLeftX,
		topLeftY:         topLeftY,
		faction:          fact,
		turret:           turr,
	}
}

func (b *building) markSelected(s bool) {
	b.isSelected = s
}

func (b *building) getDimensionsForConstructon() (int, int, int, int) {
	h := b.getStaticData().h
	// prevent closing bottom side for producing buildings
	if b.getStaticData().produces != nil || b.getStaticData().receivesResources {
		h++
	}
	return b.topLeftX, b.topLeftY, b.getStaticData().w, h
}

func (b *building) getUnitPlacementCoords() (int, int) {
	return b.topLeftX + b.getStaticData().unitPlacementX, b.topLeftY + b.getStaticData().unitPlacementY
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

func (b *building) getPhysicalCenterCoords() (float64, float64) {
	return float64(b.topLeftX) + float64(b.getStaticData().w)/2, float64(b.topLeftY) + float64(b.getStaticData().h)/2
}

func (b *building) isPresentAt(tileX, tileY int) bool {
	w, h := b.getStaticData().w, b.getStaticData().h
	return geometry.AreCoordsInTileRect(tileX, tileY, b.topLeftX, b.topLeftY, w, h)
}

func (b *building) getStaticData() *buildingStatic {
	return sTableBuildings[b.code]
}

//////////////////////////////////////

type buildingStatic struct {
	w, h          int
	displayedName string
	cost          int
	buildTime     int   // seconds
	builds        []int // buildings
	produces      []int // units
	maxHitpoints  int

	turretCode int

	receivesResources bool // is refinery
	// CanUnitBePlacedHere            bool // Removed for now, as duplicated by receivesResources
	unitPlacementX, unitPlacementY int // tile coords for placed unit draw

	givesFreeUnitOnCreation   bool
	codeForFreeUnitOnCreation int

	givesEnergy, consumesEnergy int
	storageAmount               float64

	// ui-only things:
	hotkeyToBuild string
}

const (
	BLD_BASE = iota
	BLD_POWERPLANT
	BLD_FACTORY
	BLD_TURRET
	BLD_REFINERY
	BLD_SILO
)

var sTableBuildings = map[int]*buildingStatic{
	BLD_BASE: {
		maxHitpoints:  1000,
		w:             2,
		h:             2,
		displayedName: "Construction Yard",
		cost:          2500,
		buildTime:     30,
		builds:        []int{BLD_BASE, BLD_POWERPLANT, BLD_REFINERY, BLD_FACTORY, BLD_TURRET, BLD_SILO},
		givesEnergy:   10,
		hotkeyToBuild: "B",
	},
	BLD_POWERPLANT: {
		maxHitpoints:  500,
		w:             2,
		h:             2,
		displayedName: "Power Plant",
		cost:          500,
		buildTime:     5,
		builds:        nil,
		produces:      nil,
		givesEnergy:   20,
		hotkeyToBuild: "P",
	},
	BLD_FACTORY: {
		maxHitpoints:  750,
		w:             3,
		h:             2,
		displayedName: "Factory",
		cost:          1000,
		buildTime:     7,
		builds:        nil,
		produces:      []int{UNT_TANK, UNT_QUAD, UNT_MSLTANK, UNT_HARVESTER},
		hotkeyToBuild: "F",
	},
	BLD_REFINERY: {
		maxHitpoints:  550,
		w:             3,
		h:             2,
		displayedName: "Refinery",
		cost:          2000,
		buildTime:     1,
		builds:        nil,
		produces:      nil,
		hotkeyToBuild: "R",

		receivesResources: true,
		storageAmount:     1000,
		unitPlacementX:    1, unitPlacementY: 1,

		givesFreeUnitOnCreation:   true,
		codeForFreeUnitOnCreation: UNT_HARVESTER,
	},
	BLD_TURRET: {
		maxHitpoints:  750,
		w:             1,
		h:             1,
		displayedName: "Defense tower",
		cost:          750,
		buildTime:     10,
		builds:        nil,
		produces:      nil,
		turretCode:    TRT_CANNON_BUILDING,
		hotkeyToBuild: "T",
	},
	BLD_SILO: {
		maxHitpoints:  500,
		w:             1,
		h:             2,
		displayedName: "Resource Silo",
		cost:          500,
		buildTime:     7,
		builds:        nil,
		produces:      nil,
		hotkeyToBuild: "S",
		storageAmount: 2500,
	},
}
