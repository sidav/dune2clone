package main

import (
	"dune2clone/astar"
	"dune2clone/map_generator"
)

func (b *battlefield) initFromRandomMap(rm *map_generator.GameMap) {
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
				currTileCode = TILE_SAND
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
				currTileCode = TILE_RESOURCE_VEIN
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
		b.factions = append(b.factions, createFaction(spNumber, 0, 10000, 2, 1))
		b.factions[spNumber].resetVisibilityMaps(len(b.tiles), len(b.tiles[0]))
		b.factions[spNumber].exploreAround(startPoints[spNumber][0], startPoints[spNumber][1], 2, 2, 3)
		b.addActor(createBuilding(BLD_BASE, startPoints[spNumber][0], startPoints[spNumber][1], b.factions[spNumber]))
		// b.addActor(createUnit(UNT_HARVESTER, startPoints[spNumber][0]-1, startPoints[spNumber][1]-1, b.factions[spNumber]))
	}
	// player faction settings
	b.factions[0].resourcesMultiplier = 1
	b.factions[0].buildSpeedMultiplier = 1
	// b.factions[0].visibilityCheat = true
	// b.factions[0].explorationCheat = true

	// b.ais = append(b.ais, createAi(b.factions[0], "Player-side"))
	// b.ais = append(b.ais, createAi(b.factions[1], "Enemy"))

	// unt := createUnit(AIR_TRANSPORT, startPoints[0][0], startPoints[0][1], b.factions[0])
	// unt.currentAction.code = ACTION_AIR_APPROACH_LAND_TILE
	// unt.currentAction.setTargetTileCoords(10, 10)
	// b.addActor(unt)
}
