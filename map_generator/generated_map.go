package map_generator

import (
	"dune2clone/fibrandom"
	"fmt"
	"math"
)

var rnd fibrandom.FibRandom

func SetRandom(r *fibrandom.FibRandom) {
	rnd = *r
}

type GeneratedMap struct {
	Tiles       [][]tileCode
	StartPoints [][2]int
}

func (gm *GeneratedMap) init(w, h int) {
	gm.Tiles = make([][]tileCode, w)
	for i := range gm.Tiles {
		gm.Tiles[i] = make([]tileCode, h)
		for j := range gm.Tiles[i] {
			gm.Tiles[i][j] = SAND
		}
	}
	gm.StartPoints = make([][2]int, 0)
}

func (gm *GeneratedMap) areCoordsCorrect(x, y int) bool {
	return x > 0 && y > 0 && x < len(gm.Tiles) && y < len(gm.Tiles[x])
}

func (gm *GeneratedMap) reset() {
	for i := range gm.Tiles {
		for j := range gm.Tiles[i] {
			gm.Tiles[i][j] = SAND
		}
	}
	gm.StartPoints = make([][2]int, 0)
}

func (gm *GeneratedMap) Generate(w, h, patternIndex int) {
	gm.init(w, h)
	tries := 0
	for len(gm.StartPoints) == 0 || !gm.areAllStartPointsGood() {
		tries++
		GetPatternByIndex(patternIndex).generationFunc(gm)
	}
	fmt.Printf("GENERATOR: Generated from %d try.\n", tries)
}


func (gm *GeneratedMap) getNumberOfTilesPercent(perc int) int {
	total := len(gm.Tiles) * len(gm.Tiles[0])
	fmt.Printf("w, h %d,%d total %d; %d percentage, got %d\n",len(gm.Tiles), len(gm.Tiles[0]), total, perc, perc*total/100)
	return perc * total / 100
}

func (gm *GeneratedMap) performNAutomatasLike(count int, coveragePercent, forceDrawsEach int, prototype automat, fromx, fromy, tox, toy int) {
	autsArr := make([]automat, count)
	totalDrawsEach := 0
	if coveragePercent > 0 {
		totalDrawsEach = gm.getNumberOfTilesPercent(coveragePercent)
		totalDrawsEach /= count
		if prototype.radialSymmetryCount > 0 {
			totalDrawsEach /= prototype.radialSymmetryCount
		}
	}
	if forceDrawsEach > 0 {
		totalDrawsEach = forceDrawsEach
	}
	for i := range autsArr {
		autsArr[i] = prototype
		autsArr[i].desiredTotalDraws = totalDrawsEach
		autsArr[i].x = rnd.RandInRange(fromx, tox)
		autsArr[i].y = rnd.RandInRange(fromy, toy)
		//for _, restrictedCode := range prototype.cantDrawNear {
		//	for gm.countTilesOfCodeNear(restrictedCode, autsArr[i].x, autsArr[i].y) > 0 {
		//		autsArr[i].x = rnd.RandInRange(fromx, tox)
		//		autsArr[i].y = rnd.RandInRange(fromy, toy)
		//	}
		//}
	}
	finished := false
	for !finished {
		finished = true
		for i := range autsArr {
			if !autsArr[i].perform(gm) {
				finished = false
			}
			//draw(gm)
			//drawAt(autsArr[i].x, autsArr[i].y)
		}
	}
}

func (gm *GeneratedMap) cleanupBadRadialSymmetry(times int) { // this is a workaround
	for i := 0; i < times; i++ {
		for x := range gm.Tiles {
			for y := range gm.Tiles[x] {
				if gm.Tiles[x][y] == SAND {
					if gm.countTilesOfCodeNear(BUILDABLE_TERRAIN, x, y) >= 7 {
						gm.Tiles[x][y] = BUILDABLE_TERRAIN
					}
				}
			}
		}
	}
}

func (gm *GeneratedMap) countTilesOfCodeNear(code tileCode, x, y int) int {
	count := 0
	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if gm.areCoordsCorrect(i, j) && gm.Tiles[i][j] == code {
				count++
			}
		}
	}
	return count
}

func GetListOfCoordsRadialSymmetricTo(count, x, y, mapW, mapH int) [][2]int {
	if count < 2 {
		panic("Bad count")
	}
	degreesBetweenCoords := 2*math.Pi / float64(count)
	centerFloatX, centerFloatY := float64(mapW-1)/2, float64(mapH-1)/2
	coords := make([][2]int, count)
	vectorX, vectorY := float64(x)-centerFloatX, float64(y)-centerFloatY
	//if count == 2 {
	//	return [][2]int {
	//		{x, y},
	//		{},
	//	}
	//}
	for i := 0; i < count; i++ {
		currTileX := int(math.Round(vectorX+centerFloatX))
		currTileY := int(math.Round(vectorY+centerFloatY))

		if currTileX < 0 {
			currTileX = 0
		}
		if currTileX >= mapW {
			currTileX = mapW-1
		}
		if currTileY < 0 {
			currTileY = 0
		}
		if currTileY >= mapH {
			currTileY = mapH-1
		}
		coords[i][0] = currTileX
		coords[i][1] = currTileY
		// rotate vector
		t := vectorX
		vectorX = vectorX * math.Cos(degreesBetweenCoords) - vectorY * math.Sin(degreesBetweenCoords)
		vectorY = t * math.Sin(degreesBetweenCoords) + vectorY * math.Cos(degreesBetweenCoords)
	}
	for _, c := range coords {
		if c[0] < 0 || c[0] >= mapW || c[1] < 0 || c[1] >= mapH {
			fmt.Printf("%v crashed\n", coords)
			break
		}
	}
	return coords
}
