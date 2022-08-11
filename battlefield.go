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
	buildings   list.List
	units       list.List
	projectiles list.List

	pathfinder  *astar.AStarPathfinder
	currentTick int
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
		if !bld.getStaticData().receivesResources || bld.unitPlacedInside != nil {
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
	for i := b.buildings.Front(); i != nil; i = i.Next() {
		if i.Value.(*building).isPresentAt(x, y) {
			// debugWrite("got")
			return false
		}
	}
	for i := b.units.Front(); i != nil; i = i.Next() {
		// debugWritef("req: %d,%d; act: %f, %f -> %d, %d \n", x, y, b.units[i].centerX, b.units[i].centerY, tx, ty)
		if movingUnit != nil && i.Value == movingUnit {
			continue
		}
		unt := i.Value.(*unit)
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
			currRange := geometry.SquareDistanceInt(x, y, tx, ty)
			if b.tiles[x][y].resourcesAmount > 0 && currRange < lowestRange && b.isTileClearToBeMovedInto(x, y, nil) {
				currX, currY = x, y
				lowestRange = currRange
			}
		}
	}
	return currX, currY
}
