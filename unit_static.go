package main

const (
	UNT_INFANTRY = iota
	UNT_RECONINFANTRY
	UNT_ROCKETINFANTRY
	UNT_HEAVYINFANTRY
	UNT_TANK
	UNT_TANK2
	UNT_DEVASTATOR
	UNT_MCV1
	UNT_MCV2
	UNT_QUAD
	UNT_MSLTANK
	UNT_AATANK
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
	armorType    armorCode
	visionRange  int

	movementSpeed        float64
	chassisRotationSpeed int
	maxSquadSize         int

	maxCargoAmount int // for harvesters

	defaultOrderOnCreation orderCode

	canBeDeployed bool
	deploysInto   buildingCode // building code
	isAircraft    bool
	isTransport   bool

	cost          int
	buildTime     int // seconds
	hotkeyToBuild string
}

type unitStaticTurretsData struct {
	turretCode                   int
	turretCenterX, turretCenterY float64
}

var sTableUnits = map[int]*unitStatic{
	UNT_INFANTRY: {
		displayedName:     "Infantry squad",
		chassisSpriteCode: "infantry",
		maxHitpoints:      50,
		armorType:         ARMORTYPE_INFANTRY,
		visionRange:       4,
		movementSpeed:     0.1,
		turretsData: []unitStaticTurretsData{
			{
				turretCode: TRT_INFANTRY,
			},
		},
		maxSquadSize:         5,
		chassisRotationSpeed: 90,
		cost:                 250,
		buildTime:            7,
		hotkeyToBuild:        "I",
	},
	UNT_RECONINFANTRY: {
		displayedName:     "Recon trike",
		chassisSpriteCode: "infantryrecon",
		maxHitpoints:      35,
		armorType:         ARMORTYPE_INFANTRY,
		visionRange:       6,
		movementSpeed:     0.16,
		turretsData: []unitStaticTurretsData{
			{
				turretCode: TRT_INFANTRY,
			},
		},
		maxSquadSize:         1,
		chassisRotationSpeed: 25,
		cost:                 350,
		buildTime:            7,
		hotkeyToBuild:        "C",
	},
	UNT_ROCKETINFANTRY: {
		displayedName:     "Rocketmen squad",
		chassisSpriteCode: "infantryrocket",
		maxHitpoints:      70,
		armorType:         ARMORTYPE_INFANTRY,
		visionRange:       4,
		movementSpeed:     0.075,
		turretsData: []unitStaticTurretsData{
			{
				turretCode: TRT_ROCKETINFANTRY,
			},
		},
		maxSquadSize:         4,
		chassisRotationSpeed: 90,
		cost:                 450,
		buildTime:            10,
		hotkeyToBuild:        "R",
	},
	UNT_HEAVYINFANTRY: {
		displayedName:     "Heavy infantry squad",
		chassisSpriteCode: "infantryheavy",
		maxHitpoints:      150,
		armorType:         ARMORTYPE_INFANTRY,
		visionRange:       4,
		movementSpeed:     0.06,
		turretsData: []unitStaticTurretsData{
			{
				turretCode: TRT_HEAVYINFANTRY,
			},
		},
		maxSquadSize:         3,
		chassisRotationSpeed: 90,
		cost:                 550,
		buildTime:            15,
		hotkeyToBuild:        "H",
	},
	UNT_QUAD: {
		displayedName:     "Quad",
		chassisSpriteCode: "quad",
		maxHitpoints:      100,
		armorType:         ARMORTYPE_HEAVY,
		visionRange:       5,
		movementSpeed:     0.15,
		turretsData: []unitStaticTurretsData{
			{
				turretCode: TRT_QUAD,
			},
		},
		chassisRotationSpeed: 7,
		cost:                 350,
		buildTime:            7,
		hotkeyToBuild:        "Q",
	},
	UNT_TANK: {
		displayedName:     "Super duper tank",
		chassisSpriteCode: "tank",
		movementSpeed:     0.1,
		visionRange:       4,
		maxHitpoints:      120,
		armorType:         ARMORTYPE_HEAVY,
		turretsData: []unitStaticTurretsData{
			{
				turretCode: TRT_TANK,
			},
		},
		chassisRotationSpeed: 5,
		cost:                 450,
		buildTime:            12,
		hotkeyToBuild:        "T",
	},
	UNT_TANK2: {
		displayedName:     "Anjaopterix tank",
		chassisSpriteCode: "tank2",
		movementSpeed:     0.1,
		visionRange:       4,
		maxHitpoints:      120,
		armorType:         ARMORTYPE_HEAVY,
		turretsData: []unitStaticTurretsData{
			{
				turretCode:    TRT_TANK2,
				turretCenterX: -0.14,
			},
		},
		chassisRotationSpeed: 5,
		cost:                 450,
		buildTime:            12,
		hotkeyToBuild:        "T",
	},
	UNT_DEVASTATOR: {
		displayedName:     "Devastator",
		chassisSpriteCode: "devastator",
		movementSpeed:     0.04,
		visionRange:       5,
		maxHitpoints:      500,
		armorType:         ARMORTYPE_HEAVY,
		turretsData: []unitStaticTurretsData{
			{
				turretCode:    TRT_DEVASTATOR,
			},
			{
				turretCode:    TRT_DEVASTATOR,
			},
		},
		chassisRotationSpeed: 5,
		cost:                 1500,
		buildTime:            30,
		hotkeyToBuild:        "D",
	},
	UNT_MCV1: {
		displayedName:        "BetaCorp MCV",
		chassisSpriteCode:    "placeholder",
		movementSpeed:        0.035,
		visionRange:          4,
		maxHitpoints:         300,
		armorType:            ARMORTYPE_HEAVY,
		canBeDeployed:        true,
		deploysInto:          BLD_CONYARD1,
		chassisRotationSpeed: 4,
		cost:                 750,
		buildTime:            15,
		hotkeyToBuild:        "V",
	},
	UNT_MCV2: {
		displayedName:        "Commonwealth MCV",
		chassisSpriteCode:    "placeholder",
		movementSpeed:        0.035,
		visionRange:          4,
		maxHitpoints:         300,
		armorType:            ARMORTYPE_HEAVY,
		canBeDeployed:        true,
		deploysInto:          BLD_CONYARD2,
		chassisRotationSpeed: 4,
		cost:                 750,
		buildTime:            15,
		hotkeyToBuild:        "V",
	},
	UNT_MSLTANK: {
		displayedName:     "Missile tank",
		chassisSpriteCode: "tank",
		movementSpeed:     0.05,
		visionRange:       3,
		maxHitpoints:      40,
		armorType:         ARMORTYPE_HEAVY,
		turretsData: []unitStaticTurretsData{
			{
				turretCode:    TRT_MSLTANK,
				turretCenterX: 0,
				turretCenterY: 0,
			},
		},
		chassisRotationSpeed: 8,
		cost:                 1150,
		buildTime:            25,
		hotkeyToBuild:        "M",
	},
	UNT_AATANK: {
		displayedName:     "AA tank",
		chassisSpriteCode: "quad",
		movementSpeed:     0.05,
		visionRange:       3,
		maxHitpoints:      40,
		armorType:         ARMORTYPE_HEAVY,
		turretsData: []unitStaticTurretsData{
			{
				turretCode:    TRT_AATANK,
				turretCenterX: 0,
				turretCenterY: 0,
			},
		},
		chassisRotationSpeed: 8,
		cost:                 1150,
		buildTime:            17,
		hotkeyToBuild:        "A",
	},
	UNT_HARVESTER: {
		displayedName:          "Harvester",
		chassisSpriteCode:      "harvester",
		defaultOrderOnCreation: ORDER_HARVEST,
		maxCargoAmount:         700,
		movementSpeed:          0.07,
		visionRange:            2,
		maxHitpoints:           275,
		armorType:              ARMORTYPE_HEAVY,
		turretsData: []unitStaticTurretsData{
			{
				turretCode:    TRT_HARVESTER,
				turretCenterX: 0,
				turretCenterY: 0,
			},
		},
		chassisRotationSpeed: 4,
		cost:                 1600,
		buildTime:            12,
		hotkeyToBuild:        "H",
	},
	// aircrafts
	AIR_TRANSPORT: {
		displayedName:        "Carrier aircraft",
		chassisSpriteCode:    "air_transport",
		maxHitpoints:         100,
		armorType:            ARMORTYPE_HEAVY,
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
		armorType:         ARMORTYPE_HEAVY,
		movementSpeed:     0.25,
		visionRange:       7,
		turretsData: []unitStaticTurretsData{
			{
				turretCode: TRT_AIR_GUNSHIP,
			},
		},
		chassisRotationSpeed: 3,
		cost:                 500,
		buildTime:            25,
		hotkeyToBuild:        "G",
		isAircraft:           true,
	},
}
