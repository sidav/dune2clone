package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"image"
	"image/color"
	"image/png"
	"os"
)

var (
	defaultFont rl.Font
	// index of array is faction color.
	tilesAtlaces       = map[string]*spriteAtlas{}
	buildingsAtlaces   = map[string][]*spriteAtlas{}
	unitChassisAtlaces = map[string][]*spriteAtlas{}
	turretsAtlaces     = map[string][]*spriteAtlas{}
	projectilesAtlaces = map[string][]*spriteAtlas{}
	effectsAtlaces     = map[string]*spriteAtlas{}

	uiAtlaces = map[string][]*spriteAtlas{}
)

func loadResources() {
	// a := int32(255)
	// defaultFont = rl.LoadFontEx("", 96, &a, 255)
	defaultFont = rl.LoadFont("resources/flexi.ttf")
	// rl.GenTextureMipmaps(&defaultFont.Texture)
	loadSprites()
}

func loadSprites() {
	tilesAtlaces = make(map[string]*spriteAtlas)
	buildingsAtlaces = make(map[string][]*spriteAtlas)
	unitChassisAtlaces = make(map[string][]*spriteAtlas)
	turretsAtlaces = make(map[string][]*spriteAtlas)
	projectilesAtlaces = make(map[string][]*spriteAtlas)
	uiAtlaces = make(map[string][]*spriteAtlas)

	currPath := "resources/sprites/terrain/"
	tilesAtlaces["sand1"] = CreateAtlasFromFile(currPath+"sand1.png", 0, 0, 16, 16, 16, 16, 1, false, false)[0]
	tilesAtlaces["sand2"] = CreateAtlasFromFile(currPath+"sand2.png", 0, 0, 16, 16, 16, 16, 1, false, false)[0]
	tilesAtlaces["sand3"] = CreateAtlasFromFile(currPath+"sand3.png", 0, 0, 16, 16, 16, 16, 1, false, false)[0]
	tilesAtlaces["buildable1"] = CreateAtlasFromFile(currPath+"buildable1.png", 0, 0, 16, 16, 16, 16, 1, false, false)[0]
	tilesAtlaces["rock1"] = CreateAtlasFromFile(currPath+"rocks.png", 0, 0, 16, 16, 16, 16, 1, false, false)[0]
	tilesAtlaces["buildabledamaged"] = CreateAtlasFromFile(currPath+"buildable_damaged.png", 0, 0, 16, 16, 16, 16, 1, false, false)[0]
	tilesAtlaces["resourcevein"] = CreateAtlasFromFile(currPath+"resource_vein.png", 0, 0, 16, 16, 16, 16, 1, false, false)[0]
	tilesAtlaces["melangerich"] = CreateAtlasFromFile(currPath+"melange_rich.png", 0, 0, 16, 16, 16, 16, 1, false, false)[0]
	tilesAtlaces["melangemedium"] = CreateAtlasFromFile(currPath+"melange_medium.png", 0, 0, 16, 16, 16, 16, 1, false, false)[0]
	tilesAtlaces["melangepoor"] = CreateAtlasFromFile(currPath+"melange_poor.png", 0, 0, 16, 16, 16, 16, 1, false, false)[0]

	currPath = "resources/sprites/buildings/"
	// WARNING: IT HAS FRAMES
	buildingsAtlaces["underconstruction"] = CreateAtlasFromFile(currPath+"under_construction.png", 0, 0, 32, 32, 16, 16, 4, false, false)

	buildingsAtlaces["base"] = CreateAtlasFromFile(currPath+"base.png", 0, 0, 64, 64, 32, 32, 1, false, true)
	buildingsAtlaces["powerplant"] = CreateAtlasFromFile(currPath+"powerplant.png", 0, 0, 64, 64, 32, 32, 1, false, true)
	buildingsAtlaces["barracks"] = CreateAtlasFromFile(currPath+"barracks.png", 0, 0, 32, 32, 32, 32, 1, false, true)
	buildingsAtlaces["factory"] = CreateAtlasFromFile(currPath+"factory.png", 0, 0, 96, 64, 48, 32, 1, false, true)
	buildingsAtlaces["airfactory"] = CreateAtlasFromFile(currPath+"airfactory.png", 0, 0, 64, 96, 32, 48, 1, false, true)
	buildingsAtlaces["refinery"] = CreateAtlasFromFile(currPath+"refinery.png", 0, 0, 96, 64, 48, 32, 1, false, true)
	buildingsAtlaces["depot"] = CreateAtlasFromFile(currPath+"depot.png", 0, 0, 48, 32, 48, 32, 1, false, true)
	buildingsAtlaces["silo"] = CreateAtlasFromFile(currPath+"silo.png", 0, 0, 32, 64, 16, 32, 1, false, true)
	buildingsAtlaces["turret_base"] = CreateAtlasFromFile(currPath+"turret_base.png", 0, 0, 16, 16, 16, 16, 1, false, true)
	turretsAtlaces["bld_turret_cannon"] = CreateDirectionalAtlasFromFile(currPath+"cannon_turret.png", 16, 16, 1, 2, true)
	turretsAtlaces["bld_turret_minigun"] = CreateDirectionalAtlasFromFile(currPath+"minigun_turret.png", 16, 16, 1, 2, true)
	buildingsAtlaces["fortress"] = CreateAtlasFromFile(currPath+"fortress.png", 0, 0, 32, 32, 32, 32, 1, false, true)
	turretsAtlaces["bld_fortress_cannon"] = CreateDirectionalAtlasFromFile(currPath+"fortress_turret.png", 32, 32, 1, 2, true)

	currPath = "resources/sprites/units/"
	unitChassisAtlaces["infantry"] = CreateDirectionalAtlasFromFile(currPath+"infantry.png", 16, 16, 1, 2, true)
	unitChassisAtlaces["tank"] = CreateDirectionalAtlasFromFile(currPath+"tank_chassis.png", 16, 16, 1, 2, true)
	turretsAtlaces["tank"] = CreateDirectionalAtlasFromFile(currPath+"tank_cannon.png", 16, 16, 1, 2, true)
	unitChassisAtlaces["quad"] = CreateDirectionalAtlasFromFile(currPath+"quad.png", 16, 16, 1, 2, true)
	turretsAtlaces["msltank"] = CreateDirectionalAtlasFromFile(currPath+"missiletank_turret.png", 16, 16, 1, 2, true)
	turretsAtlaces["aamsltank"] = CreateDirectionalAtlasFromFile(currPath+"aatank_turret.png", 16, 16, 1, 2, true)
	unitChassisAtlaces["harvester"] = CreateDirectionalAtlasFromFile(currPath+"harvester.png", 16, 16, 1, 2, true)

	currPath = "resources/sprites/units/aircrafts/"
	unitChassisAtlaces["air_gunship"] = CreateDirectionalAtlasFromFile(currPath+"combat_plane.png", 16, 16, 1, 2, true)
	unitChassisAtlaces["air_transport"] = CreateDirectionalAtlasFromFile(currPath+"transport_plane.png", 16, 16, 1, 2, true)

	currPath = "resources/sprites/projectiles/"
	projectilesAtlaces["shell"] = CreateDirectionalAtlasFromFile(currPath+"shell.png", 16, 8, 1, 2, false)
	projectilesAtlaces["bullets"] = CreateDirectionalAtlasFromFile(currPath+"bullets.png", 16, 8, 1, 2, false)
	projectilesAtlaces["missile"] = CreateDirectionalAtlasFromFile(currPath+"missile.png", 16, 8, 1, 2, false)
	projectilesAtlaces["aamissile"] = CreateDirectionalAtlasFromFile(currPath+"aamissile.png", 16, 8, 1, 2, false)

	currPath = "resources/sprites/ui/"
	uiAtlaces["factionflag"] = CreateAtlasFromFile(currPath+"building_faction_flag.png", 0, 0, 4, 4, 4, 4, 4, false, true)
	uiAtlaces["energyicon"] = CreateDirectionalAtlasFromFile(currPath+"energy_icon.png", 16, 8, 1, 1, false)
	uiAtlaces["readyicon"] = CreateDirectionalAtlasFromFile(currPath+"ready_icon.png", 16, 8, 1, 1, false)

	currPath = "resources/sprites/effects/"
	effectsAtlaces["smallexplosion"] = CreateAtlasFromFile(currPath+"explosion.png", 0, 0, 16, 16, 16, 16, 3, false, false)[0]
}

func extractSubimageFromImage(img image.Image, fromx, fromy, w, h int) image.Image {
	minx, miny := img.Bounds().Min.X, img.Bounds().Min.Y
	//maxx, maxy := img.Bounds().Min.X, img.Bounds().Max.Y
	switch img.(type) {
	case *image.RGBA:
		subImg := img.(*image.RGBA).SubImage(
			image.Rect(minx+fromx, miny+fromy, minx+fromx+w, miny+fromy+h),
		)
		// reset img bounds, because RayLib goes nuts about it otherwise
		subImg.(*image.RGBA).Rect = image.Rect(0, 0, w, h)
		return subImg
	case *image.NRGBA:
		subImg := img.(*image.NRGBA).SubImage(
			image.Rect(minx+fromx, miny+fromy, minx+fromx+w, miny+fromy+h),
		)
		// reset img bounds, because RayLib goes nuts about it otherwise
		subImg.(*image.NRGBA).Rect = image.Rect(0, 0, w, h)
		return subImg
	}
	panic("Unknown image type")
}

func CreateAtlasFromFile(filename string, topleftx, toplefty, originalSpriteW, originalSpriteH,
	desiredSpriteW, desiredSpriteH, totalFrames int, createAllDirections, createAllColors bool) []*spriteAtlas {

	atlases := make([]*spriteAtlas, 1)
	if createAllColors {
		atlases = make([]*spriteAtlas, len(factionColors))
	}

	file, _ := os.Open(filename)
	img, _ := png.Decode(file)
	file.Close()

	for i := range atlases {
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
			if createAllColors {
				replaceImageColorsToFactionImages(rlImg, i)
			}
			rl.ImageResizeNN(rlImg, int32(desiredSpriteW)*int32(SPRITE_SCALE_FACTOR), int32(desiredSpriteH)*int32(SPRITE_SCALE_FACTOR))
			newAtlas.atlas[0] = append(newAtlas.atlas[0], rl.LoadTextureFromImage(rlImg))
			if createAllDirections {
				for i := 1; i < 4; i++ {
					rl.ImageRotateCW(rlImg)
					newAtlas.atlas[i] = append(newAtlas.atlas[i], rl.LoadTextureFromImage(rlImg))
				}
			}
		}
		atlases[i] = &newAtlas
	}

	return atlases
}

func CreateDirectionalAtlasFromFile(filename string, originalSpriteSize, desiredSpriteSize, totalFrames, directionsInFile int, createAllColors bool) []*spriteAtlas {
	file, _ := os.Open(filename)
	img, _ := png.Decode(file)
	file.Close()

	atlases := make([]*spriteAtlas, 1)
	if createAllColors {
		atlases = make([]*spriteAtlas, len(factionColors))
	}

	for i := range atlases {

		newAtlas := spriteAtlas{
			// spriteSize: desiredSpriteSize * int(SPRITE_SCALE_FACTOR),
		}
		newAtlas.atlas = make([][]rl.Texture2D, 4*directionsInFile)

		for currFrame := 0; currFrame < totalFrames; currFrame++ {
			for currDirectionFromFile := 0; currDirectionFromFile < directionsInFile; currDirectionFromFile++ {
				currPic := extractSubimageFromImage(img, currFrame*originalSpriteSize, currDirectionFromFile*originalSpriteSize, originalSpriteSize, originalSpriteSize)
				rlImg := rl.NewImageFromImage(currPic)
				if createAllColors {
					replaceImageColorsToFactionImages(rlImg, i)
				}
				rl.ImageResizeNN(rlImg, int32(desiredSpriteSize)*int32(SPRITE_SCALE_FACTOR), int32(desiredSpriteSize)*int32(SPRITE_SCALE_FACTOR))
				newAtlas.atlas[currDirectionFromFile] = append(newAtlas.atlas[currDirectionFromFile], rl.LoadTextureFromImage(rlImg))
				for i := 1; i < 4; i++ {
					rl.ImageRotateCW(rlImg)
					newAtlas.atlas[i*directionsInFile+currDirectionFromFile] =
						append(newAtlas.atlas[i*directionsInFile+currDirectionFromFile], rl.LoadTextureFromImage(rlImg))
				}
			}
		}
		atlases[i] = &newAtlas
	}

	return atlases
}

func replaceImageColorsToFactionImages(img *rl.Image, factionColorNumber int) {
	rl.ImageColorReplace(img, color.RGBA{192, 192, 192, 255}, factionColors[factionColorNumber])
	rl.ImageColorReplace(img, color.RGBA{255, 0, 255, 255}, factionColors[factionColorNumber])
	darkerFactionTint := factionColors[factionColorNumber]
	darkerFactionTint.R /= 2
	darkerFactionTint.G /= 2
	darkerFactionTint.B /= 2
	rl.ImageColorReplace(img, color.RGBA{128, 128, 128, 255}, darkerFactionTint)
	rl.ImageColorReplace(img, color.RGBA{128, 0, 128, 255}, darkerFactionTint)
	darkestFactionTint := factionColors[factionColorNumber]
	darkestFactionTint.R /= 3
	darkestFactionTint.G /= 3
	darkestFactionTint.B /= 3
	rl.ImageColorReplace(img, color.RGBA{64, 64, 64, 255}, darkestFactionTint)
	rl.ImageColorReplace(img, color.RGBA{64, 0, 64, 255}, darkestFactionTint)
}
