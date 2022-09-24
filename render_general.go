package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (r *renderer) drawLineInfoBox(x, y, w int32, title, info string, bgColor, textColor rl.Color) {
	var textSize int32 = 32
	var textCharW int32 = 20
	titleBoxW := int32((len(title)+1)*int(textCharW) + 2)
	infoW := int32(len(info)) * textCharW
	infoBoxW := w - titleBoxW

	r.drawOutlinedRect(x, y, w, textSize+2, 4, rl.DarkGreen, bgColor)
	rl.DrawLine(x+titleBoxW, y, x+titleBoxW, y+textSize+2, rl.DarkGreen)
	// r.drawOutlinedRect(x+titleBoxW, y, w-titleBoxW, textSize+2, 2, rl.Green, bgColor)

	titlePosition := x + (titleBoxW / 2) - textCharW*int32(len(title))/2 + 4
	infoPosition := x + titleBoxW + infoBoxW/2 - infoW/2

	r.drawText(title, titlePosition, y+2, textSize, textColor)
	r.drawText(info, infoPosition, y+2, textSize, textColor)
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
	calculatedWidth := int32(curr) * w / int32(max)
	rl.DrawRectangle(x+OUTLINE_THICKNESS/2, y+OUTLINE_THICKNESS/2, calculatedWidth-OUTLINE_THICKNESS, PG_H-OUTLINE_THICKNESS, *color)
	r.drawBoldRect(x, y, w, PG_H, OUTLINE_THICKNESS, rl.Black)
}

func (r *renderer) drawDiscreteProgressBar(x, y, w int32, curr, max int, color *rl.Color) {
	const PG_H = 8
	const OUTLINE_THICKNESS = PG_H/2 - 2
	partWidth := int32(6)
	for (w-2*OUTLINE_THICKNESS)%partWidth > 2 {
		debugWritef("%d, %d %% %d = %d", w, w-2*OUTLINE_THICKNESS, partWidth, (w-2*OUTLINE_THICKNESS)%partWidth)
		partWidth++
	}
	totalParts := (w - 2*OUTLINE_THICKNESS) / partWidth
	if color == nil {
		color = &rl.Green
	}
	r.drawBoldRect(x, y, w, PG_H, OUTLINE_THICKNESS, *color)
	calculatedWidth := int32(curr) * w / int32(max)
	rl.DrawRectangle(x+OUTLINE_THICKNESS/2, y+OUTLINE_THICKNESS/2, calculatedWidth-OUTLINE_THICKNESS, PG_H-OUTLINE_THICKNESS, *color)
	for i := int32(1); i <= totalParts; i++ {
		r.drawBoldRect(x, y, i*partWidth, PG_H, OUTLINE_THICKNESS, rl.Black)
	}
}

// TODO: re-use it (it's needed)
func (r *renderer) drawBoldRect(x, y, w, h, thickness int32, color rl.Color) {
	for i := int32(0); i < thickness; i++ {
		rl.DrawRectangleLines(x+i, y+i, w-2*i, h-2*i, color)
	}
}

func (r *renderer) drawOutlinedRect(x, y, w, h, outlineThickness int32, outlineColor, fillColor rl.Color) {
	rl.DrawRectangle(x, y, w, h, fillColor)
	// draw outline
	for i := int32(0); i < outlineThickness; i++ {
		rl.DrawRectangleLines(x+i, y+i, w-2*i, h-2*i, outlineColor)
	}
}

func (r *renderer) drawText(text string, x, y, size int32, color rl.Color) {
	rl.DrawTextEx(defaultFont, text, rl.Vector2{float32(x), float32(y)}, float32(size), 0, color)
}

func (r *renderer) drawDitheredRect(x, y, w, h int32, color rl.Color) {
	pixelSize := int32(TILE_SIZE_IN_PIXELS / 16)
	for i := int32(0); i < w/pixelSize; i++ {
		for j := int32(0); j < h/pixelSize; j++ {
			if i%2 == j%2 {
				rl.DrawRectangle(x+i*pixelSize, y+j*pixelSize, pixelSize, pixelSize, color)
			}
		}
	}
}
