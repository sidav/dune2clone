package main

import (
	"dune2clone/astar"
	"dune2clone/map_generator"
	"fmt"
)

type startConditions struct {
	aiType              string
	factionName         string
	resourcesMultiplier float64
}

func (b *battlefield) initFromRandomMap(rm *map_generator.GeneratedMap, conds []*startConditions) {
	b.tiles = make([][]tile, len(rm.Tiles))
	for i := range b.tiles {
		b.tiles[i] = make([]tile, len(rm.Tiles[i]))
		for j := range b.tiles[i] {
			var currTileCode int
			switch rm.Tiles[i][j] {
			case map_generator.SAND:
				currTileCode = TILE_SAND
			case map_generator.BUILDABLE_TERRAIN:
				currTileCode = TILE_BUILDABLE
			case map_generator.ROCKS:
				currTileCode = TILE_ROCK
			case map_generator.POOR_RESOURCES:
				currTileCode = TILE_SAND
				b.tiles[i][j].resourcesAmount = rnd.RandInRange(config.Economy.ResourcesInTileMinGenerated, config.Economy.ResourceInTilePoorMax)
			case map_generator.MEDIUM_RESOURCES:
				currTileCode = TILE_SAND
				b.tiles[i][j].resourcesAmount = rnd.RandInRange(config.Economy.ResourceInTilePoorMax, config.Economy.ResourceInTileMediumMax)
			case map_generator.RICH_RESOURCES:
				currTileCode = TILE_SAND
				b.tiles[i][j].resourcesAmount = rnd.RandInRange(config.Economy.ResourceInTileMediumMax, config.Economy.ResourceInTileRichMax)
			case map_generator.RESOURCE_VEIN:
				currTileCode = TILE_SAND
				b.tiles[i][j].hasResourceVein = true
			default:
				panic("Unknown tile type!")
				currTileCode = TILE_SAND
			}
			b.tiles[i][j].code = currTileCode
		}
	}
	b.pathfinder = &astar.AStarPathfinder{
		DiagonalMoveAllowed:       true,
		ForceGetPath:              true,
		ForceIncludeFinish:        false,
		AutoAdjustDefaultMaxSteps: false,
		MapWidth:                  len(b.tiles),
		MapHeight:                 len(b.tiles[0]),
	}
	b.placeInitialStuff(rm.StartPoints, conds)
	b.finalizeTileVariants()
}

func (b *battlefield) finalizeTileVariants() {
	for i := range b.tiles {
		for j := range b.tiles[i] {
			//debugWritef("CODE: %d ", b.tiles[i][j].code)
			//debugWritef("CONTENTS: %+v\n", sTableTiles[b.tiles[i][j].code])
			b.tiles[i][j].spriteVariantIndex = rnd.Rand(len(sTableTiles[b.tiles[i][j].code].spriteCodes))
		}
	}
}

func (b *battlefield) placeInitialStuff(startPoints [][2]int, conditions []*startConditions) {
	for spNumber := range startPoints {
		b.factions = append(b.factions, createFaction(-1, 0, 10000))
		b.factions[spNumber].resourcesMultiplier = conditions[spNumber].resourcesMultiplier
		b.factions[spNumber].resetVisibilityMaps(len(b.tiles), len(b.tiles[0]))
		b.factions[spNumber].exploreAround(startPoints[spNumber][0], startPoints[spNumber][1], 2, 2, 3)
		untToCreate := UNT_MCV1
		switch conditions[spNumber].factionName {
		case "Commonwealth":
			untToCreate = UNT_MCV2
		case "Random":
			if rnd.OneChanceFrom(2) {
				untToCreate = UNT_MCV2
			}
		}
		b.addActor(createUnit(untToCreate, startPoints[spNumber][0], startPoints[spNumber][1], b.factions[spNumber]))

		if conditions[spNumber].aiType != "player" {
			b.ais = append(b.ais, createAi(b.factions[spNumber], fmt.Sprintf("Enemy %d", spNumber), conditions[spNumber].aiType))
			b.factions[spNumber].experienceMultiplier = 2
			b.factions[spNumber].storagesMultiplier = 1
			b.factions[spNumber].playerName = fmt.Sprintf("Player %d (AI %s)", spNumber+1, conditions[spNumber].aiType)
			// b.factions[spNumber].buildSpeedMultiplier = 10
			// b.factions[spNumber].visibilityCheat = true
			// b.factions[spNumber].explorationCheat = true
		} else {
			// player faction settings
			//b.factions[spNumber].buildSpeedMultiplier = 10
			//b.factions[spNumber].visibilityCheat = true
			//b.factions[spNumber].explorationCheat = true
			b.factions[spNumber].playerName = "Player"
		}
	}
	// randomize colors
	for _, f1 := range b.factions {
		colorIsUnique := false
		for !colorIsUnique {
			f1.colorNumber = rnd.Rand(len(factionColors))
			colorIsUnique = true
			for _, f2 := range b.factions {
				if f1 != f2 && f1.colorNumber == f2.colorNumber {
					colorIsUnique = false
					break
				}
			}
		}
	}

	// create all units for debugging
	//coord := 0
	//for k, _ := range sTableUnits {
	//	b.addActor(createUnit(k, startPoints[0][0]-10+coord, startPoints[0][1]+2, b.factions[0]))
	//	b.addActor(createUnit(k, startPoints[1][0]-10+coord, startPoints[1][1]+2, b.factions[1]))
	//	coord++
	//}
	// b.addActor(createUnit(UNT_JUGGERNAUT, startPoints[0][0]-1, startPoints[0][1]+2, b.factions[0]))
	//b.addActor(createUnit(UNT_TANK1, startPoints[0][0]-1, startPoints[0][1]+3, b.factions[1]))
}
