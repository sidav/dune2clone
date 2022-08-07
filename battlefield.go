package main

import (
	"container/list"
	"dune2clone/astar"
	"dune2clone/geometry"
)

type battlefield struct {
	tiles       [][]tile
	factions    []*faction
	buildings   list.List
	units       list.List
	projectiles list.List

	pathfinder  *astar.AStarPathfinder
	currentTick int
}

func (b *battlefield) create(w, h int) {
	b.tiles = make([][]tile, w)
	for i := range b.tiles {
		b.tiles[i] = make([]tile, h)
		for j := range b.tiles[i] {
			b.tiles[i][j].code = TILE_SAND
		}
	}
	b.pathfinder = &astar.AStarPathfinder{
		DiagonalMoveAllowed:       false,
		ForceGetPath:              true,
		ForceIncludeFinish:        false,
		AutoAdjustDefaultMaxSteps: false,
		MapWidth:                  len(b.tiles),
		MapHeight:                 len(b.tiles[0]),
	}
	b.placeInitialStuff()
	b.finalizeTileVariants()
}

func (b *battlefield) finalizeTileVariants() {
	for i := range b.tiles {
		for j := range b.tiles[i] {
			b.tiles[i][j].spriteVariantIndex = rnd.Rand(len(sTableTiles[b.tiles[i][j].code].spriteCodes))
		}
	}
}

func (b *battlefield) placeInitialStuff() {
	for x := 7; x < 10; x++ {
		for y := 0; y < 3; y++ {
			b.tiles[x][y].resourcesAmount = rnd.RandInRange(100, 300)
		}
	}
	b.factions = append(b.factions, &faction{
		factionColor: factionTints[0],
		money:        10000,
		team:         0,
	})
	b.factions = append(b.factions, &faction{
		factionColor: factionTints[1],
		money:        10000,
		team:         0,
	})

	b.addActor(createBuilding(BLD_BASE, 1, 1, b.factions[0]))
	b.addActor(createBuilding(BLD_BASE, 14, 8, b.factions[1]))

	b.addActor(createUnit(UNT_TANK, 3, 3, b.factions[0]))
	b.addActor(createUnit(UNT_HARVESTER, 4, 3, b.factions[0]))
	b.addActor(createUnit(UNT_TANK, 13, 7, b.factions[1]))
}

func (b *battlefield) addActor(a actor) {
	switch a.(type) {
	case *unit:
		b.units.PushBack(a)
	case *building:
		b.buildings.PushBack(a)
	default:
		panic("wat")
	}
}

func (b *battlefield) addProjectile(p *projectile) {
	b.projectiles.PushFront(p)
}

func (b *battlefield) getActorAtTileCoordinates(x, y int) actor {
	for i := b.buildings.Front(); i != nil; i = i.Next() {
		if i.Value.(*building).isPresentAt(x, y) {
			// debugWrite("got")
			return i.Value.(actor)
		}
	}
	for i := b.units.Front(); i != nil; i = i.Next() {
		// debugWritef("req: %d,%d; act: %f, %f -> %d, %d \n", x, y, b.units[i].centerX, b.units[i].centerY, tx, ty)
		if i.Value.(*unit).isPresentAt(x, y) {
			// debugWrite("got")
			return i.Value.(actor)
		}
	}
	return nil
}

func (b *battlefield) costMapForMovement(x, y int) int {
	act := b.getActorAtTileCoordinates(x, y)
	if act != nil {
		// debugWritef("At coords %d,%d there is %+v", x, y, act)
		return -1
	}
	return 10
}

func (b *battlefield) getListOfActorsInRangeFrom(x, y, r int) *list.List {
	lst := list.List{}
	for i := b.units.Front(); i != nil; i = i.Next() {
		u := i.Value.(*unit)
		tx, ty := geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
		if geometry.AreCoordsInRange(tx, ty, x, y, r) {
			lst.PushBack(u)
		}
	}
	for i := b.buildings.Front(); i != nil; i = i.Next() {
		bld := i.Value.(*building)
		if geometry.AreCoordsInRangeFromRect(x, y, bld.topLeftX, bld.topLeftY, bld.getStaticData().w, bld.getStaticData().h, r) {
			lst.PushBack(bld)
		}
	}
	return &lst
}

func (b *battlefield) findPathForUnitTo(u *unit, tileX, tileY int, forceIncludeFinish bool) *astar.Cell {
	utx, uty := geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
	b.pathfinder.ForceIncludeFinish = forceIncludeFinish
	return b.pathfinder.FindPath(
		func(x, y int) int {
			return b.costMapForMovement(x, y)
		},
		utx, uty, tileX, tileY,
	)
}

func (b *battlefield) isRectClearForBuilding(topLeftX, topLeftY, w, h int) bool {
	for x := topLeftX; x < topLeftX+w; x++ {
		for y := topLeftY; y < topLeftY+h; y++ {
			if b.costMapForMovement(x, y) == -1 {
				return false
			}
		}
	}
	return true
}

func (b *battlefield) getCoordsOfClosestEmptyTileWithResourcesTo(tx, ty int) (int, int) {
	// TODO: optimize this shit
	lowestRange := len(b.tiles)*len(b.tiles) + len(b.tiles[0])*len(b.tiles[0])
	currX, currY := -1, -1
	for x := range b.tiles {
		for y := range b.tiles[x] {
			currRange := (tx-x)*(tx-x)+(ty-y)*(ty-y)
			if b.tiles[x][y].resourcesAmount > 0 && currRange < lowestRange {
				currX, currY = x, y
				lowestRange = currRange
			}
		}
	}
	return currX, currY
}
