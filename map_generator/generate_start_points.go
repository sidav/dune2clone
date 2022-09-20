package map_generator

import "dune2clone/geometry"

func (gm *GeneratedMap) searchAndSetStartPoints(symmH, symmV bool, count int) {
	// TODO: rewrite this completely
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
		if symmV && symmH {
			gm.StartPoints = append(gm.StartPoints, [2]int{cx, cy})
			gm.StartPoints = append(gm.StartPoints, [2]int{len(gm.Tiles) - 1 - cx, len(gm.Tiles[0]) - 1 - cy})
		} else if symmV {
			gm.StartPoints = append(gm.StartPoints, [2]int{cx, cy})
			gm.StartPoints = append(gm.StartPoints, [2]int{cx, len(gm.Tiles[0]) - 1 - cy})
		} else if symmH {
			gm.StartPoints = append(gm.StartPoints, [2]int{cx, cy})
			gm.StartPoints = append(gm.StartPoints, [2]int{len(gm.Tiles) - 1 - cx, cy})
		} else if count % 2 != 0 {
			allPoints := GetListOfCoordsRadialSymmetricTo(count, cx, cy, len(gm.Tiles), len(gm.Tiles[0]))
			for _, coord := range allPoints {
				gm.StartPoints = append(gm.StartPoints, coord)
			}
		}
		return
	}
}

func (gm *GeneratedMap) areCoordsGoodForStartPoint(x, y int) bool {
	const sRange = 5
	for sx := x - sRange; sx <= x+sRange; sx++ {
		for sy := y - sRange; sy <= y+sRange; sy++ {
			if !(sx > 1 && sy > 1 && sx < len(gm.Tiles)-2 && sy < len(gm.Tiles[0])-2) || gm.Tiles[sx][sy] != BUILDABLE_TERRAIN {
				return false
			}
		}
	}
	return true
}

func (gm *GeneratedMap) areAllStartPointsGood() bool {
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

