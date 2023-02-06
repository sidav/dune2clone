package main

import (
	"fmt"
	"math"
)

const (
	MAX_VETERANCY_LEVEL = 4

	SPRITE_SCALE_FACTOR          = 4
	ORIGINAL_TILE_SIZE_IN_PIXELS = 16
	TILE_SIZE_IN_PIXELS          = ORIGINAL_TILE_SIZE_IN_PIXELS * SPRITE_SCALE_FACTOR
	TILE_PHYSICAL_SIZE           = 1 // TODO: remove, since we're using floats?
	PIXEL_TO_PHYSICAL_RATIO      = TILE_SIZE_IN_PIXELS / TILE_PHYSICAL_SIZE

	TEXT_SIZE   = TILE_SIZE_IN_PIXELS / 2
	TEXT_MARGIN = TEXT_SIZE / 4
)

var (
	WINDOW_W = int32(25 * TILE_SIZE_IN_PIXELS)
	WINDOW_H = int32(15 * TILE_SIZE_IN_PIXELS)
)

func halfPhysicalTileSize() int {
	if TILE_PHYSICAL_SIZE%2 == 1 {
		return TILE_PHYSICAL_SIZE/2 + 1
	}
	return TILE_PHYSICAL_SIZE / 2
}

// trying to overcome rounding issues
func areFloatsAlmostEqual(f, g float64) bool {
	return math.Abs(f-g) < 0.0001
}

func getPercentInt(curr, max int) int {
	return 100 * curr / max
}

func areFloatsRoughlyEqual(f, g float64) bool {
	return math.Abs(f-g) < 0.01
}

func sign(x int) int {
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	return 0
}

func getListOfRelativeCoordsForSquadMembers(squadSize int) [][2]float64 {
	// returns list of coords, relative to center, (for example, for drawing a squad of units)
	// consider that the coords won't rotate with squad
	switch squadSize {
	case 0, 1:
		return [][2]float64{{0, 0}}
	case 2:
		return [][2]float64{{0.3, -0.3}, {-0.3, 0.3}}
	case 3:
		return [][2]float64{{0.25, 0.25}, {0, -0.20}, {-0.25, 0.25}}
	case 4:
		return [][2]float64{{0, -0.32}, {0.32, 0}, {0, 0.32}, {-0.32, 0}}
	case 5:
		return [][2]float64{{0, -0.32}, {0.32, 0}, {0, 0.32}, {-0.32, 0}, {0, 0}}
	}
	return [][2]float64{}
	panic(fmt.Sprintf("No such squad size %d, renderer failed", squadSize))
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
