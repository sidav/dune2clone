package main

const (
	UNT_TANK = iota
	UNT_QUAD
	UNT_MSLTANK
	UNT_HARVESTER

	// aircrafts
	AIR_TRANSPORT
	AIR_GUNSHIP
)

type unitStatic struct {
	displayedName     string
	chassisSpriteCode string

	turretCode int

	maxHitpoints int

	movementSpeed        float64
	chassisRotationSpeed int

	maxCargoAmount int // for harvesters

	defaultOrderOnCreation orderCode

	isAircraft  bool
	isTransport bool

	cost          int
	buildTime     int // seconds
	hotkeyToBuild string
}

var sTableUnits = map[int]*unitStatic{
	UNT_QUAD: {
		displayedName:        "Quad",
		chassisSpriteCode:    "quad",
		maxHitpoints:         75,
		movementSpeed:        0.25,
		turretCode:           TRT_QUAD,
		chassisRotationSpeed: 7,
		cost:                 350,
		buildTime:            3,
		hotkeyToBuild:        "Q",
	},
	UNT_TANK: {
		displayedName:        "Super duper tank",
		chassisSpriteCode:    "tank",
		movementSpeed:        0.1,
		maxHitpoints:         120,
		turretCode:           TRT_TANK,
		chassisRotationSpeed: 5,
		cost:                 450,
		buildTime:            7,
		hotkeyToBuild:        "T",
	},
	UNT_MSLTANK: {
		displayedName:        "Missile tank",
		chassisSpriteCode:    "quad",
		movementSpeed:        0.05,
		maxHitpoints:         50,
		turretCode:           TRT_MSLTANK,
		chassisRotationSpeed: 8,
		cost:                 1150,
		buildTime:            12,
		hotkeyToBuild:        "M",
	},
	UNT_HARVESTER: {
		displayedName:          "Harvester",
		chassisSpriteCode:      "harvester",
		defaultOrderOnCreation: ORDER_HARVEST,
		maxCargoAmount:         700,
		movementSpeed:          0.07,
		maxHitpoints:           250,
		turretCode:             TRT_NONE,
		chassisRotationSpeed:   4,
		cost:                   1600,
		buildTime:              12,
		hotkeyToBuild:          "H",
	},
	// aircrafts
	AIR_TRANSPORT: {
		displayedName:        "Carrier aircraft",
		chassisSpriteCode:    "air_transport",
		maxHitpoints:         100,
		movementSpeed:        0.2,
		turretCode:           TRT_NONE,
		chassisRotationSpeed: 5,
		cost:                 500,
		buildTime:            20,
		hotkeyToBuild:        "C",
		isAircraft:           true,
		isTransport:          true,
	},
	AIR_GUNSHIP: {
		displayedName:        "Gunship",
		chassisSpriteCode:    "air_gunship",
		maxHitpoints:         50,
		movementSpeed:        0.25,
		turretCode:           TRT_AIR_GUNSHIP,
		chassisRotationSpeed: 3,
		cost:                 500,
		buildTime:            15,
		hotkeyToBuild:        "G",
		isAircraft:           true,
	},
}
