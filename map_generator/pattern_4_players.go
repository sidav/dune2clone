package map_generator

func generateByFourPlayersPattern(gm *GeneratedMap) {
	w, h := len(gm.Tiles), len(gm.Tiles[0])
	gm.reset()
	fromx, fromy, tox, toy := w/10, h/10, w-1, h-1
	tox = 90 * tox / 100
	toy = 90 * toy / 100
	var automs int

	automs = 3
	gm.performNAutomatasLike(automs,
		rnd.RandInRange(20, 60),
		automat{
			drawsChar:           BUILDABLE_TERRAIN,
			canDrawOn:           []tileCode{SAND},
			radialSymmetryCount: 4,
		},
		fromx, fromy, tox, toy,
	)

	gm.performNAutomatasLike(20,
		rnd.RandInRange(10, 25),
		automat{
			drawsChar:           POOR_RESOURCES,
			canDrawOn:           []tileCode{SAND},
			maxCodeNear:         map[tileCode]int{BUILDABLE_TERRAIN: 0},
			radialSymmetryCount: 4,
		},
		0, 0, w, h,
	)
	gm.performNAutomatasLike(20,
		rnd.RandInRange(10, 15),
		automat{
			drawsChar:           MEDIUM_RESOURCES,
			canDrawOn:           []tileCode{POOR_RESOURCES},
			maxCodeNear:         map[tileCode]int{BUILDABLE_TERRAIN: 0},
			radialSymmetryCount: 4,
		},
		0, 0, w, h,
	)
	gm.performNAutomatasLike(10,
		rnd.RandInRange(5, 10),
		automat{
			drawsChar:           RICH_RESOURCES,
			canDrawOn:           []tileCode{MEDIUM_RESOURCES},
			maxCodeNear:         map[tileCode]int{BUILDABLE_TERRAIN: 0},
			radialSymmetryCount: 4,
		},
		0, 0, w, h,
	)
	gm.performNAutomatasLike(3,
		rnd.RandInRange(1, 2),
		automat{
			drawsChar:           RESOURCE_VEIN,
			canDrawOn:           []tileCode{POOR_RESOURCES, MEDIUM_RESOURCES, RICH_RESOURCES},
			radialSymmetryCount: 4,
		},
		0, 0, w, h,
	)

	gm.performNAutomatasLike(10,
		rnd.RandInRange(5, 20),
		automat{
			drawsChar:           ROCKS,
			canDrawOn:           []tileCode{BUILDABLE_TERRAIN, SAND},
			maxCodeNear:         map[tileCode]int{ROCKS: 5},
			radialSymmetryCount: 4,
		},
		fromx, fromy, tox, toy,
	)

	gm.cleanupBadRadialSymmetry(2)
	gm.searchAndSetStartPoints(false, false, 4)
}
