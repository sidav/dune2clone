package main

import (
	"os"
	rl "github.com/gen2brain/raylib-go/raylib"
	"image"
	"image/png"
)

var (
	tilesAtlaces       = map[string]*spriteAtlas{}
	buildingsAtlaces   = map[int]*spriteAtlas{}
	unitChassisAtlaces = map[string]*spriteAtlas{}
	unitCannonsAtlaces = map[string]*spriteAtlas{}
	projectilesAtlaces = map[int]*spriteAtlas{}
)

func loadResources() {
	loadSprites()
}

func loadSprites() {
	tilesAtlaces = make(map[string]*spriteAtlas)
	tilesAtlaces["sand"] = CreateAtlasFromFile("resources/sprites/terrain/sand.png", 0, 0, 16, 16, 1, false)

	buildingsAtlaces = make(map[int]*spriteAtlas)
	buildingsAtlaces[BLD_BASE] = CreateAtlasFromFile("resources/sprites/buildings/base.png", 0, 0, 32, 32, 1, false)

	unitChassisAtlaces = make(map[string]*spriteAtlas)
	unitCannonsAtlaces = make(map[string]*spriteAtlas)
	unitChassisAtlaces["tank"] = CreateAtlasFromFile("resources/sprites/units/tank_chassis.png", 0, 0, 16, 16, 1, true)
	unitCannonsAtlaces["tank"] = CreateAtlasFromFile("resources/sprites/units/tank_cannon.png", 0, 0, 16, 16, 1, true)

}

func extractSubimageFromImage(img image.Image, fromx, fromy, w, h int) image.Image {
	minx, miny := img.Bounds().Min.X, img.Bounds().Min.Y
	//maxx, maxy := img.Bounds().Min.X, img.Bounds().Max.Y
	subImg := img.(*image.NRGBA).SubImage(
		image.Rect(minx+fromx, miny+fromy, minx+fromx+w, miny+fromy+h),
	)
	// reset img bounds, because RayLib goes nuts about it otherwise
	subImg.(*image.NRGBA).Rect = image.Rect(0, 0, w, h)
	return subImg
}

func CreateAtlasFromFile(filename string, topleftx, toplefty, originalSpriteSize, desiredSpriteSize, totalFrames int, createAllDirections bool) *spriteAtlas {

	file, _ := os.Open(filename)
	img, _ := png.Decode(file)
	file.Close()

	newAtlas := spriteAtlas{
		spriteSize: desiredSpriteSize * int(SPRITE_SCALE_FACTOR),
	}
	if createAllDirections {
		newAtlas.atlas = make([][]rl.Texture2D, 4)
	} else {
		newAtlas.atlas = make([][]rl.Texture2D, 1)
	}
	// newAtlas.atlas
	for currFrame := 0; currFrame < totalFrames; currFrame++ {
		currPic := extractSubimageFromImage(img, topleftx+currFrame*originalSpriteSize, toplefty, originalSpriteSize, originalSpriteSize)
		rlImg := rl.NewImageFromImage(currPic)
		rl.ImageResizeNN(rlImg, int32(desiredSpriteSize)*int32(SPRITE_SCALE_FACTOR), int32(desiredSpriteSize)*int32(SPRITE_SCALE_FACTOR))
		newAtlas.atlas[0] = append(newAtlas.atlas[0], rl.LoadTextureFromImage(rlImg))
		if createAllDirections {
			for i := 1; i < 4; i++ {
				rl.ImageRotateCW(rlImg)
				newAtlas.atlas[i] = append(newAtlas.atlas[i], rl.LoadTextureFromImage(rlImg))
			}
		}
	}

	return &newAtlas
}
