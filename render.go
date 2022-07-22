package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var DEFAULT_TINT = rl.RayWhite

type renderer struct {
	cameraCenterX, cameraCenterY int
	viewportW                    int
}

func (r *renderer) renderBattlefield(b *battlefield) {
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
		r.renderBuilding(b.buildings[i])
	}
	for i := range b.units {
		r.renderUnit(b.units[i])
	}

	rl.DrawText(fmt.Sprintf("TICK %d", b.currentTick), 0, 0, 24, rl.White)

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
				int32(osx),
				int32(osy),
				DEFAULT_TINT,
			)
		}
	}
}

func (r *renderer) renderBuilding(b *building) {
	x, y := tileCoordsToPhysicalCoords(b.topLeftX, b.topLeftY)
	x -= 0.5
	y -= 0.5
	osx, osy := r.physicalToOnScreenCoords(x, y)
	// fmt.Printf("%d, %d \n", osx, osy)
	if r.AreOnScreenCoordsInViewport(osx, osy) {
		rl.DrawTexture(
			b.getSprite(),
			int32(osx),
			int32(osy),
			DEFAULT_TINT,
		)
	}
}

func (r *renderer) renderUnit(u *unit) {
	x, y := u.centerX, u.centerY
	osx, osy := r.physicalToOnScreenCoords(x-0.5, y-0.5)
	// fmt.Printf("%d, %d \n", osx, osy)
	if r.AreOnScreenCoordsInViewport(osx, osy) {
		sprites := u.getPartsSprites()
		for _, s := range sprites {
			rl.DrawTexture(
				s,
				int32(osx),
				int32(osy),
				DEFAULT_TINT,
			)
		}
		if u.isSelected {
			col := rl.DarkGreen
			rl.DrawRectangleLines(int32(osx), int32(osy), TILE_SIZE_IN_PIXELS, TILE_SIZE_IN_PIXELS, col)
			rl.DrawRectangleLines(int32(osx-1), int32(osy-1), TILE_SIZE_IN_PIXELS+2, TILE_SIZE_IN_PIXELS+2, col)
			rl.DrawRectangleLines(int32(osx+1), int32(osy+1), TILE_SIZE_IN_PIXELS-2, TILE_SIZE_IN_PIXELS-2, col)
		}
	}
}

func (r *renderer) physicalToOnScreenCoords(physX, physY float64) (int, int) {
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
	return pixx, pixy
}

func (r *renderer) AreOnScreenCoordsInViewport(osx, osy int) bool {
	// fmt.Printf("%d, %d \n", osx, osy)
	return osx >= 0 && osx < r.viewportW && osy >= 0 && osy < WINDOW_H
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
