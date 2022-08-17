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
	if !r.isRectInViewport(osx, osy, int32(w*TILE_SIZE_IN_PIXELS), int32(h*TILE_SIZE_IN_PIXELS)) {
		return
	}
	if !b.canFactionSeeActor(pc.controlledFaction, bld) {
		return
	}

	// get sprite
	var sprites []rl.Texture2D
	if bld.turret != nil {
		sprites = []rl.Texture2D{
			buildingsAtlaces[bld.getStaticData().spriteCode].atlas[0][0],
			turretsAtlaces[bld.turret.getStaticData().spriteCode][bld.faction.colorNumber].
				getSpriteByDegreeAndFrameNumber(bld.turret.rotationDegree, 0),
		}
	} else {
		sprites = []rl.Texture2D{
			buildingsAtlaces[bld.getStaticData().spriteCode].atlas[0][0],
		}
	}
	// draw sprite
	for _, s := range sprites {
		rl.DrawTexture(
			s,
			osx,
			osy,
			DEFAULT_TINT,
		)
	}

	if bld.isSelected {
		col := rl.Green
		rl.DrawRectangleLines(osx, osy, TILE_SIZE_IN_PIXELS*int32(w), TILE_SIZE_IN_PIXELS*int32(h), col)
		rl.DrawRectangleLines(osx-1, osy-1, TILE_SIZE_IN_PIXELS*int32(w)+2, TILE_SIZE_IN_PIXELS*int32(h), col)
		rl.DrawRectangleLines(osx+1, osy+1, TILE_SIZE_IN_PIXELS*int32(w)-2, TILE_SIZE_IN_PIXELS*int32(h), col)
		// render rally point
		if bld.rallyTileX != -1 {
			centerX, centerY := r.physicalToOnScreenCoords(bld.getPhysicalCenterCoords())
			rallyX, rallyY := r.physicalToOnScreenCoords(geometry.TileCoordsToPhysicalCoords(bld.rallyTileX, bld.rallytileY))
			rl.DrawLine(centerX, centerY, rallyX, rallyY, rl.White)
			rl.DrawRectangleLines(rallyX-TILE_SIZE_IN_PIXELS/4, rallyY-TILE_SIZE_IN_PIXELS/4,
				2*TILE_SIZE_IN_PIXELS/4, 2*TILE_SIZE_IN_PIXELS/4, rl.White)
		}
	}
	// render completion bar
	if bld.currentAction.getCompletionPercent() >= 0 {
		//r.drawProgressCircle(osx+int32(TILE_SIZE_IN_PIXELS*w/2),
		//	osy+int32(TILE_SIZE_IN_PIXELS*h/2),
		//	TILE_SIZE_IN_PIXELS/4, bld.currentAction.getCompletionPercent(), rl.Green)
		r.drawProgressBar(osx, osy-4, int32(TILE_SIZE_IN_PIXELS*w), bld.currentAction.getCompletionPercent(), 100, &rl.Blue)
	}
	if bld.currentHitpoints < bld.getStaticData().maxHitpoints {
		r.drawProgressBar(osx, osy-4, int32(TILE_SIZE_IN_PIXELS*w), bld.currentHitpoints, bld.getStaticData().maxHitpoints,
			&factionColors[bld.getFaction().colorNumber])
	}
	// render unit inside
	if bld.unitPlacedInside != nil {
		r.renderUnit(b, pc, bld.unitPlacedInside)
	}
	// render faction flag
	if bld.faction != nil && bld.getStaticData().w > 1 || bld.getStaticData().h > 1 {
		degree := (b.currentTick * 6) % 360
		rl.DrawTexture(
			uiAtlaces["factionflag"][bld.getFaction().colorNumber].getSpriteByDegreeAndFrameNumber(degree, 0),
			osx+2,
			osy+int32(bld.getStaticData().h*TILE_SIZE_IN_PIXELS)-uiAtlaces["factionflag"][bld.faction.colorNumber].atlas[0][0].Height-2,
			DEFAULT_TINT,
		)
	}
	// render energy if not enough
	if bld.faction != nil && bld.faction.getAvailableEnergy() < 0 {
		if (b.currentTick/30)%2 == 0 {
			icon := uiAtlaces["energyicon"][0].atlas[0][0]
			rl.DrawTexture(
				icon,
				osx+int32(bld.getStaticData().w*TILE_SIZE_IN_PIXELS)/2-icon.Width/2,
				osy+int32(bld.getStaticData().h*TILE_SIZE_IN_PIXELS)/2-icon.Height/2,
				DEFAULT_TINT,
			)
		}
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

	// get sprites
	var sprites []rl.Texture2D
	if u.turret != nil && u.turret.canRotate() {
		sprites = []rl.Texture2D{
			unitChassisAtlaces[sTableUnits[u.code].chassisSpriteCode][u.faction.colorNumber].getSpriteByDegreeAndFrameNumber(u.chassisDegree, 0),
			turretsAtlaces[u.turret.getStaticData().spriteCode][u.faction.colorNumber].getSpriteByDegreeAndFrameNumber(u.turret.rotationDegree, 0),
		}
	} else {
		sprites = []rl.Texture2D{
			unitChassisAtlaces[sTableUnits[u.code].chassisSpriteCode][u.faction.colorNumber].atlas[geometry.DegreeToRotationFrameNumber(u.chassisDegree, 8)][0],
		}
	}
	// draw sprites
	for _, s := range sprites {
		rl.DrawTexture(
			s,
			osx,
			osy,
			DEFAULT_TINT,
		)
	}
	if u.currentHitpoints < u.getStaticData().maxHitpoints {
		r.drawProgressBar(osx, osy-4, int32(TILE_SIZE_IN_PIXELS), u.currentHitpoints, u.getStaticData().maxHitpoints,
			&factionColors[u.getFaction().colorNumber])
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