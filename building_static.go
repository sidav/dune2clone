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
	BLD_FACTORY1
	BLD_FACTORY2
	BLD_REPAIR_DEPOT
	BLD_AIRFACTORY
	BLD_TURRET_MINIGUN
	BLD_TURRET_CANNON
	BLD_TURRET_AA
	BLD_REFINERY1
	BLD_REFINERY2
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

	turretData *turretStatic

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
	// general
	BLD_FUSION: {
		spriteCode:        "fusionreactor",
		maxHitpoints:      1000,
		w:                 3,
		h:                 3,
		requiresTechLevel: 4,
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
		produces:                           []int{UNT_INFANTRY, UNT_RECONINFANTRY, UNT_ROCKETINFANTRY, UNT_HEAVYINFANTRY},
		hotkeyToBuild:                      "B",
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
		requiresTechLevel:                  2,
		repairsUnits:                       true,
		unitPlacementX:                     1, unitPlacementY: 1,
		consumesEnergy: 10,
		hotkeyToBuild:  "D",
	},
	BLD_AIRFACTORY: {
		spriteCode:        "airfactory",
		maxHitpoints:      750,
		w:                 2,
		h:                 3,
		displayedName:     "Aircraft Factory",
		cost:              1000,
		requiresTechLevel: 3,
		givesTechLevel:    4,
		buildTime:         10,
		produces:          []int{AIR_TRANSPORT, AIR_GUNSHIP, AIR_FIGHTER},
		consumesEnergy:    15,
		// produces:       []int{UNT_TANK, UNT_QUAD, UNT_MSLTANK, UNT_HARVESTER},
		hotkeyToBuild: "A",
	},
	BLD_TURRET_MINIGUN: {
		spriteCode:        "turret_base",
		maxHitpoints:      150,
		w:                 1,
		h:                 1,
		displayedName:     "Minigun tower",
		cost:              550,
		buildTime:         10,
		consumesEnergy:    5,
		requiresTechLevel: 1,
		turretData: &turretStatic{
			spriteCode:        "bld_turret_minigun",
			attacksLand:       true,
			attacksAir:        true,
			rotateSpeed:       17,
			fireRange:         6,
			fireSpreadDegrees: 7,
			shotRangeSpread:   0.7,
			attackCooldown:    5,
			firedProjectileData: &projectileStatic{
				spriteCode:                "bullets",
				hitDamage:                 4,
				damageType:                DAMAGETYPE_ANTI_INFANTRY,
				size:                      0.2,
				speed:                     0.7,
				createsEffectOnImpact:     true,
				effectCreatedOnImpactCode: EFFECT_SMALL_EXPLOSION,
			},
		},
		hotkeyToBuild: "M",
	},
	BLD_TURRET_CANNON: {
		spriteCode:        "turret_base",
		maxHitpoints:      200,
		w:                 1,
		h:                 1,
		displayedName:     "Heavy tower",
		cost:              750,
		buildTime:         10,
		consumesEnergy:    8,
		requiresTechLevel: 2,
		turretData: &turretStatic{
			spriteCode:        "bld_turret_cannon",
			attacksLand:       true,
			rotateSpeed:       15,
			fireRange:         6,
			fireSpreadDegrees: 7,
			shotRangeSpread:   0.7,
			attackCooldown:    50,
			firedProjectileData: &projectileStatic{
				spriteCode:                "shell",
				size:                      0.3,
				splashRadius:              0.25,
				splashDamage:              15,
				damageType:                DAMAGETYPE_HEAVY,
				speed:                     0.7,
				createsEffectOnImpact:     true,
				effectCreatedOnImpactCode: EFFECT_REGULAR_EXPLOSION,
			},
		},
		hotkeyToBuild: "T",
	},
	BLD_TURRET_AA: {
		spriteCode:        "bld_aaturret",
		maxHitpoints:      250,
		w:                 1,
		h:                 1,
		displayedName:     "AA SAM site",
		cost:              750,
		buildTime:         15,
		requiresTechLevel: 3,
		builds:            nil,
		produces:          nil,
		consumesEnergy:    15,
		hotkeyToBuild:     "M",
		turretData: &turretStatic{
			spriteCode:        "",
			attacksLand:       false,
			attacksAir:        true,
			rotateSpeed:       180,
			fireRange:         8,
			fireSpreadDegrees: 30,
			shotRangeSpread:   0.3,
			attackCooldown:    100,
			firedProjectileData: &projectileStatic{
				spriteCode:                "aamissile",
				size:                      0.3,
				hitDamage:                 20,
				damageType:                DAMAGETYPE_HEAVY,
				speed:                     0.45,
				rotationSpeed:             35,
				createsEffectOnImpact:     true,
				effectCreatedOnImpactCode: EFFECT_SMALL_EXPLOSION,
			},
		},
	},
	BLD_SILO: {
		spriteCode:        "silo",
		maxHitpoints:      500,
		w:                 1,
		h:                 2,
		displayedName:     "Resource Silo",
		cost:              500,
		buildTime:         7,
		consumesEnergy:    5,
		requiresTechLevel: 2,
		hotkeyToBuild:     "S",
		storageAmount:     2500,
	},
	BLD_FORTRESS: {
		spriteCode:        "fortress",
		maxHitpoints:      500,
		w:                 2,
		h:                 2,
		displayedName:     "Fortress",
		cost:              2500,
		buildTime:         15,
		requiresTechLevel: 3,
		builds:            nil,
		produces:          nil,
		consumesEnergy:    20,
		hotkeyToBuild:     "O",
		turretData: &turretStatic{
			spriteCode:        "bld_fortress_cannon",
			attacksLand:       true,
			rotateSpeed:       5,
			fireRange:         15,
			fireSpreadDegrees: 5,
			shotRangeSpread:   0.3,
			attackCooldown:    80,
			firedProjectileData: &projectileStatic{
				spriteCode:                "shell",
				size:                      0.3,
				splashRadius:              0.35,
				splashDamage:              15,
				damageType:                DAMAGETYPE_HEAVY,
				speed:                     0.7,
				createsEffectOnImpact:     true,
				effectCreatedOnImpactCode: EFFECT_REGULAR_EXPLOSION,
			},
		},
	},

	// faction 1
	BLD_CONYARD1: {
		spriteCode:     "base",
		maxHitpoints:   1000,
		w:              2,
		h:              2,
		displayedName:  "BetaCorp Construction Yard",
		cost:           2500,
		buildTime:      30,
		givesTechLevel: 1,
		builds: []buildingCode{BLD_POWERPLANT1, BLD_FUSION, BLD_BARRACKS, BLD_REFINERY1, BLD_FACTORY1, BLD_REPAIR_DEPOT, BLD_AIRFACTORY,
			BLD_TURRET_CANNON, BLD_TURRET_MINIGUN, BLD_SILO, BLD_FORTRESS, BLD_TURRET_AA},
		buildType:     BTYPE_BUILD_FIRST, //BTYPE_PLACE_FIRST,
		givesEnergy:   10,
		hotkeyToBuild: "Y",
	},
	BLD_POWERPLANT1: {
		spriteCode:     "powerplant1",
		maxHitpoints:   500,
		w:              2,
		h:              2,
		displayedName:  "Plasma Reactor",
		cost:           600,
		buildTime:      6,
		givesTechLevel: 2,
		builds:         nil,
		produces:       nil,
		givesEnergy:    25,
		hotkeyToBuild:  "P",
	},
	BLD_REFINERY1: {
		spriteCode:                         "refinery",
		maxHitpoints:                       550,
		w:                                  3,
		h:                                  2,
		displayedName:                      "BetaCorp Refinery",
		needsEmptyRowBelowWhenConstructing: true,
		cost:                               2000,
		buildTime:                          10,
		builds:                             nil,
		consumesEnergy:                     10,
		produces:                           nil,
		hotkeyToBuild:                      "R",
		givesTechLevel:                     2,

		receivesResources: true,
		storageAmount:     1000,
		unitPlacementX:    1, unitPlacementY: 1,

		givesFreeUnitOnCreation:   true,
		codeForFreeUnitOnCreation: UNT_FAST_HARVESTER,
	},
	BLD_FACTORY1: {
		spriteCode:                         "factory",
		maxHitpoints:                       750,
		w:                                  3,
		h:                                  2,
		displayedName:                      "Factory",
		requiresTechLevel:                  2,
		requiresToBeBuilt:                  []buildingCode{BLD_REFINERY1},
		needsEmptyRowBelowWhenConstructing: true,
		cost:                               1000,
		givesTechLevel:                     3,
		buildTime:                          12,
		builds:                             nil,
		consumesEnergy:                     15,
		produces:                           []int{UNT_TANK2, UNT_DEVASTATOR, UNT_MCV1, UNT_QUAD, UNT_MSLTANK, UNT_AATANK, UNT_FAST_HARVESTER},
		hotkeyToBuild:                      "F",
	},

	// FACTION 2
	BLD_CONYARD2: {
		spriteCode:     "base",
		maxHitpoints:   1000,
		w:              2,
		h:              2,
		displayedName:  "Commonwealth Construction Yard",
		cost:           2500,
		buildTime:      30,
		givesTechLevel: 1,
		builds: []buildingCode{BLD_POWERPLANT2, BLD_FUSION, BLD_BARRACKS, BLD_REFINERY2, BLD_FACTORY2, BLD_REPAIR_DEPOT, BLD_AIRFACTORY,
			BLD_TURRET_CANNON, BLD_TURRET_MINIGUN, BLD_SILO, BLD_FORTRESS, BLD_TURRET_AA},
		buildType:     BTYPE_PLACE_FIRST, //BTYPE_PLACE_FIRST,
		givesEnergy:   10,
		hotkeyToBuild: "Y",
	},
	BLD_POWERPLANT2: {
		spriteCode:     "powerplant2",
		maxHitpoints:   500,
		w:              2,
		h:              2,
		displayedName:  "RITEG power plant",
		cost:           500,
		buildTime:      5,
		builds:         nil,
		produces:       nil,
		givesEnergy:    20,
		hotkeyToBuild:  "P",
		givesTechLevel: 2,
	},
	BLD_REFINERY2: {
		spriteCode:                         "refinery",
		maxHitpoints:                       550,
		w:                                  3,
		h:                                  2,
		displayedName:                      "BetaCorp Refinery",
		needsEmptyRowBelowWhenConstructing: true,
		cost:                               2000,
		buildTime:                          10,
		builds:                             nil,
		consumesEnergy:                     10,
		produces:                           nil,
		hotkeyToBuild:                      "R",
		givesTechLevel:                     2,

		receivesResources: true,
		storageAmount:     1000,
		unitPlacementX:    1, unitPlacementY: 1,

		givesFreeUnitOnCreation:   true,
		codeForFreeUnitOnCreation: UNT_COMBAT_HARVESTER,
	},
	BLD_FACTORY2: {
		spriteCode:                         "factory",
		maxHitpoints:                       750,
		w:                                  3,
		h:                                  2,
		displayedName:                      "Factory",
		requiresTechLevel:                  2,
		requiresToBeBuilt:                  []buildingCode{BLD_REFINERY2},
		needsEmptyRowBelowWhenConstructing: true,
		cost:                               1000,
		givesTechLevel:                     3,
		buildTime:                          12,
		builds:                             nil,
		consumesEnergy:                     15,
		produces:                           []int{UNT_TANK2, UNT_DEVASTATOR, UNT_MCV2, UNT_QUAD, UNT_MSLTANK, UNT_AATANK, UNT_COMBAT_HARVESTER},
		hotkeyToBuild:                      "F",
	},
}
