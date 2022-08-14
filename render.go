package main

import (
	"dune2clone/geometry"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var DEFAULT_TINT = rl.RayWhite
var FOG_OF_WAR_TINT = rl.Color{
	R: 96,
	G: 96,
	B: 96,
	A: 255,
}

type renderer struct {
	camTopLeftX, camTopLeftY int32
	viewportW, viewportH     int32
	minimapRenderTextureMask rl.Texture2D
}

func (r *renderer) renderBattlefield(b *battlefield, pc *playerController) {
	if rl.IsWindowResized() || rl.IsWindowMaximized() {
		WINDOW_W = int32(rl.GetScreenWidth())
		WINDOW_H = int32(rl.GetScreenHeight())
	}
	r.viewportW = WINDOW_W
	r.viewportH = WINDOW_H
	r.camTopLeftX, r.camTopLeftY = int32(pc.camTopLeftX), int32(pc.camTopLeftY)

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	for x := range b.tiles {
		for y := range b.tiles[x] {
			r.renderTile(b, pc, x, y)
		}
	}

	// just testing
	//for i, f := range unitCannonsAtlaces[sTableUnits[b.units[0].code].cannonSpriteCode].atlas {
	//	rl.DrawTexture(
	//		f[0],
	//		int32(i * ORIGINAL_TILE_SIZE_IN_PIXELS*SPRITE_SCALE_FACTOR),
	//		int32(0),
	//		DEFAULT_TINT,
	//	)
	//}

	for i := b.buildings.Front(); i != nil; i = i.Next() {
		r.renderBuilding(b, pc, i.Value.(*building))
	}
	for i := b.units.Front(); i != nil; i = i.Next() {
		r.renderUnit(pc, i.Value.(*unit))
	}
	for p := b.projectiles.Front(); p != nil; p = p.Next() {
		r.renderProjectile(p.Value.(*projectile))
	}
	//for x := range b.tiles {
	//	for y := range b.tiles[x] {
	//		rl.DrawText(fmt.Sprintf("%d", b.costMapForMovement(x, y)),
	//			int32(x*TILE_SIZE_IN_PIXELS), int32(y * TILE_SIZE_IN_PIXELS), 24, rl.White)
	//	}
	//}

	r.renderUI(b, pc)

	rl.EndDrawing()
}

func (r *renderer) renderTile(b *battlefield, pc *playerController, x, y int) {
	t := b.tiles[x][y]
	osx, osy := r.physicalToOnScreenCoords(float64(x*TILE_PHYSICAL_SIZE), float64(y*TILE_PHYSICAL_SIZE))
	if !r.isRectInViewport(osx, osy, TILE_SIZE_IN_PIXELS, TILE_SIZE_IN_PIXELS) {
		return
	}
	if !pc.controlledFaction.exploredTilesMap[x][y] {
		// rl.DrawRectangle(osx, osy, TILE_SIZE_IN_PIXELS, TILE_SIZE_IN_PIXELS, rl.Black)
		return
	}
	if pc.controlledFaction.visibleTilesMap[x][y] {
		rl.DrawTexture(
			tilesAtlaces[sTableTiles[t.code].spriteCodes[t.spriteVariantIndex]].atlas[0][0],
			osx,
			osy,
			DEFAULT_TINT,
		)
		if t.resourcesAmount > 0 {
			rl.DrawTexture(
				tilesAtlaces["melange"].atlas[0][0],
				osx,
				osy,
				DEFAULT_TINT,
			)
		}
	} else {
		// r.drawDitheredRect(osx, osy, TILE_SIZE_IN_PIXELS, TILE_SIZE_IN_PIXELS, rl.Black)
		rl.DrawTexture(
			tilesAtlaces[sTableTiles[t.code].spriteCodes[t.spriteVariantIndex]].atlas[0][0],
			osx,
			osy,
			FOG_OF_WAR_TINT,
		)
	}

}

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
	if !pc.controlledFaction.canSeeActor(bld) {
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
		r.renderUnit(pc, bld.unitPlacedInside)
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

func (r *renderer) renderUnit(pc *playerController, u *unit) {
	x, y := u.centerX, u.centerY
	osx, osy := r.physicalToOnScreenCoords(x-0.5, y-0.5)
	// fmt.Printf("%d, %d \n", osx, osy)
	if !r.isRectInViewport(osx, osy, TILE_SIZE_IN_PIXELS, TILE_SIZE_IN_PIXELS) {
		return
	}
	if !pc.controlledFaction.canSeeActor(u) {
		return
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

func (r *renderer) renderProjectile(proj *projectile) {
	x, y := proj.centerX, proj.centerY
	osx, osy := r.physicalToOnScreenCoords(x, y)
	// fmt.Printf("%d, %d \n", osx, osy)
	if !r.isRectInViewport(osx, osy, TILE_SIZE_IN_PIXELS, TILE_SIZE_IN_PIXELS) {
		return
	}
	sprite := projectilesAtlaces[sTableProjectiles[proj.code].spriteCode][0].
		atlas[geometry.DegreeToRotationFrameNumber(proj.rotationDegree, 8)][0]
	rl.DrawTexture(
		sprite,
		osx-sprite.Width/2,
		osy-sprite.Height/2,
		DEFAULT_TINT, // proj.faction.factionColor,
	)
}

func (r *renderer) physicalToOnScreenCoords(physX, physY float64) (int32, int32) {
	pixx, pixy := r.physicalToPixelCoords(physX, physY)
	pixx -= r.camTopLeftX
	pixy -= r.camTopLeftY
	return int32(pixx), int32(pixy)
}

func (r *renderer) isRectInViewport(rx, ry, rw, rh int32) bool {
	return geometry.AreTwoCellRectsOverlapping32(rx, ry, rw, rh, 0, 0, r.viewportW, r.viewportH)
}

func (r *renderer) AreOnScreenCoordsInViewport(osx, osy int32) bool {
	// fmt.Printf("%d, %d \n", osx, osy)
	return osx >= 0 && osx < int32(r.viewportW) && osy >= 0 && osy < int32(r.viewportH)
}

func (r *renderer) physicalToPixelCoords(px, py float64) (int32, int32) {
	return int32(float32(px) * PIXEL_TO_PHYSICAL_RATIO), int32(float32(py) * PIXEL_TO_PHYSICAL_RATIO)
}
