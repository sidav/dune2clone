package main

import (
	"dune2clone/geometry"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var DEFAULT_TINT = rl.RayWhite

type renderer struct {
	cameraCenterX, cameraCenterY int
	viewportW                    int
}

func (r *renderer) renderBattlefield(b *battlefield, pc *playerController) {
	r.viewportW = WINDOW_W

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	for x := range b.tiles {
		for y := range b.tiles[x] {
			r.renderTile(b, x, y)
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

	for i := range b.buildings {
		r.renderBuilding(b, b.buildings[i])
	}
	for i := range b.units {
		r.renderUnit(b.units[i])
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

func (r *renderer) renderTile(b *battlefield, x, y int) {
	t := b.tiles[x][y]
	spr := t.getSpritesAtlas()
	if spr != nil {
		osx, osy := r.physicalToOnScreenCoords(float64(x*TILE_PHYSICAL_SIZE), float64(y*TILE_PHYSICAL_SIZE))
		if r.AreOnScreenCoordsInViewport(osx, osy) {
			rl.DrawTexture(
				t.getSprite(),
				osx,
				osy,
				DEFAULT_TINT,
			)
		}
	}
}

func (r *renderer) renderBuilding(b *battlefield, bld *building) {
	x, y := geometry.TileCoordsToPhysicalCoords(bld.topLeftX, bld.topLeftY)
	x -= 0.5
	y -= 0.5
	osx, osy := r.physicalToOnScreenCoords(x, y)
	// fmt.Printf("%d, %d \n", osx, osy)
	if r.AreOnScreenCoordsInViewport(osx, osy) {
		// get sprite
		var sprites []rl.Texture2D
		if bld.turret != nil {
			sprites = []rl.Texture2D{
				buildingsAtlaces[bld.code].atlas[0][0],
				turretsAtlaces[bld.turret.getStaticData().spriteCode].
					getSpriteByDegreeAndFrameNumber(bld.turret.rotationDegree, 0),
			}
		} else {
			sprites = []rl.Texture2D{
				buildingsAtlaces[bld.code].atlas[0][0],
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

		w, h := bld.getStaticData().w, bld.getStaticData().h
		if bld.isSelected {
			col := rl.Green
			rl.DrawRectangleLines(osx, osy, TILE_SIZE_IN_PIXELS*int32(w), TILE_SIZE_IN_PIXELS*int32(h), col)
			rl.DrawRectangleLines(osx-1, osy-1, TILE_SIZE_IN_PIXELS*int32(w)+2, TILE_SIZE_IN_PIXELS*int32(h), col)
			rl.DrawRectangleLines(osx+1, osy+1, TILE_SIZE_IN_PIXELS*int32(w)-2, TILE_SIZE_IN_PIXELS*int32(h), col)
		}
		// render completion circle
		if bld.currentAction.getCompletionPercent() >= 0 {
			r.drawProgressCircle(osx+int32(TILE_SIZE_IN_PIXELS*w/2),
				osy+int32(TILE_SIZE_IN_PIXELS*h/2),
				TILE_SIZE_IN_PIXELS/4, bld.currentAction.getCompletionPercent(), rl.Green)
		}
		// render faction flag
		if bld.faction != nil {
			degree := (b.currentTick * 6) % 360
			rl.DrawTexture(
				uiAtlaces["factionflag"].getSpriteByDegreeAndFrameNumber(degree, 0),
				osx+2,
				osy+int32(bld.getStaticData().h*TILE_SIZE_IN_PIXELS)-uiAtlaces["factionflag"].atlas[0][0].Height-2,
				bld.faction.factionColor,
			)
		}
	}
}

func (r *renderer) renderUnit(u *unit) {
	x, y := u.centerX, u.centerY
	osx, osy := r.physicalToOnScreenCoords(x-0.5, y-0.5)
	// fmt.Printf("%d, %d \n", osx, osy)
	if r.AreOnScreenCoordsInViewport(osx, osy) {
		// get sprites
		var sprites []rl.Texture2D
		if u.turret.canRotate() {
			sprites = []rl.Texture2D{
				unitChassisAtlaces[sTableUnits[u.code].chassisSpriteCode].getSpriteByDegreeAndFrameNumber(u.chassisDegree, 0),
				turretsAtlaces[u.turret.getStaticData().spriteCode].getSpriteByDegreeAndFrameNumber(u.turret.rotationDegree, 0),
			}
		} else {
			sprites = []rl.Texture2D{
				unitChassisAtlaces[sTableUnits[u.code].chassisSpriteCode].atlas[geometry.DegreeToRotationFrameNumber(u.chassisDegree, 8)][0],
			}
		}
		// draw sprites
		for _, s := range sprites {
			rl.DrawTexture(
				s,
				osx,
				osy,
				u.faction.factionColor,
			)
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
}

func (r *renderer) renderProjectile(proj *projectile) {
	x, y := proj.centerX, proj.centerY
	osx, osy := r.physicalToOnScreenCoords(x-0.5, y-0.5)
	// fmt.Printf("%d, %d \n", osx, osy)
	if r.AreOnScreenCoordsInViewport(osx, osy) {
		sprite := projectilesAtlaces[sTableProjectiles[proj.code].spriteCode].
			atlas[geometry.DegreeToRotationFrameNumber(proj.rotationDegree, 8)][0]
		rl.DrawTexture(
			sprite,
			osx,
			osy,
			proj.faction.factionColor,
		)
	}
}

func (r *renderer) physicalToOnScreenCoords(physX, physY float64) (int32, int32) {
	pixx, pixy := r.physicalToPixelCoords(physX, physY)
	if !r.doesLevelFitInScreenHorizontally() {
		if r.cameraCenterX > MAP_W*TILE_SIZE_IN_PIXELS-r.viewportW/2 {
			pixx = pixx - MAP_W*TILE_SIZE_IN_PIXELS + r.viewportW
		} else if r.cameraCenterX > r.viewportW/2 {
			pixx = pixx - r.cameraCenterX + r.viewportW/2
		}
	}
	if !r.doesLevelFitInScreenVertically() {
		if r.cameraCenterY > MAP_H*TILE_SIZE_IN_PIXELS-WINDOW_H/2 {
			pixy = pixy - MAP_H*TILE_SIZE_IN_PIXELS + WINDOW_H
		} else if r.cameraCenterY > WINDOW_H/2 {
			pixy = pixy - r.cameraCenterY + WINDOW_H/2
		}
	}
	return int32(pixx), int32(pixy)
}

func (r *renderer) AreOnScreenCoordsInViewport(osx, osy int32) bool {
	// fmt.Printf("%d, %d \n", osx, osy)
	return osx >= 0 && osx < int32(r.viewportW) && osy >= 0 && osy < WINDOW_H
}

func (r *renderer) physicalToPixelCoords(px, py float64) (int, int) {
	return int(float32(px) * PIXEL_TO_PHYSICAL_RATIO), int(float32(py) * PIXEL_TO_PHYSICAL_RATIO)
}

func (r *renderer) doesLevelFitInScreenHorizontally() bool {
	return MAP_W*TILE_SIZE_IN_PIXELS <= WINDOW_W
}

func (r *renderer) doesLevelFitInScreenVertically() bool {
	return MAP_H*TILE_SIZE_IN_PIXELS <= WINDOW_H
}
