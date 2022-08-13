package main

import (
	"dune2clone/geometry"
	"dune2clone/map_generator"
	rl "github.com/gen2brain/raylib-go/raylib"
	"strconv"
)

func drawGeneratedMap(gm *map_generator.GameMap) {
	rl.BeginDrawing()
	rl.DrawText("Select map. SPACE to generate new, ENTER to select current.", 0, 0, 28, rl.White)
	const tileSize = 10
	rl.ClearBackground(rl.Black)
	for x := range gm.Tiles {
		for y := range gm.Tiles[x] {
			color := rl.Magenta
			switch gm.Tiles[x][y] {
			case map_generator.SAND:
				color = rl.Orange
			case map_generator.RESOURCES:
				color = rl.Red
			default:
				color = rl.Brown
			}
			rl.DrawRectangle(int32(x*tileSize), 32+int32(y*tileSize), tileSize, tileSize, color)
		}
	}
	for sp := range gm.StartPoints {
		const spSize = tileSize * 4
		rl.DrawRectangle(int32(tileSize*gm.StartPoints[sp][0])-spSize/3, int32(tileSize*gm.StartPoints[sp][1]), spSize, spSize, rl.Black)
		rl.DrawText(strconv.Itoa(sp+1), int32(tileSize*gm.StartPoints[sp][0]), int32(tileSize*gm.StartPoints[sp][1]), spSize, factionColors[sp])
	}

	rl.EndDrawing()
}

func (r *renderer) drawMinimap(b *battlefield, pc *playerController, posX, posY int32, w, h int32) {
	var tileSize = int(w) / len(b.tiles)
	if h > w {
		tileSize = int(w) / len(b.tiles)
	}
	// draw random noise if energy is insufficient
	if pc.controlledFaction.getAvailableEnergy() < 0 {
		r.drawOutlinedRect(posX-2, posY-2, w, h, 2, rl.Green, rl.DarkGray)
		for i := int32(0); i < 2*(w+h); i++ {
			nx := int32(rnd.Rand(int(w)))
			ny := int32(rnd.Rand(int(h)))
			size := int32(rnd.Rand(5) + 2)
			rl.DrawRectangle(posX+nx, posY+ny, size, size, rl.LightGray)
		}
	} else {
		r.drawOutlinedRect(posX-2, posY-2, w, h, 2, rl.Green, rl.Black)
		for x := range b.tiles {
			for y := range b.tiles[x] {
				color := rl.Magenta
				switch b.tiles[x][y].code {
				case TILE_SAND:
					color = rl.Orange
					if b.tiles[x][y].resourcesAmount > 0 {
						color = rl.Red
					}
				case TILE_BUILDABLE:
					color = rl.Brown
				default:
					color = rl.Magenta
				}
				rl.DrawRectangle(posX+int32(x*tileSize), posY+int32(y*tileSize), int32(tileSize), int32(tileSize), color)
			}
		}
		// draw units and buildings TODO: optimize by reducing loop traversion?
		for i := b.buildings.Front(); i != nil; i = i.Next() {
			bld := i.Value.(*building)
			x, y, w, h := bld.topLeftX, bld.topLeftY, bld.getStaticData().w, bld.getStaticData().h
			rl.DrawRectangle(posX+int32(x*tileSize), posY+int32(y*tileSize), int32(w*tileSize), int32(h*tileSize), factionColors[bld.faction.colorNumber])
		}
		for i := b.units.Front(); i != nil; i = i.Next() {
			unt := i.Value.(actor)
			x, y := geometry.TrueCoordsToTileCoords(unt.getPhysicalCenterCoords())
			rl.DrawRectangle(posX+int32(x*tileSize), posY+int32(y*tileSize), int32(tileSize), int32(tileSize), factionColors[unt.getFaction().colorNumber])

		}
	}
	// now let's draw current screen rectangle
	scrRectX := int32(tileSize) * r.camTopLeftX / TILE_SIZE_IN_PIXELS
	scrRectY := int32(tileSize) * r.camTopLeftY / TILE_SIZE_IN_PIXELS
	scrRectW := int32(tileSize) * WINDOW_W / TILE_SIZE_IN_PIXELS
	scrRectH := int32(tileSize) * WINDOW_H / TILE_SIZE_IN_PIXELS
	rl.DrawRectangleLines(posX+scrRectX, posY+scrRectY, scrRectW, scrRectH, rl.White)
	rl.DrawRectangleLines(posX+scrRectX+1, posY+scrRectY+1, scrRectW-2, scrRectH-2, rl.White)
}
