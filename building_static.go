package main

type buildTypeCode int
type buildingCode int

const (
	BLD_NULL buildingCode = iota
	BLD_CONYARD1
	BLD_CONYARD2
	BLD_POWERPLANT1
	BLD_POWERPLANT2
	BLD_FUSION
	BLD_BARRACKS1
	BLD_BARRACKS2
	BLD_FACTORY1
	BLD_FACTORY2
	BLD_REPAIR_DEPOT
	BLD_AIRFACTORY1
	BLD_AIRFACTORY2
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
	W                 int            `json:"w,omitempty"`
	H                 int            `json:"h,omitempty"`
	DisplayedName     string         `json:"displayed_name,omitempty"`
	Cost              int            `json:"cost,omitempty"`
	BuildTime         int            `json:"build_time,omitempty"` // seconds
	RequiresToBeBuilt []buildingCode `json:"requires_to_be_built,omitempty"`
	RequiresTechLevel int            `json:"requires_tech_level,omitempty"`
	GivesTechLevel    int            `json:"gives_tech_level,omitempty"`

	Builds       []buildingCode `json:"builds,omitempty"` // buildings
	BuildType    buildTypeCode  `json:"build_type,omitempty"`
	Produces     []int          `json:"produces,omitempty"` // units
	MaxHitpoints int            `json:"max_hitpoints,omitempty"`

	TurretData *TurretStatic `json:"turret_data,omitempty"`

	ReceivesResources                  bool `json:"receives_resources,omitempty"` // is refinery
	RepairsUnits                       bool `json:"repairs_units,omitempty"`
	UnitPlacementX                     int  `json:"unit_placement_x,omitempty"`
	UnitPlacementY                     int  `json:"unit_placement_y,omitempty"` // tile coords for placed unit draw
	NeedsEmptyRowBelowWhenConstructing bool `json:"needs_empty_row_below_when_constructing,omitempty"`

	GivesFreeUnitOnCreation   bool `json:"gives_free_unit_on_creation,omitempty"`
	CodeForFreeUnitOnCreation int  `json:"code_for_free_unit_on_creation,omitempty"`

	GivesEnergy    int     `json:"gives_energy,omitempty"`
	ConsumesEnergy int     `json:"consumes_energy,omitempty"`
	StorageAmount  float64 `json:"storage_amount,omitempty"`

	SpriteCode string `json:"sprite_code,omitempty"`

	// ui-only things:
	HotkeyToBuild string `json:"hotkey_to_build,omitempty"`
}

func (bs *buildingStatic) canUnitBePlacedIn() bool {
	return bs.ReceivesResources || bs.RepairsUnits // TODO: update when needed
}

var sTableBuildings = map[buildingCode]*buildingStatic{
	// general
	BLD_FUSION: {
		SpriteCode:        "fusionreactor",
		MaxHitpoints:      1000,
		W:                 3,
		H:                 3,
		RequiresTechLevel: 4,
		// requiresToBeBuilt: []buildingCode{BLD_POWERPLANT1},
		DisplayedName: "Fusion Reactor",
		Cost:          3000,
		BuildTime:     45,
		Builds:        nil,
		Produces:      nil,
		GivesEnergy:   200,
		HotkeyToBuild: "L",
	},
	BLD_REPAIR_DEPOT: {
		SpriteCode:                         "depot",
		MaxHitpoints:                       750,
		W:                                  3,
		H:                                  2,
		DisplayedName:                      "Repair depot",
		NeedsEmptyRowBelowWhenConstructing: true,
		Cost:                               750,
		BuildTime:                          7,
		RequiresTechLevel:                  2,
		RepairsUnits:                       true,
		UnitPlacementX:                     1, UnitPlacementY: 1,
		ConsumesEnergy: 10,
		HotkeyToBuild:  "D",
	},
	BLD_TURRET_MINIGUN: {
		SpriteCode:        "turret_base",
		MaxHitpoints:      140,
		W:                 1,
		H:                 1,
		DisplayedName:     "Minigun tower",
		Cost:              550,
		BuildTime:         15,
		ConsumesEnergy:    5,
		RequiresTechLevel: 1,
		TurretData: &TurretStatic{
			SpriteCode:        "bld_turret_minigun",
			AttacksLand:       true,
			AttacksAir:        true,
			RotateSpeed:       17,
			FireRange:         6,
			FireSpreadDegrees: 7,
			ShotRangeSpread:   0.7,
			AttackCooldown:    8,
			FiredProjectileData: &projectileStatic{
				SpriteCode:                "bullets",
				HitDamage:                 4,
				DamageType:                DAMAGETYPE_ANTI_INFANTRY,
				Size:                      0.2,
				Speed:                     0.7,
				RotationSpeed:             1,
				CreatesEffectOnImpact:     true,
				EffectCreatedOnImpactCode: EFFECT_SMALL_EXPLOSION,
			},
		},
		HotkeyToBuild: "G",
	},
	BLD_TURRET_CANNON: {
		SpriteCode:        "turret_base",
		MaxHitpoints:      200,
		W:                 1,
		H:                 1,
		DisplayedName:     "Heavy tower",
		Cost:              750,
		BuildTime:         20,
		ConsumesEnergy:    8,
		RequiresTechLevel: 2,
		TurretData: &TurretStatic{
			SpriteCode:        "bld_turret_cannon",
			AttacksLand:       true,
			RotateSpeed:       15,
			FireRange:         6,
			FireSpreadDegrees: 6,
			ShotRangeSpread:   0.6,
			AttackCooldown:    45,
			FiredProjectileData: &projectileStatic{
				SpriteCode:                "shell",
				Size:                      0.3,
				SplashRadius:              0.25,
				SplashDamage:              15,
				DamageType:                DAMAGETYPE_HEAVY,
				Speed:                     0.7,
				CreatesEffectOnImpact:     true,
				EffectCreatedOnImpactCode: EFFECT_REGULAR_EXPLOSION,
			},
		},
		HotkeyToBuild: "T",
	},
	BLD_TURRET_AA: {
		SpriteCode:        "bld_aaturret",
		MaxHitpoints:      150,
		W:                 1,
		H:                 1,
		DisplayedName:     "AA SAM site",
		Cost:              750,
		BuildTime:         25,
		RequiresTechLevel: 3,
		Builds:            nil,
		Produces:          nil,
		ConsumesEnergy:    15,
		HotkeyToBuild:     "M",
		TurretData: &TurretStatic{
			SpriteCode:        "",
			AttacksLand:       false,
			AttacksAir:        true,
			RotateSpeed:       180,
			FireRange:         12,
			FireSpreadDegrees: 30,
			ShotRangeSpread:   0.3,
			AttackCooldown:    200,
			FiredProjectileData: &projectileStatic{
				SpriteCode:                "aamissile",
				Size:                      0.3,
				HitDamage:                 17,
				DamageType:                DAMAGETYPE_HEAVY,
				Speed:                     0.45,
				RotationSpeed:             35,
				CreatesEffectOnImpact:     true,
				EffectCreatedOnImpactCode: EFFECT_SMALL_EXPLOSION,
			},
		},
	},
	BLD_SILO: {
		SpriteCode:        "silo",
		MaxHitpoints:      500,
		W:                 1,
		H:                 2,
		DisplayedName:     "Resource Silo",
		Cost:              500,
		BuildTime:         7,
		ConsumesEnergy:    5,
		RequiresTechLevel: 2,
		HotkeyToBuild:     "S",
		StorageAmount:     2500,
	},
	BLD_FORTRESS: {
		SpriteCode:        "fortress",
		MaxHitpoints:      500,
		W:                 2,
		H:                 2,
		DisplayedName:     "Fortress",
		Cost:              2500,
		BuildTime:         15,
		RequiresTechLevel: 3,
		Builds:            nil,
		Produces:          nil,
		ConsumesEnergy:    20,
		HotkeyToBuild:     "O",
		TurretData: &TurretStatic{
			SpriteCode:        "bld_fortress_cannon",
			AttacksLand:       true,
			RotateSpeed:       5,
			FireRange:         15,
			FireSpreadDegrees: 5,
			ShotRangeSpread:   0.3,
			AttackCooldown:    80,
			FiredProjectileData: &projectileStatic{
				SpriteCode:                "shell",
				Size:                      0.3,
				SplashRadius:              0.35,
				SplashDamage:              15,
				DamageType:                DAMAGETYPE_HEAVY,
				Speed:                     0.7,
				CreatesEffectOnImpact:     true,
				EffectCreatedOnImpactCode: EFFECT_REGULAR_EXPLOSION,
			},
		},
	},

	// faction 1
	BLD_CONYARD1: {
		SpriteCode:     "conyard",
		MaxHitpoints:   1000,
		W:              2,
		H:              2,
		DisplayedName:  "BetaCorp Construction Yard",
		Cost:           2500,
		BuildTime:      30,
		GivesTechLevel: 1,
		Builds: []buildingCode{BLD_POWERPLANT1, BLD_FUSION, BLD_BARRACKS1, BLD_REFINERY1, BLD_FACTORY1, BLD_REPAIR_DEPOT, BLD_AIRFACTORY1,
			BLD_TURRET_CANNON, BLD_TURRET_MINIGUN, BLD_SILO, BLD_FORTRESS, BLD_TURRET_AA},
		BuildType:     BTYPE_BUILD_FIRST, //BTYPE_PLACE_FIRST,
		GivesEnergy:   10,
		HotkeyToBuild: "Y",
	},
	BLD_POWERPLANT1: {
		SpriteCode:     "powerplant1",
		MaxHitpoints:   500,
		W:              2,
		H:              2,
		DisplayedName:  "Plasma Reactor",
		Cost:           600,
		BuildTime:      6,
		GivesTechLevel: 2,
		Builds:         nil,
		Produces:       nil,
		GivesEnergy:    25,
		HotkeyToBuild:  "P",
	},
	BLD_REFINERY1: {
		SpriteCode:                         "refinery1",
		MaxHitpoints:                       550,
		W:                                  3,
		H:                                  2,
		DisplayedName:                      "BetaCorp Refinery",
		NeedsEmptyRowBelowWhenConstructing: true,
		Cost:                               2000,
		BuildTime:                          15,
		Builds:                             nil,
		ConsumesEnergy:                     10,
		Produces:                           nil,
		HotkeyToBuild:                      "R",
		GivesTechLevel:                     2,

		ReceivesResources: true,
		StorageAmount:     1000,
		UnitPlacementX:    1, UnitPlacementY: 1,

		GivesFreeUnitOnCreation:   true,
		CodeForFreeUnitOnCreation: UNT_FAST_HARVESTER,
	},
	BLD_BARRACKS1: {
		SpriteCode:                         "barracks",
		MaxHitpoints:                       500,
		W:                                  2,
		H:                                  2,
		DisplayedName:                      "Barracks",
		NeedsEmptyRowBelowWhenConstructing: true,
		Cost:                               500,
		BuildTime:                          10,
		RequiresTechLevel:                  2,
		ConsumesEnergy:                     7,
		Builds:                             nil,
		Produces:                           []int{UNT_INFANTRY, UNT_RECONINFANTRY, UNT_ROCKETINFANTRY, UNT_HEAVYINFANTRY},
		HotkeyToBuild:                      "B",
	},
	BLD_FACTORY1: {
		SpriteCode:                         "factory1",
		MaxHitpoints:                       750,
		W:                                  3,
		H:                                  2,
		DisplayedName:                      "Factory",
		RequiresTechLevel:                  2,
		RequiresToBeBuilt:                  []buildingCode{BLD_REFINERY1},
		NeedsEmptyRowBelowWhenConstructing: true,
		Cost:                               1000,
		GivesTechLevel:                     3,
		BuildTime:                          12,
		Builds:                             nil,
		ConsumesEnergy:                     15,
		Produces:                           []int{UNT_TANK1, UNT_DEVASTATOR, UNT_MCV1, UNT_QUAD, UNT_MSLTANK, UNT_AATANK, UNT_FAST_HARVESTER},
		HotkeyToBuild:                      "F",
	},
	BLD_AIRFACTORY1: {
		SpriteCode:        "airfactory",
		MaxHitpoints:      750,
		W:                 2,
		H:                 3,
		DisplayedName:     "Aircraft Factory",
		Cost:              1000,
		RequiresTechLevel: 3,
		GivesTechLevel:    4,
		BuildTime:         10,
		Produces:          []int{AIR_TRANSPORT1, AIR_GUNSHIP, AIR_FIGHTER, AIR_FORTRESS},
		ConsumesEnergy:    15,
		// produces:       []int{UNT_TANK, UNT_QUAD, UNT_MSLTANK, UNT_HARVESTER},
		HotkeyToBuild: "A",
	},

	// FACTION 2
	BLD_CONYARD2: {
		SpriteCode:     "conyard2",
		MaxHitpoints:   1000,
		W:              2,
		H:              2,
		DisplayedName:  "Commonwealth Construction Yard",
		Cost:           2500,
		BuildTime:      30,
		GivesTechLevel: 1,
		Builds: []buildingCode{BLD_POWERPLANT2, BLD_FUSION, BLD_BARRACKS2, BLD_REFINERY2, BLD_FACTORY2, BLD_REPAIR_DEPOT, BLD_AIRFACTORY2,
			BLD_TURRET_CANNON, BLD_TURRET_MINIGUN, BLD_SILO, BLD_FORTRESS, BLD_TURRET_AA},
		BuildType:     BTYPE_PLACE_FIRST, //BTYPE_PLACE_FIRST,
		GivesEnergy:   10,
		HotkeyToBuild: "Y",
	},
	BLD_POWERPLANT2: {
		SpriteCode:     "powerplant2",
		MaxHitpoints:   500,
		W:              2,
		H:              2,
		DisplayedName:  "RITEG power plant",
		Cost:           500,
		BuildTime:      5,
		Builds:         nil,
		Produces:       nil,
		GivesEnergy:    20,
		HotkeyToBuild:  "P",
		GivesTechLevel: 2,
	},
	BLD_REFINERY2: {
		SpriteCode:                         "refinery2",
		MaxHitpoints:                       600,
		W:                                  3,
		H:                                  2,
		DisplayedName:                      "Commonwealth Refinery",
		NeedsEmptyRowBelowWhenConstructing: true,
		Cost:                               2000,
		BuildTime:                          15,
		Builds:                             nil,
		ConsumesEnergy:                     10,
		Produces:                           nil,
		HotkeyToBuild:                      "R",
		GivesTechLevel:                     2,

		ReceivesResources: true,
		StorageAmount:     1000,
		UnitPlacementX:    1, UnitPlacementY: 1,

		GivesFreeUnitOnCreation:   true,
		CodeForFreeUnitOnCreation: UNT_COMBAT_HARVESTER,
	},
	BLD_BARRACKS2: {
		SpriteCode:                         "barracks2",
		MaxHitpoints:                       500,
		W:                                  2,
		H:                                  2,
		DisplayedName:                      "Barracks",
		NeedsEmptyRowBelowWhenConstructing: true,
		Cost:                               500,
		BuildTime:                          10,
		RequiresTechLevel:                  2,
		ConsumesEnergy:                     7,
		Builds:                             nil,
		Produces:                           []int{UNT_INFANTRY, UNT_RECONINFANTRY, UNT_ROCKETINFANTRY, UNT_SNIPERINFANTRY},
		HotkeyToBuild:                      "B",
	},
	BLD_FACTORY2: {
		SpriteCode:                         "factory2",
		MaxHitpoints:                       750,
		W:                                  3,
		H:                                  2,
		DisplayedName:                      "Factory",
		RequiresTechLevel:                  2,
		RequiresToBeBuilt:                  []buildingCode{BLD_REFINERY2},
		NeedsEmptyRowBelowWhenConstructing: true,
		Cost:                               1000,
		GivesTechLevel:                     3,
		BuildTime:                          12,
		Builds:                             nil,
		ConsumesEnergy:                     15,
		Produces:                           []int{UNT_TANK2, UNT_JUGGERNAUT, UNT_MCV2, UNT_QUAD, UNT_MSLTANK, UNT_AATANK, UNT_COMBAT_HARVESTER},
		HotkeyToBuild:                      "F",
	},
	BLD_AIRFACTORY2: {
		SpriteCode:        "airfactory",
		MaxHitpoints:      750,
		W:                 2,
		H:                 3,
		DisplayedName:     "Commonwealth Avionics Facility",
		Cost:              1000,
		RequiresTechLevel: 3,
		GivesTechLevel:    4,
		BuildTime:         10,
		Produces:          []int{AIR_TRANSPORT2, AIR_GUNSHIP, AIR_FIGHTER},
		ConsumesEnergy:    15,
		// produces:       []int{UNT_TANK, UNT_QUAD, UNT_MSLTANK, UNT_HARVESTER},
		HotkeyToBuild: "A",
	},
}
