package main

type buildTypeCode int
type buildingCode int

const (
	BLD_CONYARD1 buildingCode = iota
	BLD_CONYARD2
	BLD_POWERPLANT1
	BLD_POWERPLANT2
	BLD_FUSION
	BLD_BARRACKS
	BLD_FACTORY
	BLD_REPAIR_DEPOT
	BLD_AIRFACTORY
	BLD_TURRET_MINIGUN
	BLD_TURRET_CANNON
	BLD_REFINERY
	BLD_SILO
	BLD_FORTRESS

	BTYPE_BUILD_FIRST buildTypeCode = iota // like in Dune/C&C series
	BTYPE_PLACE_FIRST
)

type buildingStatic struct {
	w, h              int
	displayedName     string
	cost              int
	buildTime         int // seconds
	requiresToBeBuilt []buildingCode
	requiresTechLevel int
	givesTechLevel    int

	builds       []buildingCode // buildings
	buildType    buildTypeCode
	produces     []int // units
	maxHitpoints int

	turretCode int

	receivesResources                  bool // is refinery
	repairsUnits                       bool
	unitPlacementX, unitPlacementY     int // tile coords for placed unit draw
	needsEmptyRowBelowWhenConstructing bool

	givesFreeUnitOnCreation   bool
	codeForFreeUnitOnCreation int

	givesEnergy, consumesEnergy int
	storageAmount               float64

	spriteCode string

	// ui-only things:
	hotkeyToBuild string
}

func (bs *buildingStatic) canUnitBePlacedIn() bool {
	return bs.receivesResources || bs.repairsUnits // TODO: update when needed
}

var sTableBuildings = map[buildingCode]*buildingStatic{
	BLD_CONYARD1: {
		spriteCode:    "base",
		maxHitpoints:  1000,
		w:             2,
		h:             2,
		displayedName: "BetaCorp Construction Yard",
		cost:          2500,
		buildTime:     30,
		builds: []buildingCode{BLD_POWERPLANT1, BLD_FUSION, BLD_BARRACKS, BLD_REFINERY, BLD_FACTORY, BLD_REPAIR_DEPOT, BLD_AIRFACTORY,
			BLD_TURRET_CANNON, BLD_TURRET_MINIGUN, BLD_SILO, BLD_FORTRESS},
		buildType:     BTYPE_BUILD_FIRST, //BTYPE_PLACE_FIRST,
		givesEnergy:   10,
		hotkeyToBuild: "Y",
	},
	BLD_CONYARD2: {
		spriteCode:    "base",
		maxHitpoints:  1000,
		w:             2,
		h:             2,
		displayedName: "Commonwealth Construction Yard",
		cost:          2500,
		buildTime:     30,
		builds: []buildingCode{BLD_POWERPLANT2, BLD_FUSION, BLD_BARRACKS, BLD_REFINERY, BLD_FACTORY, BLD_REPAIR_DEPOT, BLD_AIRFACTORY,
			BLD_TURRET_CANNON, BLD_TURRET_MINIGUN, BLD_SILO, BLD_FORTRESS},
		buildType:     BTYPE_BUILD_FIRST, //BTYPE_PLACE_FIRST,
		givesEnergy:   10,
		hotkeyToBuild: "Y",
	},
	BLD_POWERPLANT1: {
		spriteCode:    "powerplant1",
		maxHitpoints:  500,
		w:             2,
		h:             2,
		displayedName: "Plasma Reactor",
		cost:          600,
		buildTime:     6,
		builds:        nil,
		produces:      nil,
		givesEnergy:   25,
		hotkeyToBuild: "P",
	},
	BLD_POWERPLANT2: {
		spriteCode:    "powerplant2",
		maxHitpoints:  500,
		w:             2,
		h:             2,
		displayedName: "RITEG power plant",
		cost:          500,
		buildTime:     5,
		builds:        nil,
		produces:      nil,
		givesEnergy:   20,
		hotkeyToBuild: "P",
	},
	BLD_FUSION: {
		spriteCode:        "fusionreactor",
		maxHitpoints:      1000,
		w:                 3,
		h:                 3,
		requiresTechLevel: 5,
		// requiresToBeBuilt: []buildingCode{BLD_POWERPLANT1},
		displayedName: "Fusion Reactor",
		cost:          3000,
		buildTime:     45,
		builds:        nil,
		produces:      nil,
		givesEnergy:   200,
		hotkeyToBuild: "L",
	},
	BLD_BARRACKS: {
		spriteCode:                         "barracks",
		maxHitpoints:                       500,
		w:                                  2,
		h:                                  2,
		displayedName:                      "Barracks",
		needsEmptyRowBelowWhenConstructing: true,
		cost:                               500,
		buildTime:                          10,
		requiresTechLevel:                  2,
		builds:                             nil,
		produces:                           []int{UNT_INFANTRY},
		hotkeyToBuild:                      "B",
	},
	BLD_FACTORY: {
		spriteCode:                         "factory",
		maxHitpoints:                       750,
		w:                                  3,
		h:                                  2,
		displayedName:                      "Factory",
		requiresTechLevel:                  2,
		requiresToBeBuilt:                  []buildingCode{BLD_REFINERY},
		needsEmptyRowBelowWhenConstructing: true,
		cost:                               1000,
		buildTime:                          12,
		builds:                             nil,
		consumesEnergy:                     15,
		produces:                           []int{UNT_TANK2, UNT_MCV1, UNT_QUAD, UNT_MSLTANK, UNT_AATANK, UNT_HARVESTER},
		hotkeyToBuild:                      "F",
	},
	BLD_REPAIR_DEPOT: {
		spriteCode:                         "depot",
		maxHitpoints:                       750,
		w:                                  3,
		h:                                  2,
		displayedName:                      "Repair depot",
		needsEmptyRowBelowWhenConstructing: true,
		cost:                               750,
		buildTime:                          7,
		builds:                             nil,
		repairsUnits:                       true,
		unitPlacementX:                     1, unitPlacementY: 1,
		consumesEnergy: 10,
		hotkeyToBuild:  "D",
	},
	BLD_AIRFACTORY: {
		spriteCode:     "airfactory",
		maxHitpoints:   750,
		w:              2,
		h:              3,
		displayedName:  "Aircraft Factory",
		cost:           1000,
		buildTime:      10,
		produces:       []int{AIR_TRANSPORT, AIR_GUNSHIP},
		consumesEnergy: 15,
		// produces:       []int{UNT_TANK, UNT_QUAD, UNT_MSLTANK, UNT_HARVESTER},
		hotkeyToBuild: "A",
	},
	BLD_REFINERY: {
		spriteCode:                         "refinery",
		maxHitpoints:                       550,
		w:                                  3,
		h:                                  2,
		displayedName:                      "Refinery",
		needsEmptyRowBelowWhenConstructing: true,
		cost:                               2000,
		buildTime:                          10,
		builds:                             nil,
		consumesEnergy:                     10,
		produces:                           nil,
		hotkeyToBuild:                      "R",

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
