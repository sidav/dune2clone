package main

import (
	"fmt"
	"math"
)

var (
	MAP_W = 20
	MAP_H = 14
)

const (
	DEBUG_OUTPUT = true

	DESIRED_FPS                 = 60
	UNIT_ACTIONS_TICK_EACH      = 2
	BUILDINGS_ACTIONS_TICK_EACH = 5

	SPRITE_SCALE_FACTOR          = 4.0
	ORIGINAL_TILE_SIZE_IN_PIXELS = 16
	TILE_SIZE_IN_PIXELS          = ORIGINAL_TILE_SIZE_IN_PIXELS * SPRITE_SCALE_FACTOR
	TILE_PHYSICAL_SIZE           = 1 // TODO: remove, since we're using floats?
	PIXEL_TO_PHYSICAL_RATIO      = TILE_SIZE_IN_PIXELS / TILE_PHYSICAL_SIZE

	WINDOW_W    = 25 * TILE_SIZE_IN_PIXELS
	WINDOW_H    = 15 * TILE_SIZE_IN_PIXELS
	TEXT_SIZE   = TILE_SIZE_IN_PIXELS / 2
	TEXT_MARGIN = TEXT_SIZE / 4
)

func halfPhysicalTileSize() int {
	if TILE_PHYSICAL_SIZE%2 == 1 {
		return TILE_PHYSICAL_SIZE/2 + 1
	}
	return TILE_PHYSICAL_SIZE / 2
}

func areTileCoordsValid(tx, ty int) bool {
	return tx >= 0 && tx < MAP_W && ty >= 0 && ty < MAP_H
}

func trueCoordsToTileCoords(tx, ty float64) (int, int) {
	return int(tx), int(ty)
}

func tileCoordsToPhysicalCoords(tx, ty int) (float64, float64) {
	//halfTileSize := TILE_PHYSICAL_SIZE/2
	//if TILE_PHYSICAL_SIZE % 2 == 1 {
	//	halfTileSize++
	//}
	//return tx * TILE_PHYSICAL_SIZE + halfTileSize, ty * TILE_PHYSICAL_SIZE + halfTileSize
	return float64(tx) + 0.5, float64(ty) + 0.5
}

func circlesOverlap(x1, y1, r1, x2, y2, r2 int) bool {
	tx := x2 - x1
	ty := y2 - y1
	r := r1 + r2

	if tx*tx+ty*ty < r*r {
		return true
	}
	return false
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

func degreeToRotationFrameNumber(degree, sectorsInCircle int) int {
	sectorWidth := 360 / sectorsInCircle
	degree += sectorWidth / 2
	for degree < 0 {
		degree += 360
	}
	for degree >= 360 {
		degree -= 360
	}
	num := 0
	for degree >= sectorWidth {
		degree -= sectorWidth
		num++
	}
	// +1 is because initial images look up (last frame number)
	return (num + 90/sectorWidth) % (360 / sectorWidth)
}

func getDiffForRotationStep(currDegree, targetDegree, rotateSpeed int) int {
	if targetDegree == currDegree {
		return 0
	}
	if targetDegree < 0 {
		targetDegree += 360
	}
	diff := currDegree - targetDegree
	for diff < 0 {
		diff += 360
	}
	if rotateSpeed > diff {
		rotateSpeed = diff
	} else if rotateSpeed > 360-diff {
		rotateSpeed = 360-diff
	}
	if diff <= 180 {
		rotateSpeed = -rotateSpeed
	}
	return rotateSpeed
}

func normalizeDegree(deg int) int {
	for deg < 0 {
		deg += 360
	}
	for deg >= 360 {
		deg -= 360
	}
	return deg
}

func isVectorDegreeEqualTo(vx, vy float64, deg int) bool {
	vectorDegree := int(180 * math.Atan2(vy, vx) / 3.14159265358)
	for vectorDegree < 0 {
		vectorDegree += 360
	}
	return deg == vectorDegree
}

func getDegreeOfIntVector(vx, vy int) int {
	return int(180 * math.Atan2(float64(vy), float64(vx)) / 3.14159265358)
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

func areCoordsInTileRect(x, y, rx, ry, w, h int) bool {
	return x >= rx && x < rx+w && y >= ry && y < ry+h
}

func areCoordsInRange(fx, fy, tx, ty, r int) bool { // border including.
	// uses more wide circle (like in Bresenham's circle) than the real geometric one.
	// It is much more handy for spaces with discrete coords (cells).
	realSqDistanceAndSqRadiusDiff := (fx-tx)*(fx-tx) + (fy-ty)*(fy-ty) - r*r
	return realSqDistanceAndSqRadiusDiff < r
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
