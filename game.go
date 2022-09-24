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
	render      renderer
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

	timeLoopStarted := time.Now()
	timeCurrentActionStarted := time.Now()
	timeLogicStarted := time.Now()

	for !rl.WindowShouldClose() {
		timeReportString := fmt.Sprintf("Tick %d. ", g.battlefield.currentTick)
		timeLoopStarted = time.Now()
		timeCurrentActionStarted = time.Now()

		pc.playerControl(&g.battlefield)

		timeCurrentActionStarted = time.Now()
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
		timeReportString += g.createTimeReportString("AI", timeCurrentActionStarted, 1)

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
			timeReportString += g.createTimeReportString("units", timeCurrentActionStarted, 2)
		}
		if g.battlefield.currentTick%BUILDINGS_ACTIONS_TICK_EACH == 0 {
			timeCurrentActionStarted = time.Now()
			for i := g.battlefield.buildings.Front(); i != nil; i = i.Next() {
				g.battlefield.executeActionForActor(i.Value.(*building))
				if i.Value.(*building).turret != nil {
					g.battlefield.actorForActorsTurret(i.Value.(*building))
				}
			}
			timeReportString += g.createTimeReportString("buildings", timeCurrentActionStarted, 2)
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
			timeReportString += g.createTimeReportString("projectiles", timeCurrentActionStarted, 2)
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
		timeReportString += g.createTimeReportString("effects", timeCurrentActionStarted, 2)
		timeReportString += g.createTimeReportString("all actions", timeLogicStarted, 2)

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

			// reset faction tech
			for _, f := range g.battlefield.factions {
				f.currTechLevel = 0
				f.hasBuildings = map[buildingCode]bool{}
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
						bld.topLeftX, bld.topLeftY, bld.getStaticData().w, bld.getStaticData().h, 20,
					)
					g.battlefield.RandomlyAddEffectInTileRect(EFFECT_BIGGER_EXPLOSION, 50,
						bld.topLeftX, bld.topLeftY, bld.getStaticData().w, bld.getStaticData().h, 20,
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
					bld.faction.hasBuildings[bld.code] = true
					if bld.getStaticData().givesTechLevel > bld.faction.currTechLevel {
						bld.faction.currTechLevel = bld.getStaticData().givesTechLevel
					}
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
			timeReportString += g.createTimeReportString("unit orders", timeCurrentActionStarted, 2)
			//if g.battlefield.currentTick%(UNIT_ACTIONS_TICK_EACH*30) == 1 {
			//	debugWritef("Tick %d, orders logic: %dms\n", g.battlefield.currentTick, time.Since(timeCurrentActionStarted)/time.Millisecond)
			//}
		}
		if g.battlefield.currentTick%BUILDINGS_ACTIONS_TICK_EACH == 1 {
			timeCurrentActionStarted = time.Now()
			for i := g.battlefield.buildings.Front(); i != nil; i = i.Next() {
				g.battlefield.executeOrderForBuilding(i.Value.(*building))
			}
			timeReportString += g.createTimeReportString("blds orders", timeCurrentActionStarted, 2)
			//if g.battlefield.currentTick%(UNIT_ACTIONS_TICK_EACH*30) == 1 {
			//	debugWritef("Tick %d, bld orders logic: %dms\n", g.battlefield.currentTick, time.Since(timeCurrentActionStarted)/time.Millisecond)
			//}
		}

		g.battlefield.currentTick++

		timeReportString += g.createTimeReportString("cleanup+orders", timeLogicStarted, 2)

		timeReportString += g.createTimeReportString("whole tick", timeLoopStarted, 5)

		timeCurrentActionStarted = time.Now()
		if g.shouldTickBeRendered(g.battlefield.currentTick, RENDERER_DESIRED_FPS, DESIRED_TPS) {
			g.render.renderBattlefield(&g.battlefield, pc)
		}
		timeReportString += fmt.Sprintf("Render/sleep %dms", time.Since(timeCurrentActionStarted)/time.Millisecond)

		if (g.battlefield.currentTick-1)%10 == 0 {
			debugWrite(timeReportString)
			// debugWrite(g.battlefield.collectStatisticsForDebug())
		}
	}
}

func (g *game) shouldTickBeRendered(tick, fps, tps int) bool {
	//g.renderedFrames += RENDERER_DESIRED_FPS
	//if g.renderedFrames > DESIRED_TPS {
	//	g.renderedFrames = g.renderedFrames % DESIRED_TPS
	//	return true
	//}
	//return false
	return fps*tick/tps != fps*(tick+1)/tps
}

// returns string, also writes the string to renderer debug lines
func (g *game) createTimeReportString(actionName string, timeSince time.Time, criticalValueMs int) string {
	if len(g.render.timeDebugInfosToRender) == 0 {
		g.render.timeDebugInfosToRender = make([]debugTimeInfo, 0)
	}
	mcs := time.Since(timeSince) / time.Microsecond
	criticalMcs := time.Duration(criticalValueMs) * time.Microsecond
	if mcs > criticalMcs {
		// time.Sleep(1000 * time.Millisecond)
		debugWritef("WARNING: %s took %d mcs!\n", actionName, mcs)
	}

	neededFound := false
	for i := range g.render.timeDebugInfosToRender {
		if g.render.timeDebugInfosToRender[i].logicName == actionName {
			g.render.timeDebugInfosToRender[i].setNewValue(mcs)
			neededFound = true
			break
		}
	}
	if !neededFound {
		g.render.timeDebugInfosToRender = append(g.render.timeDebugInfosToRender, debugTimeInfo{
			logicName:           actionName,
			duration:            mcs,
			maxRecordedDuration: mcs,
			criticalDuration:    criticalMcs,
		})
	}

	return fmt.Sprintf("%s: %dmcs", actionName, mcs) + ", "
}

func (g *game) selectMapToGenerateBattlefield() {
	map_generator.SetRandom(&rnd)
	generatedMap := &map_generator.GeneratedMap{}
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
