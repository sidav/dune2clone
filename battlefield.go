package main

import (
	"dune2clone/astar"
)

type battlefield struct {
	tiles     [][]tile
	buildings []*building
	units     []*unit

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
	f1 := &faction{
		factionColor: factionTints[0],
		money:        10000,
		team:         0,
	}
	f2 := &faction{
		factionColor: factionTints[1],
		money:        10000,
		team:         0,
	}
	b.buildings = append(b.buildings, &building{
		topLeftX: 1,
		topLeftY: 1,
		code:     BLD_BASE,
		faction: f1,
	})

	b.buildings = append(b.buildings, &building{
		topLeftX: 5,
		topLeftY: 3,
		code:     BLD_POWERPLANT,
		faction: f2,
	})
	b.buildings = append(b.buildings, &building{
		topLeftX: 8,
		topLeftY: 7,
		code:     BLD_FACTORY,
		faction: f2,
	})

	b.units = append(b.units, &unit{
		code:    UNT_TANK,
		centerX: 4.5,
		centerY: 4.5,
		faction: f1,
	})
	b.units = append(b.units, &unit{
		code:    UNT_TANK,
		centerX: 5.5,
		centerY: 5.5,
		faction: f2,
	})
}

func (b *battlefield) getActorAtTileCoordinates(x, y int) actor {
	for i := range b.buildings {
		if b.buildings[i].isPresentAt(x, y) {
			// debugWrite("got")
			return b.buildings[i]
		}
	}
	for i := range b.units {
		tx, ty := trueCoordsToTileCoords(b.units[i].centerX, b.units[i].centerY)
		// debugWritef("req: %d,%d; act: %f, %f -> %d, %d \n", x, y, b.units[i].centerX, b.units[i].centerY, tx, ty)
		if tx == x && ty == y {
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

func (b *battlefield) findPathForUnitTo(u *unit, tileX, tileY int) *astar.Cell {
	utx, uty := trueCoordsToTileCoords(u.centerX, u.centerY)
	return b.pathfinder.FindPath(
		func(x, y int) int {
			return b.costMapForMovement(x, y)
		},
		utx, uty, tileX, tileY,
	)
}
