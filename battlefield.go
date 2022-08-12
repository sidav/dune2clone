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

	pathfinder  *astar.AStarPathfinder
	currentTick int
}

func (b *battlefield) areTileCoordsValid(x, y int) bool {
	return x > 0 && y > 0 && x < len(b.tiles) && y < len(b.tiles[0])
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

func (b *battlefield) canBuildingBePlacedAt(placingBld *building, topLeftX, topLeftY int, ignoreDistanceFromBase bool) bool {
	const MAX_MARGIN_FROM_EXISTING_BUILDING = 2
	_, _, w, h := placingBld.getDimensionsForConstructon()
	for x := topLeftX; x < topLeftX+w; x++ {
		for y := topLeftY; y < topLeftY+h; y++ {
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
			if geometry.AreRectsInRange(bx, by, bw, bh, topLeftX, topLeftY, w, h, MAX_MARGIN_FROM_EXISTING_BUILDING) {
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
			return areTileCoordsValid(x, y) &&
				b.tiles[x][y].resourcesAmount > 0 &&
				b.isTileClearToBeMovedInto(x, y, nil)
		},
		tx, ty, len(b.tiles)/2, rnd.Rand(4))
}
