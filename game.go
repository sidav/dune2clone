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
			tx, ty := geometry.TrueCoordsToTileCoords(i.Value.(actor).getPhysicalCenterCoords())
			pc.centerCameraAtTile(&g.battlefield, tx, ty)
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

		if g.battlefield.currentTick%RESOURCES_GROW_EACH_TICK == 0 {
			g.battlefield.performResourceGrowth()
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
				if !g.battlefield.areTileCoordsValid(tx, ty) || proj.setToRemove {
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
		// effects
		timeCurrentActionStarted = time.Now()
		for i := g.battlefield.effects.Front(); i != nil; i = i.Next() {
			eff := i.Value.(*effect)
			g.battlefield.actForEffect(eff)
			if eff.getExpirationPercent(g.battlefield.currentTick) > 100 {
				// debugWrite("Effect deleted.")
				// deleting while iterating
				setI := i
				if i.Prev() != nil {
					i = i.Prev()
				}
				g.battlefield.effects.Remove(setI)
			}
		}
		timeReportString += fmt.Sprintf("effects: %dms, ", time.Since(timeCurrentActionStarted)/time.Millisecond)
		timeReportString += fmt.Sprintf("all actions: %dms, ", time.Since(timeLogicStarted)/time.Millisecond)

		// cleanup and faction calculations
		if g.battlefield.currentTick%TRAVERSE_ALL_ACTORS_TICK_EACH == 0 {
			for _, f := range g.battlefield.factions {
				f.resetCurrents()
				f.cleanExpiredFactionDispatchRequests(g.battlefield.currentTick)
				f.resetVisibilityMaps(len(g.battlefield.tiles), len(g.battlefield.tiles[0]))
			}
			for i := g.battlefield.units.Front(); i != nil; i = i.Next() {
				unt := i.Value.(*unit)
				tx, ty := geometry.TrueCoordsToTileCoords(unt.getPhysicalCenterCoords())
				if !unt.isAlive() {
					setI := i
					if i.Prev() != nil {
						i = i.Prev()
					}
					g.battlefield.RandomlyAddEffectInTileRect(EFFECT_SMALL_EXPLOSION, 25,
						tx, ty, 1, 1, 5,
					)
					g.battlefield.units.Remove(setI)
				} else {
					unt.faction.exploreAround(tx, ty, 1, 1, unt.getVisionRange())
				}
			}

			for i := g.battlefield.buildings.Front(); i != nil; i = i.Next() {
				bld := i.Value.(*building)
				if !bld.isAlive() {
					// deleting while iterating
					setI := i
					if i.Prev() != nil {
						i = i.Prev()
					}
					g.battlefield.RandomlyAddEffectInTileRect(EFFECT_REGULAR_EXPLOSION, 50,
						bld.topLeftX, bld.topLeftY, bld.getStaticData().w, bld.getStaticData().h, 25,
					)
					g.battlefield.RandomlyAddEffectInTileRect(EFFECT_SMALL_EXPLOSION, 50,
						bld.topLeftX, bld.topLeftY, bld.getStaticData().w, bld.getStaticData().h, 25,
					)
					g.battlefield.changeTilesCodesInRectTo(
						bld.topLeftX, bld.topLeftY, bld.getStaticData().w, bld.getStaticData().h, TILE_BUILDABLE_DAMAGED,
					)
					if bld.unitPlacedInside != nil {
						g.battlefield.addActor(bld.unitPlacedInside)
						bld.unitPlacedInside = nil
					}
					g.battlefield.buildings.Remove(setI)
				} else {
					bld.faction.energyProduction += bld.getStaticData().givesEnergy
					bld.faction.energyConsumption += bld.getStaticData().consumesEnergy
					bld.faction.resourceStorage += bld.getStaticData().storageAmount
					bld.faction.exploreAround(bld.topLeftX, bld.topLeftY, bld.getStaticData().w, bld.getStaticData().h,
						bld.getVisionRange())
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
			//if g.battlefield.currentTick%(UNIT_ACTIONS_TICK_EACH*30) == 1 {
			//	debugWritef("Tick %d, orders logic: %dms\n", g.battlefield.currentTick, time.Since(timeCurrentActionStarted)/time.Millisecond)
			//}
		}
		if g.battlefield.currentTick%BUILDINGS_ACTIONS_TICK_EACH == 1 {
			timeCurrentActionStarted = time.Now()
			for i := g.battlefield.buildings.Front(); i != nil; i = i.Next() {
				g.battlefield.executeOrderForBuilding(i.Value.(*building))
			}
			timeReportString += fmt.Sprintf("bld orders: %dms, ", time.Since(timeCurrentActionStarted)/time.Millisecond)
			//if g.battlefield.currentTick%(UNIT_ACTIONS_TICK_EACH*30) == 1 {
			//	debugWritef("Tick %d, bld orders logic: %dms\n", g.battlefield.currentTick, time.Since(timeCurrentActionStarted)/time.Millisecond)
			//}
		}

		g.battlefield.currentTick++

		timeReportString += fmt.Sprintf("Whole logic: %dms, ", time.Since(timeLogicStarted)/time.Millisecond)

		timeReportString += fmt.Sprintf("whole tick took %dms", time.Since(timeLoopStarted)/time.Millisecond)
		if (g.battlefield.currentTick-1)%100 == 0 {
			debugWrite(timeReportString)
		}
	}
}

func (g *game) selectMapToGenerateBattlefield() {
	map_generator.SetRandom(&rnd)
	generatedMap := &map_generator.GameMap{}
	w := 64
	h := 64
	generatedMap.Generate(w, h)
	for {
		rl.BeginDrawing()
		drawGeneratedMap(generatedMap)
		rl.EndDrawing()
		time.Sleep(100 * time.Millisecond)
		if rl.IsKeyDown(rl.KeyEnter) || rl.IsKeyDown(rl.KeyEscape) {
			break
		} else if rl.IsKeyDown(rl.KeySpace) {
			generatedMap.Generate(w, h)
		} else if rl.IsKeyDown(rl.KeyRight) {
			w += 16
			h += 16
			generatedMap.Generate(w, h)
		} else if rl.IsKeyDown(rl.KeyLeft) {
			w -= 16
			h -= 16
			if w < 32 {
				w = 32
			}
			if h < 32 {
				h = 32
			}
			generatedMap.Generate(w, h)
		}
	}
	g.battlefield.initFromRandomMap(generatedMap)
}
