package map_generator

func generateByFourPlayersPatternAsymm(gm *GeneratedMap) {
	w, h := len(gm.Tiles), len(gm.Tiles[0])
	gm.reset()
	fromx, fromy, tox, toy := 0, 0, w-1, h-1
	// tox = 90 * tox / 100
	// toy = 90 * toy / 100

	gm.performNAutomatasLike(40,
		rnd.RandInRange(25, 70),
		0,
		automat{
			drawsChar: BUILDABLE_TERRAIN,
			canDrawOn: []tileCode{SAND},
		},
		fromx, fromy, tox, toy,
	)

	gm.performNAutomatasLike(25,
		rnd.RandInRange(1, 5),
		0,
		automat{
			drawsChar:   ROCKS,
			canDrawOn:   []tileCode{BUILDABLE_TERRAIN, SAND},
			maxCodeNear: map[tileCode]int{ROCKS: 5},
		},
		fromx, fromy, tox, toy,
	)

	gm.performNAutomatasLike(20,
		rnd.RandInRange(10, 15),
		0,
		automat{
			drawsChar:   POOR_RESOURCES,
			canDrawOn:   []tileCode{SAND},
			maxCodeNear: map[tileCode]int{BUILDABLE_TERRAIN: 0, ROCKS: 0},
		},
		0, 0, w, h,
	)
	gm.performNAutomatasLike(20,
		rnd.RandInRange(5, 10),
		0,
		automat{
			drawsChar:   MEDIUM_RESOURCES,
			canDrawOn:   []tileCode{POOR_RESOURCES},
			maxCodeNear: map[tileCode]int{BUILDABLE_TERRAIN: 0},
		},
		0, 0, w, h,
	)
	gm.performNAutomatasLike(10,
		rnd.RandInRange(1, 5),
		0,
		automat{
			drawsChar:   RICH_RESOURCES,
			canDrawOn:   []tileCode{MEDIUM_RESOURCES},
			maxCodeNear: map[tileCode]int{BUILDABLE_TERRAIN: 0},
		},
		0, 0, w, h,
	)
	gm.performNAutomatasLike(2*rnd.RandInRange(2, 8),
		0,
		1,
		automat{
			drawsChar: RESOURCE_VEIN,
			canDrawOn: []tileCode{POOR_RESOURCES, MEDIUM_RESOURCES, RICH_RESOURCES},
		},
		0, 0, w, h,
	)

	gm.searchAndSetStartPointsAsymmetric(4)
}
