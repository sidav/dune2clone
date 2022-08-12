package map_generator

import (
	"dune2clone/fibrandom"
)

var rnd fibrandom.FibRandom

func SetRandom(r *fibrandom.FibRandom) {
	rnd = *r
}

type GameMap struct {
	Tiles       [][]int
	StartPoints [][2]int
}

func (gm *GameMap) Init(w, h int) {
	gm.Tiles = make([][]int, w)
	for i := range gm.Tiles {
		gm.Tiles[i] = make([]int, h)
		for j := range gm.Tiles[i] {
			gm.Tiles[i][j] = SAND
		}
	}
	gm.StartPoints = make([][2]int, 0)
}

func (gm *GameMap) Generate() {
	for len(gm.StartPoints) == 0 {
		w, h := len(gm.Tiles), len(gm.Tiles[0])
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
				canDrawOn:         []int{SAND},
				desiredTotalDraws: 250,
				symmV:             symmV,
				symmH:             symmH,
			},
			fromx, fromy, tox, toy,
		)

		gm.performNAutomatasLike(20,
			automat{
				drawsChar:         RESOURCES,
				canDrawOn:         []int{SAND},
				desiredTotalDraws: 25,
				symmV:             symmV,
				symmH:             symmH,
			},
			fromx, fromy, tox, toy,
		)

		gm.performNAutomatasLike(10,
			automat{
				drawsChar:         ROCKS,
				canDrawOn:         []int{BUILDABLE_TERRAIN},
				desiredTotalDraws: 3,
				symmV:             symmV,
				symmH:             symmH,
			},
			fromx, fromy, tox, toy,
		)

		gm.searchAndSetStartPoints(symmH, symmV, 2)
	}
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
	// check if not too close to center (4/8 of the map to 6/8)
	w, h := len(gm.Tiles), len(gm.Tiles[0])
	if x > 3*w/8 && x < 6*w/8 || y > 3*h/8 && y < 6*h/8 {
		return false
	}
	const sRange = 3
	for sx := x - sRange; sx <= x+sRange; sx++ {
		for sy := y - sRange; sy <= y+sRange; sy++ {
			if !(sx >= 0 && sy >= 0 && sx < len(gm.Tiles) && sy < len(gm.Tiles[0])) || gm.Tiles[sx][sy] != BUILDABLE_TERRAIN {
				return false
			}
		}
	}
	return true
}
