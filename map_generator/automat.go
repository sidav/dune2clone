package map_generator

type automat struct {
	drawsChar                                tileCode
	canDrawOn                                []tileCode
	maxCodeNear                              map[tileCode]int
	totalDraws, desiredTotalDraws, drawTries int
	x, y                                     int
	symmV, symmH                             bool
	radialSymmetryCount                      int
}

// true if finished
func (a *automat) perform(gm *GeneratedMap) bool {
	a.moveOnMap(gm)
	return a.totalDraws == a.desiredTotalDraws || a.drawTries == a.desiredTotalDraws*100
}

func (a *automat) moveOnMap(gm *GeneratedMap) {
	if a.totalDraws == a.desiredTotalDraws {
		return
	}
	w, h := len(gm.Tiles), len(gm.Tiles[0])
	vx, vy := 999, 999
	for vx*vy != 0 {
		vx, vy = rnd.RandomUnitVectorInt()
	}
	for a.x+vx < 0 || a.y+vy < 0 || a.x+vx >= w || a.y+vy >= h {
		vx, vy = rnd.RandomUnitVectorInt()
	}
	a.x += vx
	a.y += vy
	allowedToDrawHere := false
	for _, c := range a.canDrawOn {
		if c == gm.Tiles[a.x][a.y] {
			allowedToDrawHere = true
			for code, maxOfCode := range a.maxCodeNear {
				for gm.countTilesOfCodeNear(code, a.x, a.y) > maxOfCode {
					allowedToDrawHere = false
					break
				}
			}
		}
	}
	if allowedToDrawHere {
		gm.Tiles[a.x][a.y] = a.drawsChar
		if a.symmV && a.symmH {
			gm.Tiles[w-1-a.x][h-1-a.y] = a.drawsChar
		} else if a.symmV {
			gm.Tiles[a.x][h-1-a.y] = a.drawsChar
		} else if a.symmH {
			gm.Tiles[w-1-a.x][a.y] = a.drawsChar
		} else if a.radialSymmetryCount > 1 {
			allPoints := GetListOfCoordsRadialSymmetricTo(a.radialSymmetryCount, a.x, a.y, w, h)
			for _, coord := range allPoints {
				gm.Tiles[coord[0]][coord[1]] = a.drawsChar
			}
		}
		a.totalDraws++
	}
	a.drawTries++
}
