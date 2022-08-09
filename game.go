package main

import (
	"dune2clone/geometry"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
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

	timeLoopStarted := time.Now()
	timeCurrentActionStarted := time.Now()
	timeLogicStarted := time.Now()

	for !rl.WindowShouldClose() {
		timeReportString := fmt.Sprintf("Tick %d. ", g.battlefield.currentTick)
		timeLoopStarted = time.Now()
		timeCurrentActionStarted = time.Now()
		r.renderBattlefield(&g.battlefield, pc)
		timeReportString += fmt.Sprintf("render: %dms, ", time.Since(timeCurrentActionStarted) / time.Millisecond)

		pc.playerControl(&g.battlefield)

		// execute actions
		timeLogicStarted = time.Now()
		if g.battlefield.currentTick%UNIT_ACTIONS_TICK_EACH == 0 {
			timeCurrentActionStarted = time.Now()
			for i := g.battlefield.units.Front(); i != nil; i = i.Next() {
				g.battlefield.executeActionForActor(i.Value.(*unit))
				g.battlefield.actorForActorsTurret(i.Value.(*unit))
			}
			timeReportString += fmt.Sprintf("units: %dms, ", time.Since(timeCurrentActionStarted) / time.Millisecond)
		}
		if g.battlefield.currentTick%BUILDINGS_ACTIONS_TICK_EACH == 0 {
			timeCurrentActionStarted = time.Now()
			for i := g.battlefield.buildings.Front(); i != nil; i = i.Next() {
				g.battlefield.executeActionForActor(i.Value.(*building))
				if i.Value.(*building).turret != nil {
					g.battlefield.actorForActorsTurret(i.Value.(*building))
				}
			}
			timeReportString += fmt.Sprintf("buildings: %dms, ", time.Since(timeCurrentActionStarted) / time.Millisecond)
		}
		if g.battlefield.currentTick%PROJECTILES_ACTIONS_TICK_EACH == 0 {
			timeCurrentActionStarted = time.Now()
			for i := g.battlefield.projectiles.Front(); i != nil; i = i.Next() {
				proj := i.Value.(*projectile)
				g.battlefield.actForProjectile(proj)
				tx, ty := geometry.TrueCoordsToTileCoords(proj.centerX, proj.centerY)
				if !geometry.AreCoordsInTileRect(tx, ty, 0, 0, MAP_W, MAP_H) || proj.fuel <= 0 {
					// debugWrite("Projectile deleted.")
					g.battlefield.projectiles.Remove(i)
				}
			}
			timeReportString += fmt.Sprintf("projs: %dms, ", time.Since(timeCurrentActionStarted) / time.Millisecond)
		}
		timeReportString += fmt.Sprintf("all actions: %dms, ", time.Since(timeLogicStarted) / time.Millisecond)

		// cleanup
		if g.battlefield.currentTick%UNIT_ACTIONS_TICK_EACH == 1 {
			for i := g.battlefield.units.Front(); i != nil; i = i.Next() {
				if i.Value.(*unit).currentHitpoints <= 0 {
					g.battlefield.units.Remove(i)
				}
			}
		}
		if g.battlefield.currentTick%BUILDINGS_ACTIONS_TICK_EACH == 1 {
			for i := g.battlefield.buildings.Front(); i != nil; i = i.Next() {
				if i.Value.(*building).currentHitpoints <= 0 {
					g.battlefield.buildings.Remove(i)
				}
			}
		}

		// execute orders
		if g.battlefield.currentTick%UNIT_ACTIONS_TICK_EACH == 1 {
			timeCurrentActionStarted = time.Now()
			for i := g.battlefield.units.Front(); i != nil; i = i.Next() {
					g.battlefield.executeOrderForUnit(i.Value.(*unit))
			}
			timeReportString += fmt.Sprintf("orders: %dms, ", time.Since(timeLogicStarted) / time.Millisecond)
		}
		g.battlefield.currentTick++

		timeReportString += fmt.Sprintf("Whole logic: %dms, ", time.Since(timeLogicStarted) / time.Millisecond)

		timeReportString += fmt.Sprintf("whole tick took %dms", time.Since(timeLoopStarted)/time.Millisecond)
		if (g.battlefield.currentTick-1) % 10 == 0 {
			debugWrite(timeReportString)
		}
	}
}
