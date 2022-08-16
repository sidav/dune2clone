package main

import (
	"dune2clone/astar"
	"dune2clone/map_generator"
)

//func (b *battlefield) initEmpty(w, h int) {
//	b.tiles = make([][]tile, w)
//	for i := range b.tiles {
//		b.tiles[i] = make([]tile, h)
//		for j := range b.tiles[i] {
//			b.tiles[i][j].code = TILE_SAND
//		}
//	}
//	b.pathfinder = &astar.AStarPathfinder{
//		DiagonalMoveAllowed:       true,
//		ForceGetPath:              true,
//		ForceIncludeFinish:        false,
//		AutoAdjustDefaultMaxSteps: false,
//		MapWidth:                  len(b.tiles),
//		MapHeight:                 len(b.tiles[0]),
//	}
//	// place some resources
//	for x := MAP_W/5; x < 4*MAP_W/5; x++ {
//		for y := 0; y < MAP_H; y++ {
//			b.tiles[x][y].resourcesAmount = rnd.RandInRange(250, 500)
//		}
//	}
//	b.placeInitialStuff()
//	b.finalizeTileVariants()
//}

func (b *battlefield) initFromRandomMap(rm *map_generator.GameMap) {
	b.tiles = make([][]tile, len(rm.Tiles))
	for i := range b.tiles {
		b.tiles[i] = make([]tile, len(rm.Tiles[i]))
		for j := range b.tiles[i] {
			var currTileCode int
			switch rm.Tiles[i][j] {
			case map_generator.SAND:
				currTileCode = map_generator.SAND
			case map_generator.BUILDABLE_TERRAIN:
				currTileCode = map_generator.BUILDABLE_TERRAIN
			case map_generator.RESOURCES:
				currTileCode = map_generator.SAND
				b.tiles[i][j].resourcesAmount = rnd.RandInRange(100, 300)
			default:
				currTileCode = map_generator.SAND
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
			b.tiles[i][j].spriteVariantIndex = rnd.Rand(len(sTableTiles[b.tiles[i][j].code].spriteCodes))
		}
	}
}

func (b *battlefield) placeInitialStuff(startPoints [][2]int) {
	for spNumber := range startPoints {
		b.factions = append(b.factions, createFaction(spNumber, 0, 10000, 2))
		b.factions[spNumber].resetVisibilityMaps(len(b.tiles), len(b.tiles[0]))
		b.factions[spNumber].exploreAround(startPoints[spNumber][0], startPoints[spNumber][1], 2, 2, 3)
		b.addActor(createBuilding(BLD_BASE, startPoints[spNumber][0], startPoints[spNumber][1], b.factions[spNumber]))
		b.addActor(createUnit(AIR_TRANSPORT, startPoints[spNumber][0], startPoints[spNumber][1], b.factions[spNumber]))
		b.addActor(createUnit(UNT_HARVESTER, startPoints[spNumber][0]-1, startPoints[spNumber][1]-1, b.factions[spNumber]))
	}
	b.factions[0].resourcesMultiplier = 1 // for player
	//b.ais = append(b.ais, createAi(b.factions[0], "Player-side"))
	//b.ais = append(b.ais, createAi(b.factions[1], "Enemy"))
}
