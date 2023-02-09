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
