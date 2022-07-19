package main

import "fmt"

var (
	MAP_W = 25
	MAP_H = 14
)

const (
	DEBUG_OUTPUT = false

	SPRITE_SCALE_FACTOR     = 4.0
	ORIGINAL_TILE_SIZE_IN_PIXELS = 16
	TILE_SIZE_IN_PIXELS     = ORIGINAL_TILE_SIZE_IN_PIXELS*SPRITE_SCALE_FACTOR
	TILE_PHYSICAL_SIZE      = 25
	PIXEL_TO_PHYSICAL_RATIO = TILE_SIZE_IN_PIXELS / TILE_PHYSICAL_SIZE

	WINDOW_W = 25 * TILE_SIZE_IN_PIXELS
	WINDOW_H = 15 * TILE_SIZE_IN_PIXELS
	TEXT_SIZE = TILE_SIZE_IN_PIXELS/2
	TEXT_MARGIN = TEXT_SIZE/4
)

func halfPhysicalTileSize() int {
	if TILE_PHYSICAL_SIZE % 2 == 1 {
		return TILE_PHYSICAL_SIZE/2+1
	}
	return TILE_PHYSICAL_SIZE/2
}

func areTileCoordsValid(tx, ty int) bool {
	return tx >= 0 && tx < MAP_W && ty >= 0 && ty < MAP_H
}

func trueCoordsToTileCoords(tx, ty int) (int, int) {
	return (tx-1) / TILE_PHYSICAL_SIZE, (ty-1) / TILE_PHYSICAL_SIZE
}

func tileCoordsToPhysicalCoords(tx, ty int) (int, int) {
	halfTileSize := TILE_PHYSICAL_SIZE/2
	if TILE_PHYSICAL_SIZE % 2 == 1 {
		halfTileSize++
	}
	return tx * TILE_PHYSICAL_SIZE + halfTileSize, ty * TILE_PHYSICAL_SIZE + halfTileSize
}

func circlesOverlap(x1, y1, r1, x2, y2, r2 int) bool {
	tx := x2-x1
	ty := y2-y1
	r := r1+r2

	if tx*tx+ty*ty < r*r {
		return true
	}
	return false
}

func debugWrite(msg string) {
	if DEBUG_OUTPUT {
		fmt.Println(msg)
	}
}

func debugWritef(msg string, args... interface{}) {
	if DEBUG_OUTPUT {
		fmt.Printf(msg, args...)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
