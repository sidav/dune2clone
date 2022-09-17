package main

const (
	UNT_INFANTRY = iota
	UNT_TANK
	UNT_TANK2
	UNT_MCV
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
	deploysInto   int // building code
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
		chassisRotationSpeed: 360,
		cost:                 250,
		buildTime:            7,
		hotkeyToBuild:        "I",
	},
	UNT_QUAD: {
		displayedName:     "Quad",
		chassisSpriteCode: "quad",
		maxHitpoints:      75,
		armorType:         ARMORTYPE_HEAVY,
		visionRange:       6,
		movementSpeed:     0.2,
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
	UNT_MCV: {
		displayedName:        "MCV",
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
		maxHitpoints:           200,
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
