package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (r *renderer) drawLineInfoBox(x, y, w int32, title, info string, bgColor, textColor rl.Color) {
	var textSize int32 = 32
	var textCharW int32 = 22
	r.drawOutlinedRect(x, y, w, textSize+2, 2, rl.Green, bgColor)
	titleBoxW := int32((len(title)+1)*int(textCharW) + 2)
	r.drawOutlinedRect(x, y, titleBoxW, textSize+2, 2, rl.Green, bgColor)
	r.drawOutlinedRect(x+titleBoxW, y, w-titleBoxW, textSize+2, 2, rl.Green, bgColor)

	titlePosition := x + (titleBoxW / 2) - textCharW*int32(len(title))/2
	rl.DrawText(title, titlePosition-1, y+1, textSize, textColor)
	infoPosition := x + titleBoxW + int32(len(info))*textCharW/3
	rl.DrawText(info, infoPosition, y+1, textSize, textColor)
}

func (r *renderer) drawProgressCircle(x, y, radius int32, percent int, color rl.Color) {
	const OUTLINE_THICKNESS = 4
	rl.DrawCircleSector(rl.Vector2{
		X: float32(x),
		Y: float32(y),
	},
		float32(radius-OUTLINE_THICKNESS/2),
		180, 180-float32(360*percent)/100, 8, color)
	for i := -OUTLINE_THICKNESS / 2; i <= OUTLINE_THICKNESS/2; i++ {
		rl.DrawCircleLines(
			x,
			y,
			float32(radius+int32(i)),
			color)
	}
}

func (r *renderer) drawProgressBar(x, y, w int32, curr, max int, color *rl.Color) {
	const PG_H = 8
	const OUTLINE_THICKNESS = PG_H/2 - 2
	if color == nil {
		color = &rl.Green
	}
	for i := int32(0); i <= OUTLINE_THICKNESS/2; i++ {
		rl.DrawRectangleLines(x+i, y+i, w-i*2, PG_H-i*2, *color)
	}
	calculatedWidth := int32(curr) * w / int32(max)
	rl.DrawRectangle(x+OUTLINE_THICKNESS/2, y+OUTLINE_THICKNESS/2, calculatedWidth-OUTLINE_THICKNESS, PG_H-OUTLINE_THICKNESS, *color)
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
			if (i/PIXEL_SIZE)%2 == (j/PIXEL_SIZE)%2 {
				rl.DrawPixel(i, j, color)
			}
		}
	}
}
