package main

import (
	"dune2clone/geometry"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (r *renderer) drawMinimap(b *battlefield, pc *playerController, maxW, maxH int32) {
	var tileSize = int(maxW / int32(len(b.tiles)))
	if maxH < maxW {
		tileSize = int(maxH / int32(len(b.tiles[0])))
	}
	w, h := int32(tileSize * len(b.tiles)), int32(tileSize * len(b.tiles[0]))
	posX, posY := WINDOW_W-w-2, WINDOW_H-h-2
	r.drawOutlinedRect(posX-2, posY-2, w+4, h+4, 2, pc.controlledFaction.getDarkerColor(), rl.DarkGray)
	// draw random noise if energy is insufficient
	if pc.controlledFaction.getAvailableEnergy() < 0 {
		for i := int32(0); i < 2*(w+h); i++ {
			nx := int32(rnd.Rand(int(w)))
			ny := int32(rnd.Rand(int(h)))
			size := int32(rnd.Rand(5) + 2)
			rl.DrawRectangle(posX+nx, posY+ny, size, size, rl.LightGray)
		}
	} else {
		for x := range b.tiles {
			for y := range b.tiles[x] {
				color := rl.Magenta
				if pc.controlledFaction.hasTileAtCoordsExplored(x, y) {
					switch b.tiles[x][y].code {
					case TILE_SAND:
						color = rl.Orange
						if b.tiles[x][y].resourcesAmount > 0 {
							color = rl.DarkPurple
						}
					case TILE_ROCK:
						color = rl.DarkBrown
					case TILE_BUILDABLE:
						color = rl.Brown
					default:
						color = rl.Magenta
					}
					if b.tiles[x][y].hasResourceVein {
						color = rl.Magenta
					}
					if !pc.controlledFaction.seesTileAtCoords(x, y) {
						color.R /= 3
						color.G /= 3
						color.B /= 3
					}
				} else {
					color = rl.Black
				}
				rl.DrawRectangle(posX+int32(x*tileSize), posY+int32(y*tileSize), int32(tileSize), int32(tileSize), color)
			}
		}
		// draw units and buildings TODO: optimize by reducing loop traversion?
		for i := b.buildings.Front(); i != nil; i = i.Next() {
			bld := i.Value.(*building)
			if b.hasFactionExploredBuilding(pc.controlledFaction, bld) {
				x, y, w, h := bld.topLeftX, bld.topLeftY, bld.getStaticData().w, bld.getStaticData().h
				rl.DrawRectangle(posX+int32(x*tileSize), posY+int32(y*tileSize), int32(w*tileSize), int32(h*tileSize), factionColors[bld.faction.colorNumber])
			}
		}
		for i := b.units.Front(); i != nil; i = i.Next() {
			unt := i.Value.(actor)
			if b.canFactionSeeActor(pc.controlledFaction, unt) {
				x, y := geometry.TrueCoordsToTileCoords(unt.getPhysicalCenterCoords())
				rl.DrawRectangle(posX+int32(x*tileSize), posY+int32(y*tileSize), int32(tileSize), int32(tileSize), factionColors[unt.getFaction().colorNumber])
			}

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
