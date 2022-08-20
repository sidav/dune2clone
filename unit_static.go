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

	turretsData []unitStaticTurretsData

	maxHitpoints int
	visionRange  int

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

type unitStaticTurretsData struct {
	turretCode                   int
	turretCenterX, turretCenterY float64
}

var sTableUnits = map[int]*unitStatic{
	UNT_QUAD: {
		displayedName:     "Quad",
		chassisSpriteCode: "quad",
		maxHitpoints:      75,
		visionRange:       6,
		movementSpeed:     0.2,
		turretsData: []unitStaticTurretsData{
			{
				turretCode: TRT_QUAD,
			},
		},
		chassisRotationSpeed: 7,
		cost:                 350,
		buildTime:            3,
		hotkeyToBuild:        "Q",
	},
	UNT_TANK: {
		displayedName:     "Super duper tank",
		chassisSpriteCode: "tank",
		movementSpeed:     0.1,
		visionRange:       4,
		maxHitpoints:      120,
		turretsData: []unitStaticTurretsData{
			{
				turretCode: TRT_TANK,
			},
		},
		chassisRotationSpeed: 5,
		cost:                 450,
		buildTime:            7,
		hotkeyToBuild:        "T",
	},
	UNT_MSLTANK: {
		displayedName:     "Missile tank",
		chassisSpriteCode: "quad",
		movementSpeed:     0.05,
		visionRange:       3,
		maxHitpoints:      50,
		turretsData: []unitStaticTurretsData{
			{
				turretCode:    TRT_MSLTANK,
				turretCenterX: 0,
				turretCenterY: 0,
			},
		},
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
		visionRange:            2,
		maxHitpoints:           250,
		turretsData:            []unitStaticTurretsData{},
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
		visionRange:          1,
		turretsData:          []unitStaticTurretsData{},
		chassisRotationSpeed: 5,
		cost:                 500,
		buildTime:            20,
		hotkeyToBuild:        "C",
		isAircraft:           true,
		isTransport:          true,
	},
	AIR_GUNSHIP: {
		displayedName:     "Gunship",
		chassisSpriteCode: "air_gunship",
		maxHitpoints:      50,
		movementSpeed:     0.25,
		visionRange:       7,
		turretsData: []unitStaticTurretsData{
			{
				turretCode: TRT_AIR_GUNSHIP,
			},
		},
		chassisRotationSpeed: 3,
		cost:                 500,
		buildTime:            15,
		hotkeyToBuild:        "G",
		isAircraft:           true,
	},
}
