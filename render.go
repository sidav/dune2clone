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
	// render ground units
	for i := b.units.Front(); i != nil; i = i.Next() {
		if !i.Value.(*unit).getStaticData().isAircraft {
			r.renderUnit(b, pc, i.Value.(*unit))
		}
	}

	for p := b.projectiles.Front(); p != nil; p = p.Next() {
		r.renderProjectile(p.Value.(*projectile))
	}

	// render aircrafts
	for i := b.units.Front(); i != nil; i = i.Next() {
		if i.Value.(*unit).getStaticData().isAircraft {
			r.renderUnit(b, pc, i.Value.(*unit))
		}
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
			amountForPoorTile := RESOURCE_IN_TILE_MAX_GENERATED/3
			resourceTileAtlasName := "melangerich"
			if t.resourcesAmount <= amountForPoorTile {
				resourceTileAtlasName = "melangepoor"
			} else if t.resourcesAmount <= 2*amountForPoorTile {
				resourceTileAtlasName = "melangemedium"
			}
			rl.DrawTexture(
				tilesAtlaces[resourceTileAtlasName].atlas[0][0],
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
