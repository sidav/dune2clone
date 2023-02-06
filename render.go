package main

import (
	"dune2clone/geometry"
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
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
	btl                      *battlefield

	// for technical/debug output
	lastFrameRenderingTime time.Duration
	timeDebugInfosToRender []debugTimeInfo
}

func (r *renderer) renderBattlefield(b *battlefield, pc *playerController) {
	r.btl = b
	timeFrameRenderStarted := time.Now()
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

	for i := b.buildings.Front(); i != nil; i = i.Next() {
		r.renderBuilding(b, pc, i.Value.(*building))
	}
	// render ground units
	for i := b.units.Front(); i != nil; i = i.Next() {
		if !i.Value.(*unit).getStaticData().IsAircraft {
			r.renderUnit(b, pc, i.Value.(*unit))
		}
	}

	for p := b.projectiles.Front(); p != nil; p = p.Next() {
		r.renderProjectile(p.Value.(*projectile), pc)
	}

	// render aircrafts
	for i := b.units.Front(); i != nil; i = i.Next() {
		if i.Value.(*unit).getStaticData().IsAircraft {
			r.renderUnit(b, pc, i.Value.(*unit))
		}
	}

	for i := b.effects.Front(); i != nil; i = i.Next() {
		r.renderEffect(i.Value.(*effect))
	}
	// r.renderCollisionMap(b, pc)

	r.lastFrameRenderingTime = time.Since(timeFrameRenderStarted) / time.Millisecond

	r.renderUI(b, pc)
	rl.EndDrawing()
}

func (r *renderer) renderCollisionMap(b *battlefield, pc *playerController) {
	w, h := b.getSize()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			osx, osy := r.physicalToOnScreenCoords(float64(x*TILE_PHYSICAL_SIZE), float64(y*TILE_PHYSICAL_SIZE))
			if !r.isRectInViewport(osx, osy, TILE_SIZE_IN_PIXELS, TILE_SIZE_IN_PIXELS) {
				continue
			}
			if !pc.controlledFaction.hasTileAtCoordsExplored(x, y) {
				return
			}
			// Debug draw collisions
			if b.tiles[x][y].isOccupiedByActor != nil {
				r.drawBoldRect(osx, osy, TILE_SIZE_IN_PIXELS, TILE_SIZE_IN_PIXELS, 3, rl.Magenta)
				r.drawText(b.tiles[x][y].isOccupiedByActor.getName(), osx, osy, 18, rl.Magenta)
			}
		}
	}
}

func (r *renderer) renderTile(b *battlefield, pc *playerController, x, y int) {
	t := b.tiles[x][y]
	osx, osy := r.physicalToOnScreenCoords(float64(x*TILE_PHYSICAL_SIZE), float64(y*TILE_PHYSICAL_SIZE))
	if !r.isRectInViewport(osx, osy, TILE_SIZE_IN_PIXELS, TILE_SIZE_IN_PIXELS) {
		return
	}
	if !pc.controlledFaction.hasTileAtCoordsExplored(x, y) {
		// rl.DrawRectangle(osx, osy, TILE_SIZE_IN_PIXELS, TILE_SIZE_IN_PIXELS, rl.Black)
		return
	}
	tintToUse := FOG_OF_WAR_TINT
	if pc.controlledFaction.seesTileAtCoords(x, y) {
		tintToUse = DEFAULT_TINT
	}
	rl.DrawTexture(
		tilesAtlaces[sTableTiles[t.code].spriteCodes[t.spriteVariantIndex]].getSpriteByFrame(0),
		osx,
		osy,
		tintToUse,
	)
	if t.hasResourceVein {
		rl.DrawTexture(
			tilesAtlaces["resourcevein"].getSpriteByFrame(0),
			osx,
			osy,
			tintToUse,
		)
	}
	if t.resourcesAmount > 0 {
		resourceTileAtlasName := "melangerich"
		if t.resourcesAmount <= RESOURCE_IN_TILE_POOR_MAX {
			resourceTileAtlasName = "melangepoor"
		} else if t.resourcesAmount <= RESOURCE_IN_TILE_MEDIUM_MAX {
			resourceTileAtlasName = "melangemedium"
		}
		rl.DrawTexture(
			tilesAtlaces[resourceTileAtlasName].getSpriteByFrame(0),
			osx,
			osy,
			tintToUse,
		)
	}
}

func (r *renderer) renderProjectile(proj *projectile, pc *playerController) {
	x, y := proj.centerX, proj.centerY
	osx, osy := r.physicalToOnScreenCoords(x, y)
	// fmt.Printf("%d, %d \n", osx, osy)
	if !r.isRectInViewport(osx, osy, TILE_SIZE_IN_PIXELS, TILE_SIZE_IN_PIXELS) {
		return
	}
	tx, ty := geometry.TrueCoordsToTileCoords(x, y)
	if !pc.controlledFaction.seesTileAtCoords(tx, ty) {
		return
	}
	sprite := projectilesAtlaces[proj.getStaticData().SpriteCode].getSpriteByColorDegreeAndFrameNumber(0, proj.rotationDegree, 0)
	rl.DrawTexture(
		sprite,
		osx-sprite.Width/2,
		osy-sprite.Height/2,
		DEFAULT_TINT, // proj.faction.factionColor,
	)
}

func (r *renderer) renderEffect(e *effect) {
	// debugWritef("Percent is %d", e.getExpirationPercent(r.btl.currentTick))
	if e.creationTick <= r.btl.currentTick && e.getExpirationPercent(r.btl.currentTick) <= 100 {
		x, y := e.centerX, e.centerY
		osx, osy := r.physicalToOnScreenCoords(x, y)
		// fmt.Printf("%d, %d \n", osx, osy)
		if !r.isRectInViewport(osx, osy, TILE_SIZE_IN_PIXELS, TILE_SIZE_IN_PIXELS) {
			return
		}
		neededAtlas := effectsAtlaces[e.getStaticData().spriteCode]
		expPercent := e.getExpirationPercent(r.btl.currentTick)
		currentFrame := geometry.GetPartitionIndex(expPercent, 0, 100, neededAtlas.totalFrames())
		if e.splashCircleRadius > 0 {
			radius := float32(float64(expPercent*TILE_SIZE_IN_PIXELS)*e.splashCircleRadius) / 100.0
			rl.DrawCircleLines(osx, osy, radius, rl.Red)
			rl.DrawCircleLines(osx, osy, radius+1, rl.Maroon)
			rl.DrawCircleLines(osx, osy, radius+2, rl.Yellow)
		}
		rl.DrawTexture(
			neededAtlas.getSpriteByFrame(currentFrame),
			osx-neededAtlas.getSpriteByFrame(currentFrame).Width/2,
			osy-neededAtlas.getSpriteByFrame(currentFrame).Height/2,
			DEFAULT_TINT, // proj.faction.factionColor,
		)
	}
}

func (r *renderer) renderFactionFlagAt(f *faction, leftX, bottomY int32) {
	frame := (6 * r.btl.currentTick / config.TargetTPS) % uiAtlaces["factionflag"].totalFrames()
	spr := uiAtlaces["factionflag"].getSpriteByColorAndFrame(f.colorNumber, frame)
	rl.DrawTexture(
		spr,
		leftX,
		bottomY-spr.Height,
		DEFAULT_TINT,
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
