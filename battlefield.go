package main

import (
	"container/list"
	"dune2clone/astar"
)

type battlefield struct {
	tiles       [][]tile
	factions    []*faction
	buildings   []*building
	units       []*unit
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
}

func (b *battlefield) placeInitialStuff() {
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
	b.addActor(createUnit(UNT_TANK, 13, 7, b.factions[1]))
}

func (b *battlefield) addActor(a actor) {
	switch a.(type) {
	case *unit:
		b.units = append(b.units, a.(*unit))
	case *building:
		b.buildings = append(b.buildings, a.(*building))
	default:
		panic("wat")
	}
}

func (b *battlefield) addProjectile(p *projectile) {
	b.projectiles.PushFront(p)
}

func (b *battlefield) getActorAtTileCoordinates(x, y int) actor {
	for i := range b.buildings {
		if b.buildings[i].isPresentAt(x, y) {
			// debugWrite("got")
			return b.buildings[i]
		}
	}
	for i := range b.units {
		// debugWritef("req: %d,%d; act: %f, %f -> %d, %d \n", x, y, b.units[i].centerX, b.units[i].centerY, tx, ty)
		if b.units[i].isPresentAt(x, y) {
			// debugWrite("got")
			return b.units[i]
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
	for _, u := range b.units {
		tx, ty := trueCoordsToTileCoords(u.centerX, u.centerY)
		if AreCoordsInRange(tx, ty, x, y, r) {
			lst.PushBack(u)
		}
	}
	for _, bld := range b.buildings {
		if AreCoordsInRangeFromRect(x, y, bld.topLeftX, bld.topLeftY, bld.getStaticData().w, bld.getStaticData().h, r) {
			lst.PushBack(bld)
		}
	}
	return &lst
}

func AreRectsInRange(x1, y1, w1, h1, x2, y2, w2, h2, r int) bool {
	// all -1's are beacuse of TILED geometry
	x1b := x1+w1-1
	x2b := x2+w2-1
	y1b := y1+h1-1
	y2b := y2+h2-1

	left := x2b < x1
	right := x1b < x2
	bottom := y1b < y2
	top := y2b < y1
	if top && left {
		return AreCoordsInRange(x1, y1, x2b, y2b, r) // dist((x1, y1b), (x2b, y2))
	}
	if left && bottom {
		return AreCoordsInRange(x1, y1b, x2b, y2, r)
	}
	if bottom && right {
		return AreCoordsInRange(x1b, y1b, x2, y2, r)
	}
	if right && top {
		return AreCoordsInRange(x1b, y1, x2, y2b, r)
	}
	if left {
		return x1 - x2b <= r
	}
	if right {
		return x2 - x1b <= r
	}
	if bottom {
		return y2 - y1b <= r
	}
	if top {
		return y1 - y2b <= r
	}
	return true // intersect detected
}

func (b *battlefield) findPathForUnitTo(u *unit, tileX, tileY int) *astar.Cell {
	utx, uty := trueCoordsToTileCoords(u.centerX, u.centerY)
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
