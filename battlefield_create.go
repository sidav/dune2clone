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
	b.factions = append(b.factions, &faction{
		factionColor:        factionTints[0],
		money:               10000,
		team:                0,
		resourcesMultiplier: 1,
	})
	b.factions = append(b.factions, &faction{
		factionColor:        factionTints[1],
		money:               10000,
		team:                0,
		resourcesMultiplier: 10,
	})

	b.addActor(createBuilding(BLD_BASE, startPoints[0][0], startPoints[0][1], b.factions[0]))
	b.addActor(createUnit(UNT_TANK, startPoints[0][0]+2, startPoints[0][1]+2, b.factions[0]))
	b.addActor(createUnit(UNT_QUAD, startPoints[0][0]+3, startPoints[0][1]+2, b.factions[0]))

	b.addActor(createBuilding(BLD_BASE, startPoints[1][0], startPoints[1][1], b.factions[1]))
	b.addActor(createUnit(UNT_TANK, startPoints[1][0]-1, startPoints[1][1]-1, b.factions[1]))
}
