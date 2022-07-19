package main

import (
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

	for i := range b.buildings {
		r.renderBuilding(b.buildings[i])
	}

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
	x, y := b.topLeftX, b.topLeftY
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
