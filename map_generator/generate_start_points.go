package map_generator

import (
	"dune2clone/geometry"
)

func (gm *GeneratedMap) searchAndSetStartPoints(symmH, symmV bool, radialSymmetryCount int) {
	// TODO: rewrite this completely
	candidates := make([][2]int, 0)
	for cx := range gm.Tiles {
		for cy := range gm.Tiles[cx] {
			// search quadrant
			if gm.areCoordsGoodForStartPoint(cx, cy, 5) {
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
		} else if radialSymmetryCount > 1 {
			allPoints := GetListOfCoordsRadialSymmetricTo(radialSymmetryCount, cx, cy, len(gm.Tiles), len(gm.Tiles[0]))
			for _, coord := range allPoints {
				gm.StartPoints = append(gm.StartPoints, coord)
			}
		}
		return
	}
}

func (gm *GeneratedMap) searchAndSetStartPointsAsymmetric(count int) {
	// TODO: rewrite this completely
	candidates := make([][2]int, 0)
	for cx := range gm.Tiles {
		for cy := range gm.Tiles[cx] {
			// search quadrant
			if gm.areCoordsGoodForStartPoint(cx, cy, 4) {
				candidates = append(candidates, [2]int{cx, cy})
			}
		}
	}
	if len(candidates) < count*count*count {
		return
	}
	var selected [][2]int

	for try := 0; try < len(candidates)*len(candidates); try++ {
		isGood := false
		currCand := [2]int{0, 0}
		indexTry := 0
		for !isGood {
			isGood = true
			rndInd := rnd.Rand(len(candidates))
			currCand[0] = candidates[rndInd][0]
			currCand[1] = candidates[rndInd][1]
			for i := 0; i < len(selected); i++ {
				if geometry.GetApproxDistFromTo(currCand[0], currCand[1], selected[i][0], selected[i][1]) < 10 {
					isGood = false
					break
				}
			}
			if isGood {
				selected = append(selected, currCand)
			}
			indexTry++
			if indexTry > len(candidates) {
				break
			}
		}
		if len(selected) == count {
			break
		}
	}
	if len(selected) < count {
		return
	}
	for i := 0; i < len(selected) && i < count; i++ {
		gm.StartPoints = append(gm.StartPoints, selected[i])
	}
}

func (gm *GeneratedMap) areCoordsGoodForStartPoint(x, y, buildableRange int) bool {
	for sx := x - buildableRange; sx <= x+buildableRange; sx++ {
		for sy := y - buildableRange; sy <= y+buildableRange; sy++ {
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
