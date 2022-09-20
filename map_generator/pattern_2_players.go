package map_generator

func (gm *GeneratedMap) generateByTwoPlayersPattern() {
	w, h := len(gm.Tiles), len(gm.Tiles[0])
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
			maxCodeNear:       map[tileCode]int{BUILDABLE_TERRAIN: 0},
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
			maxCodeNear:       map[tileCode]int{BUILDABLE_TERRAIN: 0},
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
			maxCodeNear:       map[tileCode]int{BUILDABLE_TERRAIN: 0},
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
			canDrawOn:         []tileCode{BUILDABLE_TERRAIN, SAND},
			desiredTotalDraws: 10,
			maxCodeNear:       map[tileCode]int{ROCKS: 5},
			symmV:             symmV,
			symmH:             symmH,
		},
		fromx, fromy, tox, toy,
	)

	gm.searchAndSetStartPoints(symmH, symmV, 2)
}
