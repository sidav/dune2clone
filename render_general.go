package main

import rl "github.com/gen2brain/raylib-go/raylib"


func (r *renderer) drawProgressCircle(x, y, radius int32, percent int, color rl.Color) {
	const OUTLINE_THICKNESS = 4
	rl.DrawCircleSector(rl.Vector2{
		X: float32(x),
		Y: float32(y),
	},
		float32(radius-OUTLINE_THICKNESS/2),
		180, 180-float32(360*percent)/100, 8, color)
	for i := -OUTLINE_THICKNESS/2; i <= OUTLINE_THICKNESS/2; i++ {
		rl.DrawCircleLines(
			x,
			y,
			float32(radius+int32(i)),
			color)
	}
}

func (r *renderer) drawOutlinedRect(x, y, w, h, outlineThickness int32, outlineColor, fillColor rl.Color) {
	// draw outline
	for i := int32(0); i < outlineThickness; i++ {
		rl.DrawRectangleLines(x+i, y+i, w-i*outlineThickness, h-i*outlineThickness, outlineColor)
	}
	rl.DrawRectangle(x+outlineThickness, y+outlineThickness, w-outlineThickness*2, h-outlineThickness*2, fillColor)
}

func (r *renderer) drawDitheredRect(x, y, w, h int32, color rl.Color) {
	const PIXEL_SIZE = 4
	// draw outline
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			if (i/PIXEL_SIZE) % 2 == (j/PIXEL_SIZE) % 2 {
				rl.DrawPixel(i, j, color)
			}
		}
	}
}
