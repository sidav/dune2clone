package geometry

import (
	"math"
)

var (
	degreesInCircleInt   = 360
	degreesInCircleFloat = float64(degreesInCircleInt)
)

func SetDegreesInCircleAmount(degs int) {
	degreesInCircleInt = degs
	degreesInCircleFloat = float64(degs)
}

func TrueCoordsToTileCoords(tx, ty float64) (int, int) {
	return int(tx), int(ty)
}

func TileCoordsToPhysicalCoords(tx, ty int) (float64, float64) {
	//halfTileSize := TILE_PHYSICAL_SIZE/2
	//if TILE_PHYSICAL_SIZE % 2 == 1 {
	//	halfTileSize++
	//}
	//return tx * TILE_PHYSICAL_SIZE + halfTileSize, ty * TILE_PHYSICAL_SIZE + halfTileSize
	return float64(tx) + 0.5, float64(ty) + 0.5
}

func SquareDistanceFloat64(x1, y1, x2, y2 float64) float64 {
	return (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
}

func SquareDistanceInt(x1, y1, x2, y2 int) int {
	return (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
}

func CirclesOverlap(x1, y1, r1, x2, y2, r2 int) bool {
	tx := x2 - x1
	ty := y2 - y1
	r := r1 + r2

	if tx*tx+ty*ty < r*r {
		return true
	}
	return false
}

func RotateFloat64Vector(vx, vy float64, degrees int) (float64, float64) {
	// 6.283185307179586 is 2*pi
	radians := float64(degrees) * 6.283185307179586 / degreesInCircleFloat
	t := vx
	vx = vx*math.Cos(radians) - vy*math.Sin(radians)
	vy = t*math.Sin(radians) - vy*math.Cos(radians)
	return vx, vy
}

func SnapDegreeToFixedDirections(degree, totalDirectionsInCircle int) int {
	sectorSize := degreesInCircleInt / totalDirectionsInCircle
	index := (degree + sectorSize/2) / sectorSize
	return NormalizeDegree(index * sectorSize)
}

func DegreeToSectorNumber(degree, sectorsInCircle int) int {
	sectorSize := degreesInCircleInt / sectorsInCircle
	return (degree + sectorSize/2) / sectorSize
}

func GetDiffForRotationStep(currDegree, targetDegree, rotateSpeed int) int {
	if targetDegree == currDegree {
		return 0
	}
	if targetDegree < 0 {
		targetDegree += degreesInCircleInt
	}
	diff := currDegree - targetDegree
	for diff < 0 {
		diff += degreesInCircleInt
	}
	if rotateSpeed > diff {
		rotateSpeed = diff
	} else if rotateSpeed > degreesInCircleInt-diff {
		rotateSpeed = degreesInCircleInt - diff
	}
	if diff <= degreesInCircleInt/2 {
		rotateSpeed = -rotateSpeed
	}
	return rotateSpeed
}

func NormalizeDegree(deg int) int {
	for deg < 0 {
		deg += degreesInCircleInt
	}
	for deg >= degreesInCircleInt {
		deg -= degreesInCircleInt
	}
	return deg
}

func IsVectorDegreeEqualTo(vx, vy float64, deg int) bool {
	vectorDegree := int((degreesInCircleFloat / 2) * math.Atan2(vy, vx) / 3.14159265358)
	for vectorDegree < 0 {
		vectorDegree += degreesInCircleInt
	}
	return deg == vectorDegree
}

func AreRectsInRange(x1, y1, w1, h1, x2, y2, w2, h2, r int) bool {
	// all -1's are beacuse of TILED geometry
	x1b := x1 + w1 - 1
	x2b := x2 + w2 - 1
	y1b := y1 + h1 - 1
	y2b := y2 + h2 - 1

	left := x2b < x1
	right := x1b < x2
	bottom := y1b < y2
	top := y2b < y1
	if top && left {
		return AreCoordsInRange(x1, y1, x2b, y2b, r) // dist((x1, y1b), (x2b, y2))
	}
	if left && bottom {
		return AreCoordsInRange(x1, y1b, x2b, y2, r)
	}
	if bottom && right {
		return AreCoordsInRange(x1b, y1b, x2, y2, r)
	}
	if right && top {
		return AreCoordsInRange(x1b, y1, x2, y2b, r)
	}
	if left {
		return x1-x2b <= r
	}
	if right {
		return x2-x1b <= r
	}
	if bottom {
		return y2-y1b <= r
	}
	if top {
		return y1-y2b <= r
	}
	return true // intersect detected
}

func AreCoordsInRangeFromRect(fx, fy, tx, ty, w, h, r int) bool { // considering ANY of the tiles in the rect.
	return AreRectsInRange(fx, fy, 1, 1, tx, ty, w, h, r)
}

func GetDegreeOfFloatVector(vx, vy float64) int {
	return NormalizeDegree(int((degreesInCircleFloat / 2) * math.Atan2(vy, vx) / 3.14159265358))
}

func GetDegreeOfIntVector(vx, vy int) int {
	return GetDegreeOfFloatVector(float64(vy), float64(vx))
}

func AreTwoCellRectsOverlapping(x1, y1, w1, h1, x2, y2, w2, h2 int) bool {
	// WARNING:
	// ALL "-1"s HERE ARE BECAUSE OF WE ARE IN CELLS SPACE
	// I.E. A SINGLE CELL IS 1x1 RECTANGLE
	// SO RECTS (0, 0, 1x1) AND (1, 0, 1x1) ARE NOT OVERLAPPING IN THIS SPACE (BUT SHOULD IN EUCLIDEAN OF COURSE)
	right1 := x1 + w1 - 1
	bot1 := y1 + h1 - 1
	right2 := x2 + w2 - 1
	bot2 := y2 + h2 - 1
	return !(x2 > right1 ||
		right2 < x1 ||
		y2 > bot1 ||
		bot2 < y1)
}

func AreTwoCellRectsOverlapping32(x1, y1, w1, h1, x2, y2, w2, h2 int32) bool {
	// WARNING:
	// ALL "-1"s HERE ARE BECAUSE OF WE ARE IN CELLS SPACE
	// I.E. A SINGLE CELL IS 1x1 RECTANGLE
	// SO RECTS (0, 0, 1x1) AND (1, 0, 1x1) ARE NOT OVERLAPPING IN THIS SPACE (BUT SHOULD IN EUCLIDEAN OF COURSE)
	right1 := x1 + w1 - 1
	bot1 := y1 + h1 - 1
	right2 := x2 + w2 - 1
	bot2 := y2 + h2 - 1
	return !(x2 > right1 ||
		right2 < x1 ||
		y2 > bot1 ||
		bot2 < y1)
}

func DegreeToUnitVector(deg int) (float64, float64) {
	return math.Cos(float64(deg) * 3.14159265358 / (degreesInCircleFloat / 2)), math.Sin(float64(deg) * 3.14159265358 / (degreesInCircleFloat / 2))
}

func VectorToUnitVectorFloat64(vx, vy float64) (float64, float64) {
	length := math.Sqrt(vx*vx + vy*vy)
	if vx != 0 {
		vx /= length
	}
	if vy != 0 {
		vy /= length
	}
	return vx, vy
}

func Float64VectorToIntDirectionVector(vx, vy float64) (int, int) {
	intVx, intVy := 0, 0
	if vx < 0 {
		intVx = int(vx - 0.5)
	}
	if vx > 0 {
		intVx = int(vx + 0.5)
	}
	if vy < 0 {
		intVy = int(vy - 0.5)
	}
	if vy > 0 {
		intVy = int(vy + 0.5)
	}
	return intVx, intVy
}

func AreCoordsInTileRect(x, y, rx, ry, w, h int) bool {
	return x >= rx && x < rx+w && y >= ry && y < ry+h
}

func AreCoordsInRange(fx, fy, tx, ty, r int) bool { // border including.
	// uses more wide circle (like in Bresenham's circle) than the real geometric one.
	// It is much more handy for spaces with discrete coords (cells).
	realSqDistanceAndSqRadiusDiff := (fx-tx)*(fx-tx) + (fy-ty)*(fy-ty) - r*r
	return realSqDistanceAndSqRadiusDiff < r
}

func GetApproxDistFromTo(x1, y1, x2, y2 int) int {
	diffX := x1 - x2
	if diffX < 0 {
		diffX = -diffX
	}
	diffY := y1 - y2
	if diffY < 0 {
		diffY = -diffY
	}
	if diffX > diffY {
		return diffX + (diffY >> 1)
	} else {
		return diffY + (diffX >> 1)
	}
}

func GetApproxDistFloat64(x1, y1, x2, y2 float64) float64 {
	diffX := x1 - x2
	if diffX < 0 {
		diffX = -diffX
	}
	diffY := y1 - y2
	if diffY < 0 {
		diffY = -diffY
	}
	if diffX > diffY {
		return diffX + (diffY / 2)
	} else {
		return diffY + (diffX / 2)
	}
}

func SpiralSearchForConditionFrom(condition func(int, int) bool, startX, startY, maxSearchRadius, initialDirection int) (int, int) {
	currRadius := 1
	// direction 0
	var vx, vy, x, y int
	switch initialDirection % 4 {
	case 0:
		vx, vy = 1, 0
		x, y = startX, startY-currRadius
	case 1:
		vx, vy = 0, 1
		x, y = startX+currRadius, startY
	case 2:
		vx, vy = -1, 0
		x, y = startX, startY+currRadius
	case 3:
		vx, vy = 0, -1
		x, y = startX-currRadius, startY
	}
	currStartX, currStartY := x, y
	for {
		if condition(x, y) {
			return x, y
		}
		x += vx
		y += vy
		// rotate if we are at corner of current square
		if x == startX+currRadius && y == startY-currRadius || // right top
			x == startX+currRadius && y == startY+currRadius || // right bottom
			x == startX-currRadius && y == startY+currRadius || // left bottom
			x == startX-currRadius && y == startY-currRadius {

			t := vx
			vx = -vy
			vy = t
		}
		// increasing radius
		if x == currStartX && y == currStartY {
			currRadius++
			if currRadius > maxSearchRadius {
				return -1, -1
			}
			switch initialDirection % 4 {
			case 0:
				vx, vy = 1, 0
				x, y = startX, startY-currRadius
			case 1:
				vx, vy = 0, 1
				x, y = startX+currRadius, startY
			case 2:
				vx, vy = -1, 0
				x, y = startX, startY+currRadius
			case 3:
				vx, vy = 0, -1
				x, y = startX-currRadius, startY
			}
			currStartX, currStartY = x, y
		}
	}
	// return -1, -1
}
