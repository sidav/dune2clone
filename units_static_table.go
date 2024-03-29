package main

var sTableUnits = map[int]*unitStatic{
	UNT_INFANTRY: {
		DisplayedName:     "Infantry squad",
		ChassisSpriteCode: "infantry",
		MaxHitpoints:      50,
		ArmorType:         ARMORTYPE_INFANTRY,
		VisionRange:       4,
		MovementSpeed:     0.1,
		TurretsData: []*TurretStatic{
			{
				SpriteCode:          "",
				AttacksLand:         true,
				RotateSpeed:         0,
				FireRange:           4,
				FireSpreadDegrees:   7,
				ShotRangeSpread:     0.5,
				CooldownAfterVolley: 45,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "bullets",
					HitDamage:                 4,
					DamageType:                DAMAGETYPE_ANTI_INFANTRY,
					Size:                      0.2,
					Speed:                     0.7,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_SMALL_EXPLOSION,
				},
			},
		},
		MaxSquadSize:         5,
		ChassisRotationSpeed: 90,
		Cost:                 250,
		BuildTime:            7,
		HotkeyToBuild:        "I",
	},
	UNT_RECONINFANTRY: {
		DisplayedName:     "Recon trike",
		ChassisSpriteCode: "infantryrecon",
		MaxHitpoints:      35,
		HpRegen:           1,
		ArmorType:         ARMORTYPE_INFANTRY,
		VisionRange:       6,
		MovementSpeed:     0.16,
		TurretsData: []*TurretStatic{
			{
				SpriteCode:          "",
				AttacksLand:         true,
				RotateSpeed:         0,
				FireRange:           4,
				FireSpreadDegrees:   7,
				ShotRangeSpread:     0.5,
				CooldownAfterVolley: 45,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "bullets",
					HitDamage:                 4,
					DamageType:                DAMAGETYPE_ANTI_INFANTRY,
					Size:                      0.2,
					Speed:                     0.7,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_SMALL_EXPLOSION,
				},
			},
		},
		MaxSquadSize:         1,
		ChassisRotationSpeed: 25,
		Cost:                 350,
		BuildTime:            12,
		HotkeyToBuild:        "C",
	},
	UNT_ROCKETINFANTRY: {
		DisplayedName:     "Rocketmen squad",
		ChassisSpriteCode: "infantryrocket",
		MaxHitpoints:      70,
		ArmorType:         ARMORTYPE_INFANTRY,
		VisionRange:       4,
		MovementSpeed:     0.075,
		TurretsData: []*TurretStatic{
			{
				SpriteCode:          "",
				AttacksLand:         true,
				AttacksAir:          true,
				RotateSpeed:         0,
				FireRange:           5,
				FireSpreadDegrees:   45,
				ShotRangeSpread:     0.65,
				CooldownAfterVolley: 105,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "aamissile",
					HitDamage:                 5,
					SplashDamage:              3,
					SplashRadius:              0.25,
					DamageType:                DAMAGETYPE_HEAVY,
					Size:                      0.3,
					Speed:                     0.45,
					RotationSpeed:             35,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_SMALL_EXPLOSION,
				},
			},
		},
		MaxSquadSize:         4,
		ChassisRotationSpeed: 90,
		Cost:                 450,
		BuildTime:            15,
		HotkeyToBuild:        "R",
	},
	UNT_HEAVYINFANTRY: {
		DisplayedName:     "Heavy infantry squad",
		ChassisSpriteCode: "infantryheavy",
		MaxHitpoints:      450,
		ArmorType:         ARMORTYPE_INFANTRY,
		VisionRange:       4,
		MovementSpeed:     0.06,
		TurretsData: []*TurretStatic{
			{
				SpriteCode:          "",
				AttacksLand:         true,
				RotateSpeed:         0,
				FireRange:           5,
				FireSpreadDegrees:   8,
				ShotRangeSpread:     0.45,
				CooldownAfterVolley: 40,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "omni",
					HitDamage:                 9,
					DamageType:                DAMAGETYPE_OMNI,
					Size:                      0.3,
					Speed:                     0.45,
					RotationSpeed:             0,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_SMALL_EXPLOSION,
				},
			},
		},
		MaxSquadSize:         3,
		ChassisRotationSpeed: 90,
		Cost:                 750,
		BuildTime:            25,
		HotkeyToBuild:        "H",
		RequiresBuilding:     BLD_AIRFACTORY1,
	},
	UNT_SNIPERINFANTRY: {
		DisplayedName:     "Sniper squad",
		ChassisSpriteCode: "infantrysniper",
		MaxHitpoints:      55,
		ArmorType:         ARMORTYPE_INFANTRY,
		VisionRange:       6,
		MovementSpeed:     0.03,
		TurretsData: []*TurretStatic{
			{
				SpriteCode:          "",
				AttacksLand:         true,
				RotateSpeed:         0,
				FireRange:           8,
				FireSpreadDegrees:   2,
				ShotRangeSpread:     0.45,
				CooldownAfterVolley: 40,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "bullets",
					HitDamage:                 40,
					DamageType:                DAMAGETYPE_ANTI_INFANTRY,
					Size:                      0.3,
					Speed:                     0.85,
					RotationSpeed:             0,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_SMALL_EXPLOSION,
				},
			},
		},
		MaxSquadSize:         3,
		ChassisRotationSpeed: 90,
		Cost:                 750,
		BuildTime:            25,
		HotkeyToBuild:        "S",
		RequiresBuilding:     BLD_AIRFACTORY2,
	},
	UNT_QUAD: {
		DisplayedName:     "Quad",
		ChassisSpriteCode: "quad",
		MaxHitpoints:      100,
		ArmorType:         ARMORTYPE_HEAVY,
		VisionRange:       5,
		MovementSpeed:     0.15,
		TurretsData: []*TurretStatic{
			{
				SpriteCode:          "",
				AttacksLand:         true,
				RotateSpeed:         0,
				FireRange:           4,
				FireSpreadDegrees:   6,
				ShotRangeSpread:     0.3,
				CooldownAfterVolley: 5,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "bullets",
					HitDamage:                 4,
					DamageType:                DAMAGETYPE_ANTI_INFANTRY,
					Size:                      0.2,
					Speed:                     0.7,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_SMALL_EXPLOSION,
				},
			},
		},
		ChassisRotationSpeed: 7,
		Cost:                 350,
		BuildTime:            15,
		HotkeyToBuild:        "Q",
	},
	UNT_TANK1: {
		DisplayedName:     "Medium tank",
		ChassisSpriteCode: "tank",
		MovementSpeed:     0.09,
		VisionRange:       5,
		MaxHitpoints:      125,
		ArmorType:         ARMORTYPE_HEAVY,
		TurretsData: []*TurretStatic{
			{
				SpriteCode:          "tank",
				AttacksLand:         true,
				RotateSpeed:         7,
				FireRange:           5,
				FireSpreadDegrees:   6,
				ShotRangeSpread:     0.7,
				CooldownAfterVolley: 55,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "shell",
					DamageType:                DAMAGETYPE_HEAVY,
					SplashRadius:              0.25,
					HitDamage:                 20,
					SplashDamage:              10,
					Size:                      0.3,
					Speed:                     0.7,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_REGULAR_EXPLOSION,
				},
			},
		},
		ChassisRotationSpeed: 5,
		Cost:                 500,
		BuildTime:            20,
		HotkeyToBuild:        "T",
	},
	UNT_JUGGERNAUT: {
		DisplayedName:     "Juggernaut",
		ChassisSpriteCode: "juggernaut",
		MovementSpeed:     0.02,
		VisionRange:       5,
		MaxHitpoints:      1000,
		ArmorType:         ARMORTYPE_HEAVY,
		TurretsData: []*TurretStatic{
			{
				SpriteCode:          "juggernautmain",
				AttacksLand:         true,
				RotateSpeed:         5,
				FireRange:           5,
				FireSpreadDegrees:   8,
				ShotRangeSpread:     0.7,
				CooldownAfterVolley: 250,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "shell",
					DamageType:                DAMAGETYPE_HEAVY,
					SplashRadius:              1.2,
					HitDamage:                 35,
					SplashDamage:              10,
					Size:                      0.3,
					Speed:                     0.65,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_BIGGER_EXPLOSION,
				},
			},
			{
				SpriteCode:          "juggernautsec",
				AttacksLand:         true,
				AttacksAir:          false,
				RotateSpeed:         15,
				FireRange:           5,
				FireSpreadDegrees:   12,
				ShotRangeSpread:     0.7,
				CooldownAfterVolley: 45,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "aamissile",
					DamageType:                DAMAGETYPE_ANTI_INFANTRY,
					SplashRadius:              0.5,
					RotationSpeed:             3,
					SplashDamage:              10,
					Size:                      0.3,
					Speed:                     0.7,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_SMALL_EXPLOSION,
				},
			},
			{
				SpriteCode:          "juggernautsec",
				AttacksLand:         false,
				AttacksAir:          true,
				RotateSpeed:         15,
				FireRange:           8,
				FireSpreadDegrees:   12,
				ShotRangeSpread:     0.7,
				CooldownAfterVolley: 45,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "aamissile",
					DamageType:                DAMAGETYPE_ANTI_INFANTRY,
					SplashRadius:              0.5,
					RotationSpeed:             3,
					SplashDamage:              10,
					Size:                      0.3,
					Speed:                     1.5,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_SMALL_EXPLOSION,
				},
			},
		},
		ChassisRotationSpeed: 1,
		Cost:                 2500,
		BuildTime:            60,
		RequiresBuilding:     BLD_AIRFACTORY2,
		HotkeyToBuild:        "J",
	},
	UNT_TANK2: {
		DisplayedName:     "Anjaopterix tank",
		ChassisSpriteCode: "tank2",
		MovementSpeed:     0.1,
		VisionRange:       5,
		MaxHitpoints:      120,
		ArmorType:         ARMORTYPE_HEAVY,
		TurretsData: []*TurretStatic{
			{
				SpriteCode:          "tank2",
				TurretCenterX:       -0.14,
				AttacksLand:         true,
				RotateSpeed:         7,
				FireRange:           5,
				FireSpreadDegrees:   7,
				ShotRangeSpread:     0.7,
				CooldownAfterVolley: 55,
				MaxShotsInVolley:    3,
				CooldownPerShot:     10,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "shell",
					Size:                      0.3,
					SplashRadius:              0.25,
					HitDamage:                 15,
					SplashDamage:              10,
					DamageType:                DAMAGETYPE_HEAVY,
					Speed:                     0.7,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_REGULAR_EXPLOSION,
				},
			},
		},
		ChassisRotationSpeed: 5,
		Cost:                 450,
		BuildTime:            20,
		HotkeyToBuild:        "T",
	},
	UNT_DEVASTATOR: {
		DisplayedName:     "Devastator",
		ChassisSpriteCode: "devastator",
		HasEliteVersion:   true,
		EliteVersionCode:  EUNT_DEVASTATOR,
		MovementSpeed:     0.04,
		VisionRange:       6,
		MaxHitpoints:      500,
		ArmorType:         ARMORTYPE_HEAVY,
		TurretsData: []*TurretStatic{
			{
				SpriteCode:          "devastator",
				AttacksLand:         true,
				RotateSpeed:         0,
				FireRange:           6,
				FireSpreadDegrees:   7,
				ShotRangeSpread:     0.5,
				CooldownAfterVolley: 85,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "shell",
					Size:                      0.3,
					SplashRadius:              0.5,
					HitDamage:                 15,
					SplashDamage:              10,
					DamageType:                DAMAGETYPE_HEAVY,
					Speed:                     0.7,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_REGULAR_EXPLOSION,
				},
			},
			{
				SpriteCode:          "",
				AttacksLand:         true,
				RotateSpeed:         0,
				FireRange:           6,
				FireSpreadDegrees:   7,
				ShotRangeSpread:     0.5,
				CooldownAfterVolley: 85,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "shell",
					Size:                      0.3,
					SplashRadius:              0.5,
					HitDamage:                 15,
					SplashDamage:              10,
					DamageType:                DAMAGETYPE_HEAVY,
					Speed:                     0.7,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_REGULAR_EXPLOSION,
				},
			},
		},
		ChassisRotationSpeed: 5,
		RequiresBuilding:     BLD_AIRFACTORY1,
		Cost:                 1500,
		BuildTime:            35,
		HotkeyToBuild:        "D",
	},
	EUNT_DEVASTATOR: {
		DisplayedName:     "Elite Devastator",
		ChassisSpriteCode: "devastator",
		MovementSpeed:     0.04,
		VisionRange:       6,
		MaxHitpoints:      550,
		ArmorType:         ARMORTYPE_HEAVY,
		TurretsData: []*TurretStatic{
			{
				SpriteCode:          "devastator",
				AttacksLand:         true,
				RotateSpeed:         5,
				FireRange:           7,
				FireSpreadDegrees:   6,
				ShotRangeSpread:     0.45,
				MaxShotsInVolley:    3,
				CooldownPerShot:     15,
				CooldownAfterVolley: 120,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "shell",
					Size:                      0.3,
					SplashRadius:              0.75,
					HitDamage:                 15,
					SplashDamage:              10,
					DamageType:                DAMAGETYPE_HEAVY,
					Speed:                     0.7,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_BIGGER_EXPLOSION,
				},
			},
		},
		ChassisRotationSpeed: 5,
		Cost:                 1500,
		BuildTime:            45,
		HotkeyToBuild:        "D",
	},
	UNT_MCV1: {
		DisplayedName:        "BetaCorp MCV",
		ChassisSpriteCode:    "mcv",
		MovementSpeed:        0.035,
		VisionRange:          4,
		MaxHitpoints:         300,
		ArmorType:            ARMORTYPE_HEAVY,
		CanBeDeployed:        true,
		DeploysInto:          BLD_CONYARD1,
		ChassisRotationSpeed: 4,
		Cost:                 750,
		BuildTime:            30,
		HotkeyToBuild:        "V",
	},
	UNT_MCV2: {
		DisplayedName:        "Commonwealth MCV",
		ChassisSpriteCode:    "mcv",
		MovementSpeed:        0.035,
		VisionRange:          4,
		MaxHitpoints:         300,
		ArmorType:            ARMORTYPE_HEAVY,
		CanBeDeployed:        true,
		DeploysInto:          BLD_CONYARD2,
		ChassisRotationSpeed: 4,
		Cost:                 750,
		BuildTime:            30,
		HotkeyToBuild:        "V",
	},
	UNT_MSLTANK1: {
		DisplayedName:     "Missile tank",
		ChassisSpriteCode: "tank",
		MovementSpeed:     0.05,
		VisionRange:       3,
		MaxHitpoints:      65,
		ArmorType:         ARMORTYPE_HEAVY,
		TurretsData: []*TurretStatic{
			{
				SpriteCode:          "msltank",
				AttacksLand:         true,
				RotateSpeed:         15,
				FireRange:           18,
				FireSpreadDegrees:   25,
				ShotRangeSpread:     0.85,
				CooldownAfterVolley: 200,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "missile",
					Size:                      0.3,
					Speed:                     0.3,
					SplashRadius:              1.5,
					HitDamage:                 25,
					SplashDamage:              25,
					DamageType:                DAMAGETYPE_ANTI_BUILDING,
					RotationSpeed:             1,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_BIGGER_EXPLOSION,
				},
			},
		},
		ChassisRotationSpeed: 8,
		Cost:                 1150,
		BuildTime:            35,
		HotkeyToBuild:        "M",
	},
	UNT_MSLTANK2: {
		DisplayedName:     "Commonwealth MRLS",
		ChassisSpriteCode: "tank2",
		MovementSpeed:     0.05,
		VisionRange:       3,
		MaxHitpoints:      65,
		ArmorType:         ARMORTYPE_HEAVY,
		TurretsData: []*TurretStatic{
			{
				TurretCenterX:       -0.14,
				SpriteCode:          "msltank",
				AttacksLand:         true,
				RotateSpeed:         15,
				FireRange:           18,
				FireSpreadDegrees:   25,
				ShotRangeSpread:     1.85,
				MaxShotsInVolley:    12,
				CooldownPerShot:     8,
				CooldownAfterVolley: 300,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "aamissile",
					Size:                      0.3,
					Speed:                     0.4,
					SplashRadius:              0.8,
					HitDamage:                 10,
					SplashDamage:              10,
					DamageType:                DAMAGETYPE_ANTI_BUILDING,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_REGULAR_EXPLOSION,
				},
			},
		},
		ChassisRotationSpeed: 8,
		Cost:                 1150,
		BuildTime:            35,
		HotkeyToBuild:        "M",
	},
	UNT_AATANK: {
		DisplayedName:     "AA tank",
		ChassisSpriteCode: "tank",
		MovementSpeed:     0.05,
		VisionRange:       3,
		MaxHitpoints:      75,
		ArmorType:         ARMORTYPE_HEAVY,
		TurretsData: []*TurretStatic{
			{
				SpriteCode:          "aamsltank",
				AttacksAir:          true,
				RotateSpeed:         15,
				FireRange:           8,
				FireSpreadDegrees:   35,
				ShotRangeSpread:     0.7,
				CooldownAfterVolley: 45,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "aamissile",
					Size:                      0.3,
					Speed:                     0.65,
					SplashRadius:              0.5,
					HitDamage:                 40,
					SplashDamage:              15,
					DamageType:                DAMAGETYPE_HEAVY,
					RotationSpeed:             45,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_REGULAR_EXPLOSION,
				},
			},
		},
		ChassisRotationSpeed: 8,
		Cost:                 1150,
		BuildTime:            25,
		HotkeyToBuild:        "A",
	},
	UNT_COMBAT_HARVESTER: {
		DisplayedName:          "Combat Harvester",
		ChassisSpriteCode:      "combatharvester",
		DefaultOrderOnCreation: ORDER_HARVEST,
		MaxCargoAmount:         700,
		MovementSpeed:          0.062,
		VisionRange:            3,
		MaxHitpoints:           275,
		ArmorType:              ARMORTYPE_HEAVY,
		TurretsData: []*TurretStatic{
			{
				SpriteCode:          "combatharvester",
				AttacksLand:         true,
				RotateSpeed:         25,
				FireRange:           5,
				FireSpreadDegrees:   11,
				ShotRangeSpread:     0.4,
				CooldownAfterVolley: 14,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "bullets",
					HitDamage:                 5,
					DamageType:                DAMAGETYPE_ANTI_INFANTRY,
					Size:                      0.2,
					Speed:                     0.7,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_SMALL_EXPLOSION,
				},
			},
		},
		ChassisRotationSpeed: 4,
		Cost:                 1600,
		BuildTime:            20,
		HotkeyToBuild:        "H",
	},
	UNT_FAST_HARVESTER: {
		DisplayedName:          "Patented Harvester",
		ChassisSpriteCode:      "fastharvester",
		DefaultOrderOnCreation: ORDER_HARVEST,
		MaxCargoAmount:         500,
		MovementSpeed:          0.072,
		VisionRange:            3,
		MaxHitpoints:           200,
		HpRegen:                1,
		ArmorType:              ARMORTYPE_HEAVY,
		ChassisRotationSpeed:   7,
		Cost:                   1600,
		BuildTime:              20,
		HotkeyToBuild:          "H",
	},
	// aircrafts
	AIR_TRANSPORT1: {
		DisplayedName:        "Carrier aircraft",
		ChassisSpriteCode:    "airtransport",
		MaxHitpoints:         100,
		ArmorType:            ARMORTYPE_HEAVY,
		MovementSpeed:        0.205,
		VisionRange:          3,
		HpRegen:              1,
		ChassisRotationSpeed: 6,
		Cost:                 650,
		BuildTime:            25,
		HotkeyToBuild:        "C",
		IsAircraft:           true,
		IsTransport:          true,
	},
	AIR_TRANSPORT2: {
		DisplayedName:        "Combat carrier",
		ChassisSpriteCode:    "airtransport2",
		MaxHitpoints:         120,
		ArmorType:            ARMORTYPE_HEAVY,
		MovementSpeed:        0.2,
		VisionRange:          3,
		ChassisRotationSpeed: 5,
		Cost:                 500,
		BuildTime:            25,
		HotkeyToBuild:        "C",
		IsAircraft:           true,
		IsTransport:          true,
		TurretsData: []*TurretStatic{
			{
				SpriteCode:          "",
				AttacksLand:         true,
				AttacksAir:          true,
				RotateSpeed:         25,
				FireRange:           5,
				FireSpreadDegrees:   45,
				ShotRangeSpread:     0.4,
				CooldownAfterVolley: 16,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "aamissile",
					HitDamage:                 5,
					DamageType:                DAMAGETYPE_ANTI_INFANTRY,
					Size:                      0.2,
					Speed:                     0.5,
					RotationSpeed:             5,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_SMALL_EXPLOSION,
				},
			},
		},
	},
	AIR_GUNSHIP: {
		DisplayedName:     "Gunship",
		ChassisSpriteCode: "air_gunship",
		MaxHitpoints:      120,
		ArmorType:         ARMORTYPE_HEAVY,
		MovementSpeed:     0.20,
		VisionRange:       5,
		TurretsData: []*TurretStatic{
			{
				SpriteCode:          "",
				AttacksLand:         true,
				RotateSpeed:         180,
				FireRange:           4,
				FireSpreadDegrees:   12,
				ShotRangeSpread:     1.0,
				CooldownAfterVolley: 35,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "shell",
					Size:                      0.3,
					DamageType:                DAMAGETYPE_HEAVY,
					HitDamage:                 25,
					SplashDamage:              15,
					SplashRadius:              0.5,
					Speed:                     0.7,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_REGULAR_EXPLOSION,
				},
			},
		},
		ChassisRotationSpeed: 2,
		Cost:                 600,
		BuildTime:            30,
		HotkeyToBuild:        "G",
		IsAircraft:           true,
	},
	AIR_FIGHTER: {
		DisplayedName:     "Fighter",
		ChassisSpriteCode: "air_fighter",
		MaxHitpoints:      75,
		ArmorType:         ARMORTYPE_HEAVY,
		MovementSpeed:     0.25,
		VisionRange:       7,
		TurretsData: []*TurretStatic{
			{
				SpriteCode:          "",
				AttacksLand:         false,
				AttacksAir:          true,
				RotateSpeed:         0,
				FireRange:           6,
				FireSpreadDegrees:   15,
				ShotRangeSpread:     2.0,
				CooldownAfterVolley: 25,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "aamissile",
					Size:                      0.3,
					Speed:                     0.65,
					HitDamage:                 25,
					DamageType:                DAMAGETYPE_HEAVY,
					RotationSpeed:             45,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_REGULAR_EXPLOSION,
				},
			},
		},
		ChassisRotationSpeed: 3,
		Cost:                 700,
		BuildTime:            30,
		HotkeyToBuild:        "F",
		IsAircraft:           true,
	},
	AIR_FORTRESS: {
		DisplayedName:     "Big Sister",
		ChassisSpriteCode: "air_fortress",
		EliteVersionCode:  0,
		TurretsData: []*TurretStatic{
			{
				SpriteCode:          "",
				AttacksLand:         true,
				AttacksAir:          true,
				RotateSpeed:         45,
				FireRange:           7,
				FireSpreadDegrees:   35,
				ShotRangeSpread:     1.5,
				CooldownAfterVolley: 125,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "aamissile",
					Size:                      0.3,
					Speed:                     0.6,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_REGULAR_EXPLOSION,
					HitDamage:                 25,
					RotationSpeed:             1,
					SplashDamage:              15,
					SplashRadius:              0.5,
					DamageType:                DAMAGETYPE_HEAVY,
				},
			},
			{
				SpriteCode:          "",
				AttacksLand:         true,
				AttacksAir:          true,
				RotateSpeed:         45,
				FireRange:           5,
				FireSpreadDegrees:   25,
				ShotRangeSpread:     1.5,
				CooldownAfterVolley: 12,
				FiredProjectileData: &projectileStatic{
					SpriteCode:                "bullets",
					Size:                      0.3,
					Speed:                     0.7,
					CreatesEffectOnImpact:     true,
					EffectCreatedOnImpactCode: EFFECT_SMALL_EXPLOSION,
					HitDamage:                 4,
					RotationSpeed:             1,
					DamageType:                DAMAGETYPE_ANTI_INFANTRY,
				},
			},
		},
		MaxHitpoints:         750,
		HpRegen:              1,
		ArmorType:            ARMORTYPE_HEAVY,
		VisionRange:          12,
		MovementSpeed:        0.05,
		ChassisRotationSpeed: 3,
		CanBeDeployed:        false,
		IsAircraft:           true,
		IsTransport:          false,
		Cost:                 4000,
		BuildTime:            60,
		HotkeyToBuild:        "A",
	},
}
