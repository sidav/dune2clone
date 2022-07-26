package main

import (
	"image/color"
)

type faction struct {
	factionColor color.RGBA
	money        float64 // float because of division when spending
	energy       int
	team         int // 0 means "enemy to all"
}

const zeroTiltColor = 48
const strongerTiltColor = 128

var factionTints = []color.RGBA{
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
