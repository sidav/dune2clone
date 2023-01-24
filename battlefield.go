package main

import (
	"container/list"
	"dune2clone/astar"
	"dune2clone/geometry"
	"fmt"
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
			if b.tiles[tx][ty].hasResourceVein {
				// select resource tile to grow
				growX, growY := -1, -1
				// if b.tiles[tx][ty].resourcesAmount >= RESOURCES_GROWTH_MAX {
				growX, growY = geometry.SpiralSearchForClosestConditionFrom(
					func(x, y int) bool {
						return b.areTileCoordsValid(x, y) && b.tiles[x][y].canHaveResources() &&
							b.tiles[x][y].resourcesAmount < RESOURCE_IN_TILE_MEDIUM_MAX &&
							geometry.GetApproxDistFromTo(tx, ty, x, y) <= RESOURCES_GROWTH_RADIUS
					},
					tx, ty, RESOURCES_GROWTH_RADIUS, rnd.Rand(4),
				)
				// }
				if growX != -1 {
					b.tiles[growX][growY].resourcesAmount += rnd.RandInRange(RESOURCES_GROWTH_MIN, RESOURCES_GROWTH_MAX)
				}
			}
		}
	}
}

func (b *battlefield) hasFactionExploredBuilding(f *faction, bld *building) bool {
	for tx := bld.topLeftX; tx < bld.topLeftX+bld.getStaticData().w; tx++ {
		for ty := bld.topLeftY; ty < bld.topLeftY+bld.getStaticData().h; ty++ {
			if f.hasTileAtCoordsExplored(tx, ty) {
				return true
			}
		}
	}
	return false
}

func (b *battlefield) canFactionSeeActor(f *faction, a actor) bool {
	switch a.(type) {
	case *unit:
		x, y := geometry.TrueCoordsToTileCoords(a.getPhysicalCenterCoords())
		if !b.areTileCoordsValid(x, y) {
			return false
		}
		return f.seesTileAtCoords(x, y)
	case *building:
		bld := a.(*building)
		for tx := bld.topLeftX; tx < bld.topLeftX+bld.getStaticData().w; tx++ {
			for ty := bld.topLeftY; ty < bld.topLeftY+bld.getStaticData().h; ty++ {
				if f.seesTileAtCoords(tx, ty) {
					return true
				}
			}
		}
		return false
	default:
		panic("wat")
	}
	return false
}

func (b *battlefield) changeTilesCodesInRectTo(x, y, w, h, code int) {
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			if b.areTileCoordsValid(i, j) {
				b.tiles[i][j].code = code
			}
		}
	}
}

func (b *battlefield) addActor(a actor) {
	switch a.(type) {
	case *unit:
		a.getFaction().gameStatistics.totalProduced++
		b.units.PushBack(a)
		x, y := a.(*unit).getTileCoords()
		if !a.isInAir() {
			b.setTilesOccupiedByActor(x, y, 1, 1, a)
		}
	case *building:
		a.getFaction().gameStatistics.totalBuilt++
		bld := a.(*building)
		if bld.getStaticData().givesFreeUnitOnCreation {
			x, y := bld.getUnitPlacementAbsoluteCoords()
			unt := createUnit(bld.getStaticData().codeForFreeUnitOnCreation, x, y, bld.getFaction())
			unt.chassisDegree = 90 // looking down
			bld.unitPlacedInside = unt
		}
		b.setTilesOccupiedByActor(bld.topLeftX, bld.topLeftY, bld.getStaticData().w, bld.getStaticData().h, a)
		b.buildings.PushBack(a)
	default:
		panic("wat")
	}
}

func (b *battlefield) isActorPresentOnBattlefield(a actor) bool {
	switch a.(type) {
	case *unit:
		for i := b.units.Front(); i != nil; i = i.Next() {
			if i.Value == a {
				return true
			}
		}
	case *building:
		for i := b.buildings.Front(); i != nil; i = i.Next() {
			if i.Value == a {
				return true
			}
		}
	default:
		panic("wat")
	}
	return false
}

func (b *battlefield) removeActor(a actor) {
	removals := 0          // for debugging
	var next *list.Element // for deletion while iterating
	switch a.(type) {
	case *unit:
		for i := b.units.Front(); i != nil; i = next {
			next = i.Next()
			if i.Value == a {
				b.clearTileOccupationForActor(a)
				b.units.Remove(i)
			}
		}
	case *building:
		for i := b.buildings.Front(); i != nil; i = next {
			next = i.Next()
			if i.Value == a {
				b.clearTileOccupationForActor(a)
				b.buildings.Remove(i)
			}
		}
	default:
		panic("wat")
	}
	if removals > 1 {
		debugWritef("WAT?! %s removed %d times\n", a.getName(), removals)
	}
}

func (b *battlefield) addProjectile(p *projectile) {
	b.projectiles.PushFront(p)
}

func (b *battlefield) addEffect(e *effect) {
	b.effects.PushFront(e)
}

func (b *battlefield) setTilesOccupiedByActor(x, y, w, h int, a actor) {
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			b.tiles[i][j].isOccupiedByActor = a
		}
	}
	if bld, ok := a.(*building); ok {
		if bld.getStaticData().canUnitBePlacedIn() && bld.unitPlacedInside == nil {
			freeX, freeY := bld.getStaticData().unitPlacementX, bld.getStaticData().unitPlacementY
			b.tiles[x+freeX][y+freeY].isOccupiedByActor = nil
		}
	}
}

func (b *battlefield) clearTilesOccupationInRect(x, y, w, h int) {
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			b.tiles[i][j].isOccupiedByActor = nil
		}
	}
}

func (b *battlefield) clearTileOccupationForActor(a actor) {
	if unt, ok := a.(*unit); ok {
		if !unt.isInAir() {
			x, y := unt.getTileCoords()
			b.clearTilesOccupationInRect(x, y, 1, 1)
			// just to be sure, maybe clear a tile at which the unit moves
			// The logic ON PURPOSE doesn't look at unit's action
			// check unit displacement
			tileCx, tileCy := geometry.TileCoordsToTrueCoords(x, y)
			ux, uy := unt.getPhysicalCenterCoords()
			dispX, dispY := geometry.Float64VectorToIntUnitVector(ux-tileCx, uy-tileCy)
			// clear the tile if it is set occupied by this unit
			if b.tiles[x+dispX][y+dispY].isOccupiedByActor == unt {
				b.clearTilesOccupationInRect(x+dispX, y+dispY, 1, 1)
			}
		}
	}
	if bld, ok := a.(*building); ok {
		b.clearTilesOccupationInRect(bld.topLeftX, bld.topLeftY, bld.getStaticData().w, bld.getStaticData().h)
	}
}

func (b *battlefield) getActorAtTileCoordinates(x, y int) actor {
	if b.areTileCoordsValid(x, y) && b.tiles[x][y].isOccupiedByActor != nil {
		return b.tiles[x][y].isOccupiedByActor
	}
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

func (b *battlefield) isTileClearToBeMovedInto(x, y int, movingUnit *unit) bool {
	if !b.areTileCoordsValid(x, y) {
		return false
	}
	if !b.tiles[x][y].getStaticData().canBeWalkedOn {
		return false
	}
	if b.tiles[x][y].isOccupiedByActor != nil && b.tiles[x][y].isOccupiedByActor != movingUnit {
		return false
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

//func (b *battlefield) getListOfActorsInRangeFromActor(x, y, r int) *list.List {
//	lst := list.List{}
//	for i := b.units.Front(); i != nil; i = i.Next() {
//		u := i.Value.(*unit)
//		tx, ty := geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
//		if geometry.GetApproxDistFromTo(tx, ty, x, y) <= r {
//			lst.PushBack(u)
//		}
//	}
//	for i := b.buildings.Front(); i != nil; i = i.Next() {
//		bld := i.Value.(*building)
//		if geometry.AreCoordsInRangeFromRect(x, y, bld.topLeftX, bld.topLeftY, bld.getStaticData().w, bld.getStaticData().h, r) {
//			lst.PushBack(bld)
//		}
//	}
//	return &lst
//}

func (b *battlefield) getListOfActorsInRangeFromActor(a actor, r int) *list.List {
	lst := list.List{}
	for i := b.units.Front(); i != nil; i = i.Next() {
		u := i.Value.(*unit)
		if b.areActorsInRangeFromEachOther(a, u, r) {
			lst.PushBack(u)
		}
	}
	for i := b.buildings.Front(); i != nil; i = i.Next() {
		bld := i.Value.(*building)
		if b.areActorsInRangeFromEachOther(a, bld, r) {
			lst.PushBack(bld)
		}
	}
	return &lst
}

func (b *battlefield) areActorsInRangeFromEachOther(a1, a2 actor, r int) bool {
	switch a1.(type) {
	case *unit:
		switch a2.(type) {
		case *unit:
			x1, y1 := a1.getPhysicalCenterCoords()
			x2, y2 := a2.getPhysicalCenterCoords()
			return int(geometry.GetApproxDistFloat64(x1, y1, x2, y2)) <= r
		case *building:
			x1, y1 := a1.(*unit).getTileCoords()
			bld := a2.(*building)
			x2, y2, w, h := bld.topLeftX, bld.topLeftY, bld.getStaticData().w, bld.getStaticData().h
			return geometry.AreRectsInRange(x1, y1, 1, 1, x2, y2, w, h, r)
		}
	case *building:
		switch a2.(type) {
		case *unit:
			unitX, unitY := a2.(*unit).getTileCoords()
			bld := a1.(*building)
			x2, y2, w, h := bld.topLeftX, bld.topLeftY, bld.getStaticData().w, bld.getStaticData().h
			return geometry.AreRectsInRange(unitX, unitY, 1, 1, x2, y2, w, h, r)
		case *building:
			bld1 := a1.(*building)
			x1, y1, w1, h1 := bld1.topLeftX, bld1.topLeftY, bld1.getStaticData().w, bld1.getStaticData().h
			bld2 := a2.(*building)
			x2, y2, w2, h2 := bld2.topLeftX, bld2.topLeftY, bld2.getStaticData().w, bld2.getStaticData().h
			return geometry.AreRectsInRange(x1, y1, w1, h1, x2, y2, w2, h2, r)
		}
	}
	panic("How is it possible?")
}

func (b *battlefield) getListOfActorsInTilesRect(x, y, w, h int) *list.List {
	lst := list.List{}
	for i := b.units.Front(); i != nil; i = i.Next() {
		u := i.Value.(*unit)
		tx, ty := geometry.TrueCoordsToTileCoords(u.centerX, u.centerY)
		if geometry.AreCoordsInTileRect(tx, ty, x, y, w, h) {
			lst.PushBack(u)
		}
	}
	for i := b.buildings.Front(); i != nil; i = i.Next() {
		bld := i.Value.(*building)
		if geometry.AreTwoCellRectsOverlapping(x, y, w, h,
			bld.topLeftX, bld.topLeftY, bld.getStaticData().w, bld.getStaticData().h) {

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

func (b *battlefield) canUnitBeDeployedAt(u *unit, x, y int) bool {
	// tx, ty := u.getTileCoords()
	bld := createBuilding(u.getStaticData().deploysInto, 0, 0, u.faction)
	b.removeActor(u)
	can := b.canBuildingBePlacedAt(bld, x, y, 0, true)
	b.addActor(u)
	return can
}

func (b *battlefield) canActorAttackActor(attacker, target actor) bool {
	switch attacker.(type) {
	case *building:
		if attacker.(*building).turret == nil {
			return false
		}
		return b.canTurretAttackActor(attacker.(*building).turret, target)
	case *unit:
		if len(attacker.(*unit).turrets) == 0 {
			return false
		}
		// TODO: consider other turrets too?
		return b.canTurretAttackActor(attacker.(*unit).turrets[0], target)
	}
	panic("wat")
}

func (b *battlefield) canBuildingBePlacedAt(placingBld *building, topLeftX, topLeftY,
	forcedDistanceFromOtherBuildings int, ignoreDistanceFromBase bool) bool {

	const MAX_MARGIN_FROM_EXISTING_BUILDING = 2
	_, _, w, h := placingBld.getDimensionsForConstructon()
	if forcedDistanceFromOtherBuildings > 0 {
		w = placingBld.getStaticData().w
		h = placingBld.getStaticData().h
	}
	for x := topLeftX - forcedDistanceFromOtherBuildings; x < topLeftX+w+forcedDistanceFromOtherBuildings; x++ {
		for y := topLeftY - forcedDistanceFromOtherBuildings; y < topLeftY+h+forcedDistanceFromOtherBuildings; y++ {
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
			if geometry.AreRectsInTaxicabRange(bx, by, bw, bh, topLeftX, topLeftY, placingBld.getStaticData().w, placingBld.getStaticData().h, MAX_MARGIN_FROM_EXISTING_BUILDING) {
				return true
			}
		}
	}
	return false
}

func (b *battlefield) getCoordsOfClosestEmptyTileWithResourcesTo(tx, ty int) (int, int) {
	return geometry.SpiralSearchForClosestConditionFrom(
		func(x, y int) bool {
			return b.areTileCoordsValid(x, y) &&
				b.tiles[x][y].resourcesAmount > 0 &&
				b.isTileClearToBeMovedInto(x, y, nil)
		},
		tx, ty, len(b.tiles)/2, rnd.Rand(4))
}

// not intended to be optimized, so don't call it frequently.
func (b *battlefield) collectStatisticsForDebug() string {
	str := ""
	counter := 0
	for i := b.units.Front(); i != nil; i = i.Next() {
		counter++
	}
	str += fmt.Sprintf("Total units: %d ", counter)
	counter = 0
	for i := b.buildings.Front(); i != nil; i = i.Next() {
		counter++
	}
	str += fmt.Sprintf("buildings: %d", counter)
	return str
}
