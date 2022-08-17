package main

import (
	"container/list"
	"dune2clone/astar"
	"dune2clone/geometry"
	"math"
)

type battlefield struct {
	tiles       [][]tile
	factions    []*faction
	ais         []*aiStruct
	buildings   list.List
	units       list.List
	projectiles list.List
	effects     list.List

	pathfinder  *astar.AStarPathfinder
	currentTick int
}

func (b *battlefield) tickToNonImportantRandom(mod int) int {
	return (b.currentTick / 73) % mod
}

func (b *battlefield) getSize() (int, int) {
	return len(b.tiles), len(b.tiles[0])
}

func (b *battlefield) areTileCoordsValid(x, y int) bool {
	return x >= 0 && y >= 0 && x < len(b.tiles) && y < len(b.tiles[0])
}

func (b *battlefield) performResourceGrowth() {
	for tx := range b.tiles {
		for ty := range b.tiles[tx] {
			if b.tiles[tx][ty].getStaticData().growsResources {
				// select resource tile to grow
				growX, growY := tx, ty
				if b.tiles[tx][ty].resourcesAmount >= RESOURCES_GROWTH_MAX {
					growX, growY = geometry.SpiralSearchForConditionFrom(
						func(x, y int) bool {
							return b.areTileCoordsValid(x, y) && b.tiles[x][y].getStaticData().canHaveResources &&
								b.tiles[x][y].resourcesAmount < RESOURCE_IN_TILE_MEDIUM_MAX &&
								geometry.GetApproxDistFromTo(tx, ty, x, y) <= RESOURCES_GROWTH_RADIUS
						},
						tx, ty, RESOURCES_GROWTH_RADIUS, rnd.Rand(4),
					)
				}
				if growX != -1 {
					b.tiles[growX][growY].resourcesAmount += rnd.RandInRange(RESOURCES_GROWTH_MIN, RESOURCES_GROWTH_MAX)
				}
			}
		}
	}
}

func (b *battlefield) canFactionSeeActor(f *faction, a actor) bool {
	switch a.(type) {
	case *unit:
		x, y := geometry.TrueCoordsToTileCoords(a.getPhysicalCenterCoords())
		if !b.areTileCoordsValid(x, y) {
			return false
		}
		return f.visibleTilesMap[x][y]
	case *building:
		bld := a.(*building)
		for tx := bld.topLeftX; tx < bld.topLeftX+bld.getStaticData().w; tx++ {
			for ty := bld.topLeftY; ty < bld.topLeftY+bld.getStaticData().h; ty++ {
				if f.visibleTilesMap[tx][ty] {
					return true
				}
			}
		}
	default:
		panic("wat")
	}
	return false
}

func (b *battlefield) changeTilesCodeInRectTo(x, y, w, h, code int) {
	for i := x; i < x+w; i++ {
		for j := y; j < j+h; j++ {
			if b.areTileCoordsValid(i, j) {
				b.tiles[i][j].code = code
			}
		}
	}
}

func (b *battlefield) addActor(a actor) {
	switch a.(type) {
	case *unit:
		b.units.PushBack(a)
	case *building:
		bld := a.(*building)
		if bld.getStaticData().givesFreeUnitOnCreation {
			x, y := bld.getUnitPlacementCoords()
			unt := createUnit(bld.getStaticData().codeForFreeUnitOnCreation, x, y, bld.getFaction())
			bld.unitPlacedInside = unt
		}
		b.buildings.PushBack(a)
	default:
		panic("wat")
	}
}

func (b *battlefield) removeActor(a actor) {
	switch a.(type) {
	case *unit:
		for i := b.units.Front(); i != nil; i = i.Next() {
			if i.Value == a {
				b.units.Remove(i)
			}
		}
	case *building:
		for i := b.buildings.Front(); i != nil; i = i.Next() {
			if i.Value == a {
				b.buildings.Remove(i)
			}
		}
	default:
		panic("wat")
	}
}

func (b *battlefield) addProjectile(p *projectile) {
	b.projectiles.PushFront(p)
}

func (b *battlefield) addEffect(e *effect) {
	b.effects.PushFront(e)
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

func (b *battlefield) getClosestEmptyFactionRefineryFromCoords(f *faction, x, y float64) actor {
	var selected actor = nil
	closestDist := math.MaxFloat64
	for i := b.buildings.Front(); i != nil; i = i.Next() {
		bld := i.Value.(*building)
		if bld.faction != f || !bld.getStaticData().receivesResources || bld.unitPlacedInside != nil {
			continue
		}
		bldCX, bldCY := bld.getPhysicalCenterCoords()
		distFromBld := geometry.SquareDistanceFloat64(x, y, bldCX, bldCY)
		if selected == nil || distFromBld < closestDist {
			closestDist = distFromBld
			selected = bld
		}
	}
	return selected
}

func (b *battlefield) isTileClearToBeMovedInto(x, y int, movingUnit *unit) bool {
	if !b.areTileCoordsValid(x, y) {
		return false
	}
	for i := b.buildings.Front(); i != nil; i = i.Next() {
		bld := i.Value.(*building)
		if bld.isPresentAt(x, y) {
			if bld.unitPlacedInside == nil {
				px, py := bld.getUnitPlacementCoords()
				if px == x && py == y {
					continue
				}
			}
			return false
		}
	}
	for i := b.units.Front(); i != nil; i = i.Next() {
		// debugWritef("req: %d,%d; act: %f, %f -> %d, %d \n", x, y, b.units[i].centerX, b.units[i].centerY, tx, ty)
		if movingUnit != nil && i.Value == movingUnit {
			continue
		}
		unt := i.Value.(*unit)
		if unt.getStaticData().isAircraft {
			continue
		}
		if unt.isPresentAt(x, y) {
			return false
		}
		if unt.currentAction.code == ACTION_MOVE {
			// x, y := geometry.TrueCoordsToTileCoords(unt.centerX, unt.centerY)
			if unt.currentAction.targetTileX == x && unt.currentAction.targetTileY == y {
				return false
			}
		}
	}
	return true
}

func (b *battlefield) costMapForMovement(x, y int) int {
	if !b.isTileClearToBeMovedInto(x, y, nil) {
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

func (b *battlefield) canBuildingBePlacedAt(placingBld *building, topLeftX, topLeftY, additionalDistance int, ignoreDistanceFromBase bool) bool {
	const MAX_MARGIN_FROM_EXISTING_BUILDING = 2
	_, _, w, h := placingBld.getDimensionsForConstructon()
	if additionalDistance > 0 {
		w = placingBld.getStaticData().w
		h = placingBld.getStaticData().h
	}
	for x := topLeftX - additionalDistance; x < topLeftX+w+additionalDistance; x++ {
		for y := topLeftY - additionalDistance; y < topLeftY+h+additionalDistance; y++ {
			if !b.areTileCoordsValid(x, y) {
				return false
			}
			if !b.tiles[x][y].getStaticData().canBuildHere {
				return false
			}
			if b.costMapForMovement(x, y) == -1 {
				return false
			}
			for i := b.buildings.Front(); i != nil; i = i.Next() {
				if bld, ok := i.Value.(*building); ok {
					bx, by, bw, bh := bld.getDimensionsForConstructon()
					if geometry.AreCoordsInTileRect(x, y, bx, by, bw, bh) {
						return false
					}
				}
			}
		}
	}
	if ignoreDistanceFromBase {
		return true
	}
	// check distance requirement
	for i := b.buildings.Front(); i != nil; i = i.Next() {
		if bld, ok := i.Value.(*building); ok {
			if bld.getFaction() != placingBld.getFaction() {
				continue
			}
			bx, by, bw, bh := bld.topLeftX, bld.topLeftY, bld.getStaticData().w, bld.getStaticData().h
			if geometry.AreRectsInRange(bx, by, bw, bh, topLeftX, topLeftY, placingBld.getStaticData().w, placingBld.getStaticData().h, MAX_MARGIN_FROM_EXISTING_BUILDING) {
				return true
			}
		}
	}
	return false
}

func (b *battlefield) getCoordsOfClosestEmptyTileWithResourcesTo(tx, ty int) (int, int) {
	// TODO: optimize this shit
	return geometry.SpiralSearchForConditionFrom(
		func(x, y int) bool {
			return b.areTileCoordsValid(x, y) &&
				b.tiles[x][y].resourcesAmount > 0 &&
				b.isTileClearToBeMovedInto(x, y, nil)
		},
		tx, ty, len(b.tiles)/2, rnd.Rand(4))
}
