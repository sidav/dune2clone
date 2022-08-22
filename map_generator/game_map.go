package map_generator

import (
	"dune2clone/fibrandom"
	"dune2clone/geometry"
	"fmt"
)

var rnd fibrandom.FibRandom

func SetRandom(r *fibrandom.FibRandom) {
	rnd = *r
}

type GameMap struct {
	Tiles       [][]tileCode
	StartPoints [][2]int
}

func (gm *GameMap) init(w, h int) {
	gm.Tiles = make([][]tileCode, w)
	for i := range gm.Tiles {
		gm.Tiles[i] = make([]tileCode, h)
		for j := range gm.Tiles[i] {
			gm.Tiles[i][j] = SAND
		}
	}
	gm.StartPoints = make([][2]int, 0)
}

func (gm *GameMap) reset() {
	for i := range gm.Tiles {
		for j := range gm.Tiles[i] {
			gm.Tiles[i][j] = SAND
		}
	}
	gm.StartPoints = make([][2]int, 0)
}

func (gm *GameMap) Generate(w, h int) {
	gm.init(w, h)
	tries := 0
	for len(gm.StartPoints) == 0 || !gm.areAllStartPointsGood() {
		tries++
		gm.reset()
		symmV := rnd.OneChanceFrom(2)
		symmH := true // rnd.OneChanceFrom(2) || !symmV
		fromx, fromy, tox, toy := 0, 0, w-1, h-1
		tox = 90 * tox / 100
		toy = 90 * toy / 100
		if symmH {
			tox /= 2
		}
		if symmV {
			toy /= 2
		}

		gm.performNAutomatasLike(3,
			automat{
				drawsChar:         BUILDABLE_TERRAIN,
				canDrawOn:         []tileCode{SAND},
				desiredTotalDraws: 250,
				symmV:             symmV,
				symmH:             symmH,
			},
			fromx, fromy, tox, toy,
		)

		gm.performNAutomatasLike(20,
			automat{
				drawsChar:         POOR_RESOURCES,
				canDrawOn:         []tileCode{SAND},
				desiredTotalDraws: 25,
				symmV:             symmV,
				symmH:             symmH,
			},
			0, 0, w, h,
		)
		gm.performNAutomatasLike(20,
			automat{
				drawsChar:         MEDIUM_RESOURCES,
				canDrawOn:         []tileCode{POOR_RESOURCES},
				desiredTotalDraws: 15,
				symmV:             symmV,
				symmH:             symmH,
			},
			0, 0, w, h,
		)
		gm.performNAutomatasLike(10,
			automat{
				drawsChar:         RICH_RESOURCES,
				canDrawOn:         []tileCode{MEDIUM_RESOURCES},
				desiredTotalDraws: 5,
				symmV:             symmV,
				symmH:             symmH,
			},
			0, 0, w, h,
		)
		gm.performNAutomatasLike(5,
			automat{
				drawsChar:         RESOURCE_VEIN,
				canDrawOn:         []tileCode{RICH_RESOURCES},
				desiredTotalDraws: 1,
				symmV:             symmV,
				symmH:             symmH,
			},
			0, 0, w, h,
		)

		gm.performNAutomatasLike(10,
			automat{
				drawsChar:         ROCKS,
				canDrawOn:         []tileCode{BUILDABLE_TERRAIN},
				desiredTotalDraws: 3,
				symmV:             symmV,
				symmH:             symmH,
			},
			fromx, fromy, tox, toy,
		)

		gm.searchAndSetStartPoints(symmH, symmV, 2)
	}
	fmt.Printf("GENERATOR: Generated from %d try.\n", tries)
}

func (gm *GameMap) performNAutomatasLike(count int, prototype automat, fromx, fromy, tox, toy int) {
	autsArr := make([]automat, count)
	for i := range autsArr {
		autsArr[i] = prototype
		autsArr[i].x = rnd.RandInRange(fromx, tox)
		autsArr[i].y = rnd.RandInRange(fromy, toy)
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

func (gm *GameMap) searchAndSetStartPoints(symmH, symmV bool, count int) {
	candidates := make([][2]int, 0)
	for cx := range gm.Tiles {
		for cy := range gm.Tiles[cx] {
			// search quadrant
			if gm.areCoordsGoodForStartPoint(cx, cy) {
				candidates = append(candidates, [2]int{cx, cy})
			}
		}
	}
	if len(candidates) > 0 {
		selectedCandidate := candidates[rnd.Rand(len(candidates))]
		cx, cy := selectedCandidate[0], selectedCandidate[1]
		gm.StartPoints = append(gm.StartPoints, [2]int{cx, cy})
		if symmV && symmH {
			gm.StartPoints = append(gm.StartPoints, [2]int{len(gm.Tiles) - 1 - cx, len(gm.Tiles[0]) - 1 - cy})
		} else if symmV {
			gm.StartPoints = append(gm.StartPoints, [2]int{cx, len(gm.Tiles[0]) - 1 - cy})
		} else if symmH {
			gm.StartPoints = append(gm.StartPoints, [2]int{len(gm.Tiles) - 1 - cx, cy})
		}
		return
	}
}

func (gm *GameMap) areCoordsGoodForStartPoint(x, y int) bool {
	const sRange = 5
	for sx := x - sRange; sx <= x+sRange; sx++ {
		for sy := y - sRange; sy <= y+sRange; sy++ {
			if !(sx >= 0 && sy >= 0 && sx < len(gm.Tiles) && sy < len(gm.Tiles[0])) || gm.Tiles[sx][sy] != BUILDABLE_TERRAIN {
				return false
			}
		}
	}
	return true
}

func (gm *GameMap) areAllStartPointsGood() bool {
	w, h := len(gm.Tiles), len(gm.Tiles[0])
	minDistance := w / len(gm.StartPoints)
	if h < w {
		minDistance = h / len(gm.StartPoints)
	}
	for i := range gm.StartPoints {
		for j := range gm.StartPoints {
			if i == j {
				continue
			}
			if geometry.GetApproxDistFromTo(gm.StartPoints[i][0], gm.StartPoints[i][1], gm.StartPoints[j][0], gm.StartPoints[j][1]) < minDistance {
				return false
			}
		}
	}
	return true
}
