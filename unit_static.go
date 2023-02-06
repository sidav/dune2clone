package main

const (
	UNT_INFANTRY = iota
	UNT_RECONINFANTRY
	UNT_ROCKETINFANTRY
	UNT_HEAVYINFANTRY
	UNT_SNIPERINFANTRY
	UNT_TANK1
	UNT_TANK2
	UNT_DEVASTATOR
	EUNT_DEVASTATOR
	UNT_JUGGERNAUT
	UNT_MCV1
	UNT_MCV2
	UNT_QUAD
	UNT_MSLTANK
	UNT_AATANK
	UNT_COMBAT_HARVESTER
	UNT_FAST_HARVESTER

	// aircrafts
	AIR_TRANSPORT1
	AIR_TRANSPORT2
	AIR_GUNSHIP
	AIR_FIGHTER
	AIR_FORTRESS
)

type unitStatic struct {
	DisplayedName     string `json:"displayed_name"`
	ChassisSpriteCode string `json:"chassis_sprite_code,omitempty"`
	HasEliteVersion   bool   `json:"has_elite_version,omitempty"`
	EliteVersionCode  int    `json:"elite_version_code,omitempty"`

	TurretsData []*TurretStatic `json:"turrets_data,omitempty"`

	MaxHitpoints int       `json:"max_hitpoints,omitempty"`
	HpRegen      int       `json:"hp_regen,omitempty"`
	ArmorType    armorCode `json:"armor_type,omitempty"`
	VisionRange  int       `json:"vision_range,omitempty"`

	MovementSpeed        float64 `json:"movement_speed,omitempty"`
	ChassisRotationSpeed int     `json:"chassis_rotation_speed,omitempty"`
	MaxSquadSize         int     `json:"max_squad_size,omitempty"`

	MaxCargoAmount int `json:"max_cargo_amount,omitempty"` // for harvesters

	DefaultOrderOnCreation orderCode `json:"default_order_on_creation,omitempty"`

	CanBeDeployed    bool         `json:"can_be_deployed,omitempty"`
	DeploysInto      buildingCode `json:"deploys_into,omitempty"`
	RequiresBuilding buildingCode `json:"requires_building,omitempty"`
	IsAircraft       bool         `json:"is_aircraft,omitempty"`
	IsTransport      bool         `json:"is_transport,omitempty"`

	Cost          int    `json:"cost,omitempty"`
	BuildTime     int    `json:"build_time,omitempty"` // seconds
	HotkeyToBuild string `json:"hotkey_to_build,omitempty"`
}
