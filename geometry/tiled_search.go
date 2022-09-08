package geometry

func SpiralSearchForClosestConditionFrom(condition func(int, int) bool, startX, startY, maxSearchRadius, initialDirection int) (int, int) {
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

func SpiralSearchForFarthestConditionFrom(condition func(int, int) bool, startX, startY, maxSearchRadius, initialDirection int) (int, int) {
	currRadius := 1
	currFoundX, currFoundY := -1, -1
	currMaxRadius := 0
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
		if condition(x, y) && currMaxRadius < currRadius {
			currFoundX, currFoundY = x, y
			currMaxRadius = currRadius
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
				return currFoundX, currFoundY
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

func SpiralSearchForHighestScoreFrom(score func(int, int) int, startX, startY, maxSearchRadius, initialDirection int) (int, int) {
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
	somethingFound := false
	currFoundX, currFoundY, currMaxScore := -1, -1, 0
	for {
		currScore := score(x, y)
		if currScore != -1 && (!somethingFound || currMaxScore < currScore) {
			somethingFound = true
			currFoundX, currFoundY = x, y
			currMaxScore = currScore
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
				return currFoundX, currFoundY
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
