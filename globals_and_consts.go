package main

import (
	"fmt"
	"math"
)

var (
	MAP_W = 64
	MAP_H = 64
)

const (
	DEBUG_OUTPUT = true

	DESIRED_FPS                   = 60
	UNIT_ACTIONS_TICK_EACH        = 2
	PROJECTILES_ACTIONS_TICK_EACH = 2
	BUILDINGS_ACTIONS_TICK_EACH   = 5

	AI_ACTS_EACH = 60

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

func debugWrite(msg string) {
	if DEBUG_OUTPUT {
		fmt.Println(msg)
	}
}

func debugWritef(msg string, args ...interface{}) {
	if DEBUG_OUTPUT {
		fmt.Printf(msg, args...)
	}
}

func areTileCoordsValid(tx, ty int) bool {
	return tx >= 0 && tx < MAP_W && ty >= 0 && ty < MAP_H
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
