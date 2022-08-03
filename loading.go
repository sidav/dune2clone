package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"image"
	"image/png"
	"os"
)

var (
	tilesAtlaces       = map[string]*spriteAtlas{}
	buildingsAtlaces   = map[int]*spriteAtlas{}
	unitChassisAtlaces = map[string]*spriteAtlas{}
	turretsAtlaces     = map[string]*spriteAtlas{}
	projectilesAtlaces = map[string]*spriteAtlas{}

	uiAtlaces = map[string]*spriteAtlas{}
)

func loadResources() {
	loadSprites()
}

func loadSprites() {
	tilesAtlaces = make(map[string]*spriteAtlas)
	unitChassisAtlaces = make(map[string]*spriteAtlas)
	turretsAtlaces = make(map[string]*spriteAtlas)

	tilesAtlaces["sand"] = CreateAtlasFromFile("resources/sprites/terrain/sand.png", 0, 0, 16, 16, 16, 16, 1, false)

	buildingsAtlaces = make(map[int]*spriteAtlas)
	buildingsAtlaces[BLD_BASE] = CreateAtlasFromFile("resources/sprites/buildings/base.png", 0, 0, 32, 32, 32, 32, 1, false)
	buildingsAtlaces[BLD_POWERPLANT] = CreateAtlasFromFile("resources/sprites/buildings/powerplant.png", 0, 0, 32, 32, 32, 32, 1, false)
	buildingsAtlaces[BLD_FACTORY] = CreateAtlasFromFile("resources/sprites/buildings/factory.png", 0, 0, 48, 32, 48, 32, 1, false)
	buildingsAtlaces[BLD_TURRET] = CreateAtlasFromFile("resources/sprites/buildings/cannon.png", 0, 0, 16, 16, 16, 16, 1, false)
	turretsAtlaces["cannon_turret"] = CreateDirectionalAtlasFromFile("resources/sprites/buildings/cannon_turret.png", 16, 16, 1, 2)


	unitChassisAtlaces["tank"] = CreateDirectionalAtlasFromFile("resources/sprites/units/tank_chassis.png", 16, 16, 1, 2)
	turretsAtlaces["tank"] = CreateDirectionalAtlasFromFile("resources/sprites/units/tank_cannon.png", 16, 16, 1, 2)

	unitChassisAtlaces["quad"] = CreateDirectionalAtlasFromFile("resources/sprites/units/quad.png", 16, 16, 1, 2)

	projectilesAtlaces = make(map[string]*spriteAtlas)
	projectilesAtlaces["shell"] = CreateDirectionalAtlasFromFile("resources/sprites/projectiles/shell.png", 16, 16, 1, 2)
	projectilesAtlaces["missile"] = CreateDirectionalAtlasFromFile("resources/sprites/projectiles/missile.png", 16, 16, 1, 2)

	uiAtlaces = make(map[string]*spriteAtlas)
	uiAtlaces["factionflag"] = CreateDirectionalAtlasFromFile("resources/sprites/ui/building_faction_flag.png", 8, 8, 1, 2)
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

func CreateAtlasFromFile(filename string, topleftx, toplefty, originalSpriteW, originalSpriteH,
	desiredSpriteW, desiredSpriteH, totalFrames int, createAllDirections bool) *spriteAtlas {

	file, _ := os.Open(filename)
	img, _ := png.Decode(file)
	file.Close()

	newAtlas := spriteAtlas{
		// spriteSize: desiredSpriteSize * int(SPRITE_SCALE_FACTOR),
	}
	if createAllDirections {
		newAtlas.atlas = make([][]rl.Texture2D, 4)
	} else {
		newAtlas.atlas = make([][]rl.Texture2D, 1)
	}
	// newAtlas.atlas
	for currFrame := 0; currFrame < totalFrames; currFrame++ {
		currPic := extractSubimageFromImage(img, topleftx+currFrame*originalSpriteW, toplefty, originalSpriteW, originalSpriteH)
		rlImg := rl.NewImageFromImage(currPic)
		rl.ImageResizeNN(rlImg, int32(desiredSpriteW)*int32(SPRITE_SCALE_FACTOR), int32(desiredSpriteH)*int32(SPRITE_SCALE_FACTOR))
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

func CreateDirectionalAtlasFromFile(filename string, originalSpriteSize, desiredSpriteSize, totalFrames, directionsInFile int) *spriteAtlas {

	file, _ := os.Open(filename)
	img, _ := png.Decode(file)
	file.Close()

	newAtlas := spriteAtlas{
		// spriteSize: desiredSpriteSize * int(SPRITE_SCALE_FACTOR),
	}
	newAtlas.atlas = make([][]rl.Texture2D, 4*directionsInFile)

	for currFrame := 0; currFrame < totalFrames; currFrame++ {
		for currDirectionFromFile := 0; currDirectionFromFile < directionsInFile; currDirectionFromFile++ {
			currPic := extractSubimageFromImage(img, currFrame*originalSpriteSize, currDirectionFromFile*originalSpriteSize, originalSpriteSize, originalSpriteSize)
			rlImg := rl.NewImageFromImage(currPic)
			rl.ImageResizeNN(rlImg, int32(desiredSpriteSize)*int32(SPRITE_SCALE_FACTOR), int32(desiredSpriteSize)*int32(SPRITE_SCALE_FACTOR))
			newAtlas.atlas[currDirectionFromFile] = append(newAtlas.atlas[currDirectionFromFile], rl.LoadTextureFromImage(rlImg))
			for i := 1; i < 4; i++ {
				rl.ImageRotateCW(rlImg)
				newAtlas.atlas[i*directionsInFile+currDirectionFromFile] =
					append(newAtlas.atlas[i*directionsInFile+currDirectionFromFile], rl.LoadTextureFromImage(rlImg))
			}
		}
	}

	return &newAtlas
}
