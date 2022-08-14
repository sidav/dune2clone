package main

import (
	"dune2clone/geometry"
	"image/color"
	"math"
)

type faction struct {
	colorNumber                         int
	currentResources, resourceStorage   float64
	money                               float64 // float because of division when spending
	energyProduction, energyConsumption int

	team                int // 0 means "enemy to all"
	resourcesMultiplier float64

	exploredTilesMap, visibleTilesMap [][]bool
}

func createFaction(colorNumber, team int, initialMoney, resourcesMultiplier float64) *faction {
	f := &faction{
		colorNumber:         colorNumber,
		money:               initialMoney,
		energyProduction:    999, // will be overwritten anyway
		team:                team,
		resourcesMultiplier: resourcesMultiplier,
	}
	return f
}

func (f *faction) resetVisibilityMaps(mapW, mapH int) {
	if f.visibleTilesMap == nil {
		f.visibleTilesMap = make([][]bool, mapW)
		for i := range f.visibleTilesMap {
			f.visibleTilesMap[i] = make([]bool, mapH)
		}
	}
	for x := range f.visibleTilesMap {
		for y := range f.visibleTilesMap[x] {
			f.visibleTilesMap[x][y] = false
		}
	}

	if f.exploredTilesMap == nil {
		f.exploredTilesMap = make([][]bool, mapW)
		for i := range f.exploredTilesMap {
			f.exploredTilesMap[i] = make([]bool, mapH)
		}
	}
	//for x := range f.exploredTilesMap {
	//	for y := range f.exploredTilesMap[x] {
	//		f.exploredTilesMap[x][y] = false
	//	}
	//}
}

func (f *faction) exploreAround(tileX, tileY, fromW, fromH, radius int) {
	for x := tileX - radius; x < tileX+fromW+radius; x++ {
		for y := tileY - radius; y < tileY+fromH+radius; y++ {
			if x >= 0 && x < len(f.visibleTilesMap) && y >= 0 && y < len(f.visibleTilesMap) {
				if geometry.GetApproxDistFromTo(x, y, tileX, tileY) <= radius ||
					geometry.GetApproxDistFromTo(x, y, tileX+fromW-1, tileY+fromH-1) <= radius ||
					geometry.GetApproxDistFromTo(x, y, tileX+fromW-1, tileY) <= radius ||
					geometry.GetApproxDistFromTo(x, y, tileX, tileY+fromH-1) <= radius {

					f.exploredTilesMap[x][y] = true
					f.visibleTilesMap[x][y] = true
				}
			}
		}
	}
}

func (f *faction) canSeeActor(a actor) bool {
	switch a.(type) {
	case *unit:
		x, y := geometry.TrueCoordsToTileCoords(a.getPhysicalCenterCoords())
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

func (f *faction) getMoney() float64 {
	return f.currentResources + f.money
}

func (f *faction) getStorageRemaining() float64 {
	return f.resourceStorage - f.currentResources
}

func (f *faction) getAvailableEnergy() int {
	//if f.energyProduction <= f.energyConsumption {
	//	return 0
	//}
	return f.energyProduction - f.energyConsumption
}

func (f *faction) getEnergyProductionMultiplier() float64 {
	if f.energyProduction >= f.energyConsumption {
		return 1.0
	}
	factor := float64(f.energyProduction) / float64(f.energyConsumption)
	if factor < 0.25 {
		factor = 0.25
	}
	return factor
}

func (f *faction) spendMoney(spent float64) {
	if f.currentResources > 0 {
		spentFromResources := math.Min(f.currentResources, spent)
		f.currentResources -= spentFromResources
		spent -= spentFromResources
	}
	if spent > 0 {
		f.money -= spent
	}
}
func (f *faction) receiveResources(amount float64, asMoney bool) {
	if asMoney {
		f.money += amount * f.resourcesMultiplier
	} else {
		f.currentResources += amount * f.resourcesMultiplier
		if f.currentResources > f.resourceStorage {
			f.currentResources = f.resourceStorage
		}
	}
}

func (f *faction) resetCurrents() {
	f.resourceStorage = 0
	f.energyProduction = 0
	f.energyConsumption = 0
}

const zeroTiltColor = 32
const strongerTiltColor = 128

var factionColors = []color.RGBA{
	{
		R: zeroTiltColor,
		G: zeroTiltColor,
		B: 255,
		A: 255,
	},
	{
		R: 255,
		G: zeroTiltColor,
		B: zeroTiltColor,
		A: 255,
	},
	{
		R: zeroTiltColor,
		G: 255,
		B: zeroTiltColor,
		A: 255,
	},
	{
		R: 255,
		G: 255,
		B: zeroTiltColor,
		A: 255,
	},
}
