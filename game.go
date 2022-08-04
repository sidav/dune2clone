package main

import (
	"dune2clone/geometry"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type game struct {
	battlefield battlefield
}

func (g *game) startGame() {
	g.battlefield = battlefield{}
	g.battlefield.create(MAP_W, MAP_H)
	pc := &playerController{
		controlledFaction: g.battlefield.factions[0],
		selection:         nil,
	}
	r := renderer{}

	for !rl.WindowShouldClose() {
		r.renderBattlefield(&g.battlefield, pc)

		pc.playerControl(&g.battlefield)

		// execute actions
		if g.battlefield.currentTick%UNIT_ACTIONS_TICK_EACH == 0 {
			for i := g.battlefield.units.Front(); i != nil; i = i.Next() {
				g.battlefield.executeActionForActor(i.Value.(*unit))
				g.battlefield.actorForActorsTurret(i.Value.(*unit))
			}
		}
		if g.battlefield.currentTick%BUILDINGS_ACTIONS_TICK_EACH == 0 {
			for i := g.battlefield.buildings.Front(); i != nil; i = i.Next() {
				g.battlefield.executeActionForActor(i.Value.(*building))
				if i.Value.(*building).turret != nil {
					g.battlefield.actorForActorsTurret(i.Value.(*building))
				}
			}
		}
		if g.battlefield.currentTick%PROJECTILES_ACTIONS_TICK_EACH == 0 {
			for i := g.battlefield.projectiles.Front(); i != nil; i = i.Next() {
				proj := i.Value.(*projectile)
				g.battlefield.actForProjectile(proj)
				tx, ty := geometry.TrueCoordsToTileCoords(proj.centerX, proj.centerY)
				if !geometry.AreCoordsInTileRect(tx, ty, 0, 0, MAP_W, MAP_H) || proj.fuel <= 0 {
					// debugWrite("Projectile deleted.")
					g.battlefield.projectiles.Remove(i)
				}
			}
		}
		// execute orders
		if g.battlefield.currentTick%UNIT_ACTIONS_TICK_EACH == 1 {
			for i := g.battlefield.units.Front(); i != nil; i = i.Next() {
					g.battlefield.executeOrderForUnit(i.Value.(*unit))
			}
		}
		g.battlefield.currentTick++
	}
}
