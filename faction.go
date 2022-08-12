package main

import (
	"image/color"
)

type faction struct {
	colorNumber             int
	resources, maxResources float64
	// money                         float64 // float because of division when spending
	currentEnergy, requiredEnergy int

	team                int // 0 means "enemy to all"
	resourcesMultiplier float64
}

func (f *faction) getMoney() float64 {
	return f.resources
}

func (f *faction) spendMoney(spent float64) {
	f.resources -= spent
}
func (f *faction) receiveResources(amount float64) {
	f.resources += amount * f.resourcesMultiplier
}

func (f *faction) resetCurrents() {
	f.maxResources = 0
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
