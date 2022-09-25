package map_generator

func generateByTwoPlayersPattern(gm *GeneratedMap) {
	w, h := len(gm.Tiles), len(gm.Tiles[0])
	gm.reset()
	symmV := false // rnd.OneChanceFrom(2)
	symmH := false // true // rnd.OneChanceFrom(2) || !symmV
	fromx, fromy, tox, toy := 0, 0, w-1, h-1

	gm.performNAutomatasLike(3,
		rnd.RandInRange(20, 60),
		automat{
			drawsChar:         BUILDABLE_TERRAIN,
			canDrawOn:         []tileCode{SAND},
			symmV:             symmV,
			symmH:             symmH,
			radialSymmetryCount: 2,
		},
		fromx, fromy, tox, toy,
	)

	gm.performNAutomatasLike(20,
		rnd.RandInRange(10, 15),
		automat{
			drawsChar:         POOR_RESOURCES,
			canDrawOn:         []tileCode{SAND},
			maxCodeNear:       map[tileCode]int{BUILDABLE_TERRAIN: 0},
			symmV:             symmV,
			symmH:             symmH,
			radialSymmetryCount: 2,
		},
		0, 0, w, h,
	)

	gm.performNAutomatasLike(20,
		rnd.RandInRange(5, 10),
		automat{
			drawsChar:         MEDIUM_RESOURCES,
			canDrawOn:         []tileCode{POOR_RESOURCES},
			maxCodeNear:       map[tileCode]int{BUILDABLE_TERRAIN: 0},
			symmV:             symmV,
			symmH:             symmH,
			radialSymmetryCount: 2,
		},
		0, 0, w, h,
	)

	gm.performNAutomatasLike(10,
		rnd.RandInRange(2, 5),
		automat{
			drawsChar:         RICH_RESOURCES,
			canDrawOn:         []tileCode{MEDIUM_RESOURCES},
			maxCodeNear:       map[tileCode]int{BUILDABLE_TERRAIN: 0},
			symmV:             symmV,
			symmH:             symmH,
			radialSymmetryCount: 2,
		},
		0, 0, w, h,
	)

	gm.performNAutomatasLike(5,
		rnd.RandInRange(1, 2),
		automat{
			drawsChar:         RESOURCE_VEIN,
			canDrawOn:         []tileCode{RICH_RESOURCES},
			symmV:             symmV,
			symmH:             symmH,
			radialSymmetryCount: 2,
		},
		0, 0, w, h,
	)

	gm.performNAutomatasLike(10,
		rnd.RandInRange(5, 25),
		automat{
			drawsChar:         ROCKS,
			canDrawOn:         []tileCode{BUILDABLE_TERRAIN, SAND},
			maxCodeNear:       map[tileCode]int{ROCKS: 5},
			symmV:             symmV,
			symmH:             symmH,
			radialSymmetryCount: 2,
		},
		fromx, fromy, tox, toy,
	)

	gm.searchAndSetStartPoints(symmH, symmV, 2)
}
