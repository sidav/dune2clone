package map_generator

func generateByTwoPlayersPattern(gm *GeneratedMap) {
	w, h := len(gm.Tiles), len(gm.Tiles[0])
	gm.reset()
	fromx, fromy, tox, toy := 0, 0, w-1, h-1

	gm.performNAutomatasLike(3,
		rnd.RandInRange(20, 60),
		0,
		automat{
			drawsChar:           BUILDABLE_TERRAIN,
			canDrawOn:           []tileCode{SAND},
			radialSymmetryCount: 2,
		},
		fromx, fromy, tox, toy,
	)

	gm.performNAutomatasLike(20,
		rnd.RandInRange(1, 15),
		0,
		automat{
			drawsChar:           ROCKS,
			canDrawOn:           []tileCode{BUILDABLE_TERRAIN, SAND},
			maxCodeNear:         map[tileCode]int{ROCKS: 5},
			radialSymmetryCount: 2,
		},
		fromx, fromy, tox, toy,
	)

	gm.performNAutomatasLike(20,
		rnd.RandInRange(10, 15),
		0,
		automat{
			drawsChar:           POOR_RESOURCES,
			canDrawOn:           []tileCode{SAND},
			maxCodeNear:         map[tileCode]int{BUILDABLE_TERRAIN: 0, ROCKS: 0},
			radialSymmetryCount: 2,
		},
		0, 0, w, h,
	)

	gm.performNAutomatasLike(20,
		rnd.RandInRange(5, 10),
		0,
		automat{
			drawsChar:           MEDIUM_RESOURCES,
			canDrawOn:           []tileCode{POOR_RESOURCES},
			maxCodeNear:         map[tileCode]int{BUILDABLE_TERRAIN: 0},
			radialSymmetryCount: 2,
		},
		0, 0, w, h,
	)

	gm.performNAutomatasLike(10,
		rnd.RandInRange(2, 5),
		0,
		automat{
			drawsChar:           RICH_RESOURCES,
			canDrawOn:           []tileCode{MEDIUM_RESOURCES},
			maxCodeNear:         map[tileCode]int{BUILDABLE_TERRAIN: 0},
			radialSymmetryCount: 2,
		},
		0, 0, w, h,
	)

	gm.performNAutomatasLike(2*rnd.RandInRange(1, 4),
		0,
		1,
		automat{
			drawsChar:           RESOURCE_VEIN,
			canDrawOn:           []tileCode{RICH_RESOURCES},
			radialSymmetryCount: 2,
		},
		0, 0, w, h,
	)

	// gm.cleanupBadRadialSymmetry(2)
	gm.searchAndSetStartPoints(false, false, 2)
}
