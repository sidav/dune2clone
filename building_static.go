package main

const (
	BLD_BASE = iota
	BLD_POWERPLANT
	BLD_FACTORY
	BLD_AIRFACTORY
	BLD_TURRET_MINIGUN
	BLD_TURRET_CANNON
	BLD_REFINERY
	BLD_SILO
	BLD_FORTRESS
)

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
	// CanUnitBePlacedHere            bool // Removed, implemented in method
	unitPlacementX, unitPlacementY int // tile coords for placed unit draw

	givesFreeUnitOnCreation   bool
	codeForFreeUnitOnCreation int

	givesEnergy, consumesEnergy int
	storageAmount               float64

	spriteCode string

	// ui-only things:
	hotkeyToBuild string
}

func (bs *buildingStatic) canUnitBePlacedIn() bool {
	return bs.receivesResources // TODO: update when needed
}

var sTableBuildings = map[int]*buildingStatic{
	BLD_BASE: {
		spriteCode:    "base",
		maxHitpoints:  1000,
		w:             2,
		h:             2,
		displayedName: "Construction Yard",
		cost:          2500,
		buildTime:     30,
		builds: []int{BLD_BASE, BLD_POWERPLANT, BLD_REFINERY, BLD_FACTORY, BLD_AIRFACTORY,
			BLD_TURRET_CANNON, BLD_TURRET_MINIGUN, BLD_SILO, BLD_FORTRESS},
		givesEnergy:   10,
		hotkeyToBuild: "B",
	},
	BLD_POWERPLANT: {
		spriteCode:    "powerplant",
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
		spriteCode:     "factory",
		maxHitpoints:   750,
		w:              3,
		h:              2,
		displayedName:  "Factory",
		cost:           1000,
		buildTime:      7,
		builds:         nil,
		consumesEnergy: 15,
		produces:       []int{UNT_TANK, UNT_QUAD, UNT_MSLTANK, UNT_HARVESTER},
		hotkeyToBuild:  "F",
	},
	BLD_AIRFACTORY: {
		spriteCode:     "airfactory",
		maxHitpoints:   750,
		w:              3,
		h:              2,
		displayedName:  "Aircraft Factory",
		cost:           1000,
		buildTime:      10,
		produces:       []int{AIR_TRANSPORT, AIR_GUNSHIP},
		consumesEnergy: 15,
		// produces:       []int{UNT_TANK, UNT_QUAD, UNT_MSLTANK, UNT_HARVESTER},
		hotkeyToBuild: "A",
	},
	BLD_REFINERY: {
		spriteCode:     "refinery",
		maxHitpoints:   550,
		w:              3,
		h:              2,
		displayedName:  "Refinery",
		cost:           2000,
		buildTime:      10,
		builds:         nil,
		consumesEnergy: 10,
		produces:       nil,
		hotkeyToBuild:  "R",

		receivesResources: true,
		storageAmount:     1000,
		unitPlacementX:    1, unitPlacementY: 1,

		givesFreeUnitOnCreation:   true,
		codeForFreeUnitOnCreation: UNT_HARVESTER,
	},
	BLD_TURRET_MINIGUN: {
		spriteCode:     "turret_base",
		maxHitpoints:   500,
		w:              1,
		h:              1,
		displayedName:  "Minigun tower",
		cost:           550,
		buildTime:      10,
		consumesEnergy: 5,
		turretCode:     TRT_MINIGUN_BUILDING,
		hotkeyToBuild:  "M",
	},
	BLD_TURRET_CANNON: {
		spriteCode:     "turret_base",
		maxHitpoints:   500,
		w:              1,
		h:              1,
		displayedName:  "Heavy tower",
		cost:           750,
		buildTime:      10,
		consumesEnergy: 8,
		turretCode:     TRT_CANNON_BUILDING,
		hotkeyToBuild:  "T",
	},
	BLD_SILO: {
		spriteCode:     "silo",
		maxHitpoints:   500,
		w:              1,
		h:              2,
		displayedName:  "Resource Silo",
		cost:           500,
		buildTime:      7,
		consumesEnergy: 5,
		hotkeyToBuild:  "S",
		storageAmount:  2500,
	},
	BLD_FORTRESS: {
		spriteCode:     "fortress",
		maxHitpoints:   1500,
		w:              2,
		h:              2,
		displayedName:  "Fortress",
		cost:           2500,
		buildTime:      15,
		builds:         nil,
		produces:       nil,
		consumesEnergy: 20,
		hotkeyToBuild:  "O",
		turretCode:     TRT_BUILDING_FORTRESS,
	},
}