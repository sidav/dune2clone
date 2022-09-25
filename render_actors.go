package main

import (
	"dune2clone/geometry"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (r *renderer) renderBuilding(b *battlefield, pc *playerController, bld *building) {
	x, y := geometry.TileCoordsToPhysicalCoords(bld.topLeftX, bld.topLeftY)
	x -= 0.5
	y -= 0.5
	osx, osy := r.physicalToOnScreenCoords(x, y)
	w, h := bld.getStaticData().w, bld.getStaticData().h
	// fmt.Printf("%d, %d \n", osx, osy)
	// render rally point. Called BEFORE viewport check.
	if bld.rallyTileX != -1 {
		centerX, centerY := r.physicalToOnScreenCoords(bld.getPhysicalCenterCoords())
		rallyX, rallyY := r.physicalToOnScreenCoords(geometry.TileCoordsToPhysicalCoords(bld.rallyTileX, bld.rallytileY))
		rl.DrawLine(centerX, centerY, rallyX, rallyY, rl.White)
		rl.DrawRectangleLines(rallyX-TILE_SIZE_IN_PIXELS/4, rallyY-TILE_SIZE_IN_PIXELS/4,
			2*TILE_SIZE_IN_PIXELS/4, 2*TILE_SIZE_IN_PIXELS/4, rl.White)
		r.renderFactionFlagAt(bld.faction, rallyX, rallyY)
	}
	if !r.isRectInViewport(osx, osy, int32(w*TILE_SIZE_IN_PIXELS), int32(h*TILE_SIZE_IN_PIXELS)) {
		return
	}
	if !b.canFactionSeeActor(pc.controlledFaction, bld) {
		return
	}

	// draw sprite
	if bld.isUnderConstruction() && (b.currentTick/10)%2 != 0 {
		// under construction
		underConstructionAtlas := buildingsAtlaces["underconstruction"]
		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				frameNumber := (x + y + b.currentTick/250) % underConstructionAtlas.totalFrames()
				rl.DrawTexture(
					underConstructionAtlas.getSpriteByFrame(frameNumber),
					osx+int32(x)*TILE_SIZE_IN_PIXELS,
					osy+int32(y)*TILE_SIZE_IN_PIXELS,
					DEFAULT_TINT,
				)
			}
		}
	} else {
		var sprites []rl.Texture2D
		frameNumber := b.currentTick/(DESIRED_TPS /4)
		if bld.turret != nil && bld.turret.getStaticData().spriteCode != "" {
			sprites = []rl.Texture2D{
				buildingsAtlaces[bld.getStaticData().spriteCode].getSpriteByColorAndFrame(bld.getFaction().colorNumber, frameNumber),
				turretsAtlaces[bld.turret.getStaticData().spriteCode].getSpriteByColorDegreeAndFrameNumber(bld.faction.colorNumber, bld.turret.rotationDegree, 0),
			}
		} else {
			sprites = []rl.Texture2D{
				buildingsAtlaces[bld.getStaticData().spriteCode].getSpriteByColorAndFrame(bld.getFaction().colorNumber, frameNumber),
			}
		}
		for _, s := range sprites {
			rl.DrawTexture(
				s,
				osx,
				osy,
				DEFAULT_TINT,
			)
		}
	}

	if bld.isSelected {
		col := rl.Green
		rl.DrawRectangleLines(osx, osy, TILE_SIZE_IN_PIXELS*int32(w), TILE_SIZE_IN_PIXELS*int32(h), col)
		rl.DrawRectangleLines(osx-1, osy-1, TILE_SIZE_IN_PIXELS*int32(w)+2, TILE_SIZE_IN_PIXELS*int32(h), col)
		rl.DrawRectangleLines(osx+1, osy+1, TILE_SIZE_IN_PIXELS*int32(w)-2, TILE_SIZE_IN_PIXELS*int32(h), col)
	}
	// render completion bar
	if bld.currentAction.getCompletionPercent() >= 0 {
		//r.drawProgressCircle(osx+int32(TILE_SIZE_IN_PIXELS*w/2),
		//	osy+int32(TILE_SIZE_IN_PIXELS*h/2),
		//	TILE_SIZE_IN_PIXELS/4, bld.currentAction.getCompletionPercent(), rl.Green)
		r.drawProgressBar(osx, osy+2, int32(TILE_SIZE_IN_PIXELS*w), bld.currentAction.getCompletionPercent(), 100, &rl.Blue)
	}
	if bld.currentHitpoints < bld.getStaticData().maxHitpoints {
		r.drawProgressBar(osx, osy-4, int32(TILE_SIZE_IN_PIXELS*w), bld.currentHitpoints, bld.getStaticData().maxHitpoints,
			&factionColors[bld.getFaction().colorNumber])
	}
	// render unit inside
	if bld.unitPlacedInside != nil && bld.currentAction.code != ACTION_BEING_BUILT {
		r.renderUnit(b, pc, bld.unitPlacedInside)
	}
	// render faction flag
	if bld.faction != nil && bld.getStaticData().w > 1 || bld.getStaticData().h > 1 {
		r.renderFactionFlagAt(
			bld.faction,
			osx+4,
			osy+int32(bld.getStaticData().h*TILE_SIZE_IN_PIXELS),
		)
	}
	// render energy if not enough
	if bld.faction != nil && bld.faction.getAvailableEnergy() < 0 {
		r.renderBlinkingIconCenteredAt("energyicon",
			osx+int32(bld.getStaticData().w*TILE_SIZE_IN_PIXELS)/2,
			osy+int32(bld.getStaticData().h*TILE_SIZE_IN_PIXELS)/2,
			0,
		)
	}
	if bld.currentOrder.code == ORDER_WAIT_FOR_BUILDING_PLACEMENT {
		r.renderBlinkingIconCenteredAt("readyicon",
			osx+int32(bld.getStaticData().w*TILE_SIZE_IN_PIXELS)/2,
			osy+int32(bld.getStaticData().h*TILE_SIZE_IN_PIXELS)/2,
			1,
		)
	}
}

func (r *renderer) renderUnit(b *battlefield, pc *playerController, u *unit) {
	x, y := u.centerX, u.centerY
	osx, osy := r.physicalToOnScreenCoords(x-0.5, y-0.5)
	// fmt.Printf("%d, %d \n", osx, osy)
	if !r.isRectInViewport(osx, osy, TILE_SIZE_IN_PIXELS, TILE_SIZE_IN_PIXELS) {
		return
	}
	if !b.canFactionSeeActor(pc.controlledFaction, u) {
		return
	}

	// render unit inside
	if u.carriedUnit != nil {
		r.renderUnit(b, pc, u.carriedUnit)
	}

	// draw chassis sprite
	relativeSquadCoords := getListOfRelativeCoordsForSquadMembers(u.squadSize)
	for _, relSquadCoord := range relativeSquadCoords {
		osx, osy := r.physicalToOnScreenCoords(x+relSquadCoord[0]-0.5, y+relSquadCoord[1]-0.5)
		rl.DrawTexture(
			unitChassisAtlaces[sTableUnits[u.code].chassisSpriteCode].getSpriteByColorDegreeAndFrameNumber(u.faction.colorNumber, u.chassisDegree, 0),
			osx,
			osy,
			DEFAULT_TINT,
		)
		// draw turrets
		for turrIndex := range u.turrets {
			if u.turrets[turrIndex].getStaticData().spriteCode == "" {
				continue
			}
			usd := u.getStaticData()
			// calculate turret displacement
			dsplX, dsplY := usd.turretsData[turrIndex].turretCenterX, usd.turretsData[turrIndex].turretCenterY
			if dsplX != 0 || dsplY != 0 {
				// rotate according to units face
				chassisShownDegree := geometry.SnapDegreeToFixedDirections(u.chassisDegree, 8)
				dsplX, dsplY = geometry.RotateFloat64Vector(dsplX, dsplY, chassisShownDegree)
			}
			turrOsX, turrOsY := r.physicalToOnScreenCoords(x+dsplX, y+dsplY)
			sprite := turretsAtlaces[u.turrets[turrIndex].getStaticData().spriteCode].getSpriteByColorDegreeAndFrameNumber(u.faction.colorNumber, u.turrets[turrIndex].rotationDegree, 0)
			rl.DrawTexture(
				sprite,
				turrOsX-sprite.Width/2,
				turrOsY-sprite.Height/2,
				DEFAULT_TINT,
			)
		}
	}

	if u.currentHitpoints < u.getStaticData().maxHitpoints {
		r.drawProgressBar(osx, osy-4, int32(TILE_SIZE_IN_PIXELS), u.currentHitpoints, u.getStaticData().maxHitpoints,
			&factionColors[u.getFaction().colorNumber])
	}
	// render completion bar
	if u.currentAction.getCompletionPercent() >= 0 {
		r.drawProgressBar(osx, osy+2, int32(TILE_SIZE_IN_PIXELS), u.currentAction.getCompletionPercent(), 100, &rl.Blue)
	}
	if u.isSelected {
		col := rl.DarkGreen
		circleX := osx + TILE_SIZE_IN_PIXELS/2
		circleY := osy + TILE_SIZE_IN_PIXELS/2
		rl.DrawCircleLines(circleX, circleY, TILE_SIZE_IN_PIXELS/2, col)
		rl.DrawCircleLines(circleX, circleY, TILE_SIZE_IN_PIXELS/2-1, col)
		//rl.DrawCircleLines(circleX, circleY, TILE_SIZE_IN_PIXELS/2-2, col)
		rl.DrawCircleLines(circleX, circleY, TILE_SIZE_IN_PIXELS/2-3, col)
		rl.DrawCircleLines(circleX, circleY, TILE_SIZE_IN_PIXELS/2-4, col)
		//rl.DrawRectangleLines(osx, osy, TILE_SIZE_IN_PIXELS, TILE_SIZE_IN_PIXELS, col)
		//rl.DrawRectangleLines(int32(osx-1), int32(osy-1), TILE_SIZE_IN_PIXELS+2, TILE_SIZE_IN_PIXELS+2, col)
		//rl.DrawRectangleLines(int32(osx+1), int32(osy+1), TILE_SIZE_IN_PIXELS-2, TILE_SIZE_IN_PIXELS-2, col)
	}
}
