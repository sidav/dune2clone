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
	UNT_COMBAT_HARVESTER
	UNT_FAST_HARVESTER

	// aircrafts
	AIR_TRANSPORT
	AIR_GUNSHIP
	AIR_FIGHTER
)

type unitStatic struct {
	displayedName     string
	chassisSpriteCode string

	turretsData []*turretStatic

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

var sTableUnits = map[int]*unitStatic{
	UNT_INFANTRY: {
		displayedName:     "Infantry squad",
		chassisSpriteCode: "infantry",
		maxHitpoints:      50,
		armorType:         ARMORTYPE_INFANTRY,
		visionRange:       4,
		movementSpeed:     0.1,
		turretsData: []*turretStatic{
			{
				spriteCode:            "",
				firesProjectileOfCode: PRJ_BULLETS,
				attacksLand:           true,
				rotateSpeed:           0,
				fireRange:             4,
				fireSpreadDegrees:     7,
				shotRangeSpread:       0.5,
				attackCooldown:        45,
				projectileDamage:      4,
				projectileDamageType:  DAMAGETYPE_ANTI_INFANTRY,
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
		turretsData: []*turretStatic{
			{
				spriteCode:            "",
				firesProjectileOfCode: PRJ_BULLETS,
				attacksLand:           true,
				rotateSpeed:           0,
				fireRange:             4,
				fireSpreadDegrees:     7,
				shotRangeSpread:       0.5,
				attackCooldown:        45,
				projectileDamage:      4,
				projectileDamageType:  DAMAGETYPE_ANTI_INFANTRY,
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
		turretsData: []*turretStatic{
			{
				spriteCode:            "",
				firesProjectileOfCode: PRJ_INFANTRY_MISSILE,
				attacksLand:           true,
				attacksAir:            true,
				rotateSpeed:           0,
				fireRange:             5,
				fireSpreadDegrees:     45,
				shotRangeSpread:       0.5,
				attackCooldown:        105,
				projectileDamage:      8,
				projectileDamageType:  DAMAGETYPE_HEAVY,
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
		turretsData: []*turretStatic{
			{
				spriteCode:            "",
				firesProjectileOfCode: PRJ_OMNI,
				attacksLand:           true,
				rotateSpeed:           0,
				fireRange:             5,
				fireSpreadDegrees:     8,
				shotRangeSpread:       0.45,
				attackCooldown:        40,
				projectileDamage:      6,
				projectileDamageType:  DAMAGETYPE_OMNI,
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
		turretsData: []*turretStatic{
			{
				spriteCode:            "",
				firesProjectileOfCode: PRJ_BULLETS,
				attacksLand:           true,
				rotateSpeed:           0,
				fireRange:             4,
				fireSpreadDegrees:     6,
				shotRangeSpread:       0.3,
				attackCooldown:        5,
				projectileDamage:      4,
				projectileDamageType:  DAMAGETYPE_ANTI_INFANTRY,
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
		turretsData: []*turretStatic{
			{
				spriteCode:            "tank",
				firesProjectileOfCode: PRJ_SHELL,
				attacksLand:           true,
				rotateSpeed:           7,
				fireRange:             5,
				fireSpreadDegrees:     7,
				shotRangeSpread:       0.7,
				attackCooldown:        45,
				projectileDamage:      30,
				projectileDamageType:  DAMAGETYPE_HEAVY,
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
		turretsData: []*turretStatic{
			{
				spriteCode:            "tank2",
				firesProjectileOfCode: PRJ_SHELL,
				turretCenterX:         -0.14,
				attacksLand:           true,
				rotateSpeed:           7,
				fireRange:             5,
				fireSpreadDegrees:     7,
				shotRangeSpread:       0.7,
				attackCooldown:        45,
				projectileDamage:      30,
				projectileDamageType:  DAMAGETYPE_HEAVY,
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
		turretsData: []*turretStatic{
			{
				spriteCode:            "devastator",
				firesProjectileOfCode: PRJ_SHELL,
				attacksLand:           true,
				rotateSpeed:           0,
				fireRange:             6,
				fireSpreadDegrees:     7,
				shotRangeSpread:       0.5,
				attackCooldown:        75,
				projectileDamage:      47,
				projectileDamageType:  DAMAGETYPE_HEAVY,
			},
			{
				spriteCode:            "",
				firesProjectileOfCode: PRJ_SHELL,
				attacksLand:           true,
				rotateSpeed:           0,
				fireRange:             6,
				fireSpreadDegrees:     7,
				shotRangeSpread:       0.5,
				attackCooldown:        75,
				projectileDamage:      47,
				projectileDamageType:  DAMAGETYPE_HEAVY,
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
		turretsData: []*turretStatic{
			{
				spriteCode:            "msltank",
				firesProjectileOfCode: PRJ_MISSILE,
				attacksLand:           true,
				rotateSpeed:           15,
				fireRange:             10,
				fireSpreadDegrees:     25,
				shotRangeSpread:       0.7,
				attackCooldown:        150,
				projectileDamage:      50,
				projectileDamageType:  DAMAGETYPE_ANTI_BUILDING,
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
		turretsData: []*turretStatic{
			{
				spriteCode:            "aamsltank",
				firesProjectileOfCode: PRJ_AA_MISSILE,
				attacksAir:            true,
				rotateSpeed:           15,
				fireRange:             10,
				fireSpreadDegrees:     35,
				shotRangeSpread:       0.7,
				attackCooldown:        75,
				projectileDamage:      45,
				projectileDamageType:  DAMAGETYPE_HEAVY,
			},
		},
		chassisRotationSpeed: 8,
		cost:                 1150,
		buildTime:            17,
		hotkeyToBuild:        "A",
	},
	UNT_COMBAT_HARVESTER: {
		displayedName:          "Combat Harvester",
		chassisSpriteCode:      "combatharvester",
		defaultOrderOnCreation: ORDER_HARVEST,
		maxCargoAmount:         700,
		movementSpeed:          0.07,
		visionRange:            2,
		maxHitpoints:           275,
		armorType:              ARMORTYPE_HEAVY,
		turretsData: []*turretStatic{
			{
				spriteCode:            "",
				firesProjectileOfCode: PRJ_BULLETS,
				attacksLand:           true,
				rotateSpeed:           90,
				fireRange:             5,
				fireSpreadDegrees:     11,
				shotRangeSpread:       0.4,
				attackCooldown:        15,
				projectileDamage:      3,
				projectileDamageType:  DAMAGETYPE_ANTI_INFANTRY,
			},
		},
		chassisRotationSpeed: 4,
		cost:                 1600,
		buildTime:            12,
		hotkeyToBuild:        "H",
	},
	UNT_FAST_HARVESTER: {
		displayedName:          "Patented Harvester",
		chassisSpriteCode:      "fastharvester",
		defaultOrderOnCreation: ORDER_HARVEST,
		maxCargoAmount:         500,
		movementSpeed:          0.075,
		visionRange:            2,
		maxHitpoints:           200,
		armorType:              ARMORTYPE_HEAVY,
		chassisRotationSpeed:   7,
		cost:                   1600,
		buildTime:              12,
		hotkeyToBuild:          "H",
	},
	// aircrafts
	AIR_TRANSPORT: {
		displayedName:     "Carrier aircraft",
		chassisSpriteCode: "air_transport",
		maxHitpoints:      100,
		armorType:         ARMORTYPE_HEAVY,
		movementSpeed:     0.2,
		visionRange:       1,
		turretsData: []*turretStatic{

		},
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
		maxHitpoints:      100,
		armorType:         ARMORTYPE_HEAVY,
		movementSpeed:     0.20,
		visionRange:       7,
		turretsData: []*turretStatic{
			{
				spriteCode:            "",
				firesProjectileOfCode: PRJ_SHELL,
				attacksLand:           true,
				rotateSpeed:           180,
				fireRange:             6,
				fireSpreadDegrees:     12,
				shotRangeSpread:       1.0,
				attackCooldown:        15,
				projectileDamage:      12,
				projectileDamageType:  DAMAGETYPE_HEAVY,
			},
		},
		chassisRotationSpeed: 2,
		cost:                 600,
		buildTime:            25,
		hotkeyToBuild:        "G",
		isAircraft:           true,
	},
	AIR_FIGHTER: {
		displayedName:     "Fighter",
		chassisSpriteCode: "air_fighter",
		maxHitpoints:      75,
		armorType:         ARMORTYPE_HEAVY,
		movementSpeed:     0.25,
		visionRange:       7,
		turretsData: []*turretStatic{
			{
				spriteCode:            "",
				firesProjectileOfCode: PRJ_AA_MISSILE,
				attacksLand:           false,
				attacksAir:            true,
				rotateSpeed:           0,
				fireRange:             6,
				fireSpreadDegrees:     15,
				shotRangeSpread:       2.0,
				attackCooldown:        25,
				projectileDamage:      40,
				projectileDamageType:  DAMAGETYPE_HEAVY,
			},
		},
		chassisRotationSpeed: 3,
		cost:                 700,
		buildTime:            25,
		hotkeyToBuild:        "F",
		isAircraft:           true,
	},
}
