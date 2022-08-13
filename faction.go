package main

import (
	"image/color"
	"math"
)

type faction struct {
	colorNumber                       int
	currentResources, resourceStorage float64
	money                             float64 // float because of division when spending
	currentEnergy, requiredEnergy     int

	team                int // 0 means "enemy to all"
	resourcesMultiplier float64
}

func (f *faction) getMoney() float64 {
	return f.currentResources + f.money
}

func (f *faction) getStorageRemaining() float64 {
	return f.resourceStorage - f.currentResources
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
	f.currentEnergy = 0
	f.requiredEnergy = 0
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
