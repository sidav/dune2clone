package main

import (
	"dune2clone/geometry"
	"dune2clone/map_generator"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
)

type game struct {
	battlefield battlefield
}

func (g *game) startGame() {
	g.selectMapToGenerateBattlefield()

	pc := &playerController{
		controlledFaction: g.battlefield.factions[0],
		selection:         nil,
	}
	for i := g.battlefield.buildings.Front(); i != nil; i = i.Next() {
		if i.Value.(actor).getFaction() == pc.controlledFaction {
			pc.centerCameraAtTile(geometry.TrueCoordsToTileCoords(i.Value.(actor).getPhysicalCenterCoords()))
		}
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
		timeReportString += fmt.Sprintf("render: %dms, ", time.Since(timeCurrentActionStarted)/time.Millisecond)

		pc.playerControl(&g.battlefield)

		if g.battlefield.currentTick%AI_ANALYZES_EACH == 0 {
			for i := range g.battlefield.ais {
				g.battlefield.ais[i].aiAnalyze(&g.battlefield)
			}
		}
		if g.battlefield.currentTick%AI_ACTS_EACH == 0 {
			for i := range g.battlefield.ais {
				g.battlefield.ais[i].aiControl(&g.battlefield)
			}
		}

		// execute actions
		timeLogicStarted = time.Now()
		if g.battlefield.currentTick%UNIT_ACTIONS_TICK_EACH == 0 {
			timeCurrentActionStarted = time.Now()
			for i := g.battlefield.units.Front(); i != nil; i = i.Next() {
				g.battlefield.executeActionForActor(i.Value.(*unit))
				g.battlefield.actorForActorsTurret(i.Value.(*unit))
			}
			timeReportString += fmt.Sprintf("units: %dms, ", time.Since(timeCurrentActionStarted)/time.Millisecond)
		}
		if g.battlefield.currentTick%BUILDINGS_ACTIONS_TICK_EACH == 0 {
			timeCurrentActionStarted = time.Now()
			for i := g.battlefield.buildings.Front(); i != nil; i = i.Next() {
				g.battlefield.executeActionForActor(i.Value.(*building))
				if i.Value.(*building).turret != nil {
					g.battlefield.actorForActorsTurret(i.Value.(*building))
				}
			}
			timeReportString += fmt.Sprintf("buildings: %dms, ", time.Since(timeCurrentActionStarted)/time.Millisecond)
		}
		if g.battlefield.currentTick%PROJECTILES_ACTIONS_TICK_EACH == 0 {
			timeCurrentActionStarted = time.Now()
			for i := g.battlefield.projectiles.Front(); i != nil; i = i.Next() {
				proj := i.Value.(*projectile)
				g.battlefield.actForProjectile(proj)
				tx, ty := geometry.TrueCoordsToTileCoords(proj.centerX, proj.centerY)
				if !geometry.AreCoordsInTileRect(tx, ty, 0, 0, MAP_W, MAP_H) || proj.setToRemove {
					// debugWrite("Projectile deleted.")
					// deleting while iterating
					setI := i
					if i.Prev() != nil {
						i = i.Prev()
					}
					g.battlefield.projectiles.Remove(setI)
				}
			}
			timeReportString += fmt.Sprintf("projs: %dms, ", time.Since(timeCurrentActionStarted)/time.Millisecond)
		}
		timeReportString += fmt.Sprintf("all actions: %dms, ", time.Since(timeLogicStarted)/time.Millisecond)

		// cleanup
		if g.battlefield.currentTick%UNIT_ACTIONS_TICK_EACH == 1 {
			for i := g.battlefield.units.Front(); i != nil; i = i.Next() {
				if i.Value.(*unit).currentHitpoints <= 0 {
					setI := i
					if i.Prev() != nil {
						i = i.Prev()
					}
					g.battlefield.units.Remove(setI)
				}
			}
		}
		if g.battlefield.currentTick%BUILDINGS_ACTIONS_TICK_EACH == 1 {
			for i := g.battlefield.buildings.Front(); i != nil; i = i.Next() {
				if i.Value.(*building).currentHitpoints <= 0 {
					// deleting while iterating
					setI := i
					if i.Prev() != nil {
						i = i.Prev()
					}
					g.battlefield.buildings.Remove(setI)
				}
			}
		}

		// execute orders
		if g.battlefield.currentTick%UNIT_ACTIONS_TICK_EACH == 1 {
			timeCurrentActionStarted = time.Now()
			for i := g.battlefield.units.Front(); i != nil; i = i.Next() {
				g.battlefield.executeOrderForUnit(i.Value.(*unit))
			}
			timeReportString += fmt.Sprintf("orders: %dms, ", time.Since(timeCurrentActionStarted)/time.Millisecond)
			if g.battlefield.currentTick%(UNIT_ACTIONS_TICK_EACH*30) == 1 {
				debugWritef("Tick %d, orders logic: %dms\n", g.battlefield.currentTick, time.Since(timeCurrentActionStarted)/time.Millisecond)
			}
		}
		if g.battlefield.currentTick%BUILDINGS_ACTIONS_TICK_EACH == 1 {
			timeCurrentActionStarted = time.Now()
			for i := g.battlefield.buildings.Front(); i != nil; i = i.Next() {
				g.battlefield.executeOrderForBuilding(i.Value.(*building))
			}
			timeReportString += fmt.Sprintf("bld orders: %dms, ", time.Since(timeCurrentActionStarted)/time.Millisecond)
			if g.battlefield.currentTick%(UNIT_ACTIONS_TICK_EACH*30) == 1 {
				debugWritef("Tick %d, bld orders logic: %dms\n", g.battlefield.currentTick, time.Since(timeCurrentActionStarted)/time.Millisecond)
			}
		}


		g.battlefield.currentTick++

		timeReportString += fmt.Sprintf("Whole logic: %dms, ", time.Since(timeLogicStarted)/time.Millisecond)

		timeReportString += fmt.Sprintf("whole tick took %dms", time.Since(timeLoopStarted)/time.Millisecond)
		if (g.battlefield.currentTick-1)%10 == 0 {
			debugWrite(timeReportString)
		}
	}
}

func (g *game) selectMapToGenerateBattlefield() {
	map_generator.SetRandom(&rnd)
	generatedMap := &map_generator.GameMap{}
	generatedMap.Init(64, 64)
	generatedMap.Generate()
	for {
		rl.BeginDrawing()
		rl.EndDrawing()
		drawGeneratedMap(generatedMap)
		if rl.IsKeyDown(rl.KeyEnter) || rl.IsKeyDown(rl.KeyEscape) {
			break
		} else if rl.GetKeyPressed() != 0 {
			generatedMap.Init(64, 64)
			generatedMap.Generate()
		}
	}
	g.battlefield.initFromRandomMap(generatedMap)
}
