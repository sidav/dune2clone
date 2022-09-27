package main

import (
	"dune2clone/astar"
	"dune2clone/map_generator"
	"fmt"
)

func (b *battlefield) initFromRandomMap(rm *map_generator.GeneratedMap) {
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
				b.tiles[i][j].resourcesAmount = rnd.RandInRange(RESOURCE_IN_TILE_MIN_GENERATED, RESOURCE_IN_TILE_POOR_MAX)
			case map_generator.MEDIUM_RESOURCES:
				currTileCode = TILE_SAND
				b.tiles[i][j].resourcesAmount = rnd.RandInRange(RESOURCE_IN_TILE_POOR_MAX, RESOURCE_IN_TILE_MEDIUM_MAX)
			case map_generator.RICH_RESOURCES:
				currTileCode = TILE_SAND
				b.tiles[i][j].resourcesAmount = rnd.RandInRange(RESOURCE_IN_TILE_MEDIUM_MAX, RESOURCE_IN_TILE_RICH_MAX)
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
	b.placeInitialStuff(rm.StartPoints)
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

func (b *battlefield) placeInitialStuff(startPoints [][2]int) {
	for spNumber := range startPoints {
		b.factions = append(b.factions, createFaction(spNumber, 0, 10000, 4, 1))
		b.factions[spNumber].resetVisibilityMaps(len(b.tiles), len(b.tiles[0]))
		b.factions[spNumber].exploreAround(startPoints[spNumber][0], startPoints[spNumber][1], 2, 2, 3)
		// TODO: faction selection
		if rnd.OneChanceFrom(2) {
			b.addActor(createUnit(UNT_MCV1, startPoints[spNumber][0], startPoints[spNumber][1], b.factions[spNumber]))
		} else {
			b.addActor(createUnit(UNT_MCV2, startPoints[spNumber][0], startPoints[spNumber][1], b.factions[spNumber]))
		}
		// b.addActor(createUnit(UNT_HARVESTER, startPoints[spNumber][0]-1, startPoints[spNumber][1]-1, b.factions[spNumber]))
	}
	// player faction settings
	//b.factions[0].resourcesMultiplier = 1
	//b.factions[0].buildSpeedMultiplier = 10
	// b.factions[0].visibilityCheat = true
	// b.factions[0].explorationCheat = true

	b.ais = append(b.ais, createAi(b.factions[0], "Player-side", "random"))
	for i := 1; i < len(b.factions); i++ {
		b.ais = append(b.ais, createAi(b.factions[i], fmt.Sprintf("Enemy %d", i), "random"))
	}
	for i := range b.ais {
		b.ais[i].controlsFaction.resourcesMultiplier = 4
		b.ais[i].controlsFaction.storagesMultiplier = 2
	}
	//b.ais[len(b.ais)-1].controlsFaction.resourcesMultiplier = 1.0
	//b.ais[len(b.ais)-1].controlsFaction.money = 5000

	//for i := 0; i < 14; i++ {
	//	b.addActor(createUnit(UNT_INFANTRY, startPoints[0][0]-1, startPoints[0][1], b.factions[0]))
	//}
	// unt.currentHitpoints = 1
	// unt.currentAction.setTargetTileCoords(10, 10)
}
