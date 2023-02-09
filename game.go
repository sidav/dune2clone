package main

import (
	"container/list"
	"dune2clone/geometry"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"time"
)

type game struct {
	gameIsRunning bool
	battlefield   battlefield
	render        *renderer
}

func (g *game) startGame() {
	g.gameIsRunning = true
	g.selectMapToGenerateBattlefield()
	if !g.gameIsRunning {
		return
	}

	pc := &playerController{
		controlledFaction: g.battlefield.factions[0],
		selection:         nil,
	}
	g.centerPlayerCameraAtStartPosition(pc)

	timeLoopStarted := time.Now()
	timeCurrentActionStarted := time.Now()
	timeLogicStarted := time.Now()

	for !rl.WindowShouldClose() && g.gameIsRunning {
		timeReportString := fmt.Sprintf("Tick %d. ", g.battlefield.currentTick)
		timeLoopStarted = time.Now()
		timeCurrentActionStarted = time.Now()

		pc.playerControl(&g.battlefield)

		timeCurrentActionStarted = time.Now()
		g.performAiActions()
		timeReportString += g.createTimeReportString("AI", timeCurrentActionStarted, 1)

		if g.battlefield.currentTick%config.Economy.ResourcesGrowthPeriod == 0 {
			g.battlefield.performResourceGrowth()
		}

		// execute actions
		timeLogicStarted = time.Now()
		if g.battlefield.currentTick%config.Engine.UnitsActionPeriod == 0 {
			timeCurrentActionStarted = time.Now()
			for i := g.battlefield.units.Front(); i != nil; i = i.Next() {
				g.battlefield.executeActionForActor(i.Value.(*unit))
				g.battlefield.actorForActorsTurret(i.Value.(*unit))
			}
			timeReportString += g.createTimeReportString("units", timeCurrentActionStarted, 2)
		}
		if g.battlefield.currentTick%config.Engine.BuildingsActionPeriod == 0 {
			timeCurrentActionStarted = time.Now()
			for i := g.battlefield.buildings.Front(); i != nil; i = i.Next() {
				g.battlefield.executeActionForActor(i.Value.(*building))
				if i.Value.(*building).turret != nil {
					g.battlefield.actorForActorsTurret(i.Value.(*building))
				}
			}
			timeReportString += g.createTimeReportString("buildings", timeCurrentActionStarted, 2)
		}
		if g.battlefield.currentTick%config.Engine.ProjectilesActionPeriod == 0 {
			timeCurrentActionStarted = time.Now()
			// "next" is for deletion while iterating
			var next *list.Element
			for i := g.battlefield.projectiles.Front(); i != nil; i = next {
				next = i.Next()
				proj := i.Value.(*projectile)
				g.battlefield.actForProjectile(proj)
				tx, ty := geometry.TrueCoordsToTileCoords(proj.centerX, proj.centerY)
				if !g.battlefield.areTileCoordsValid(tx, ty) || proj.setToRemove {
					// debugWrite("Projectile deleted.")
					g.battlefield.projectiles.Remove(i)
				}
			}
			timeReportString += g.createTimeReportString("projectiles", timeCurrentActionStarted, 2)
		}
		// effects
		timeCurrentActionStarted = time.Now()
		var next *list.Element
		for i := g.battlefield.effects.Front(); i != nil; i = next {
			next = i.Next()
			eff := i.Value.(*effect)
			g.battlefield.actForEffect(eff)
			if eff.getExpirationPercent(g.battlefield.currentTick) > 100 {
				// deleting while iterating
				g.battlefield.effects.Remove(i)
			}
		}
		timeReportString += g.createTimeReportString("effects", timeCurrentActionStarted, 2)
		timeReportString += g.createTimeReportString("all actions", timeLogicStarted, 2)

		// cleanup and faction calculations
		if g.battlefield.currentTick%config.Engine.CleanupPeriod == 0 {
			g.traverseAllActors()
		}

		// execute orders
		if g.battlefield.currentTick%config.Engine.UnitsListenOrderPeriod == 0 {
			timeCurrentActionStarted = time.Now()
			for i := g.battlefield.units.Front(); i != nil; i = i.Next() {
				g.battlefield.executeOrderForUnit(i.Value.(*unit))
			}
			timeReportString += g.createTimeReportString("unit orders", timeCurrentActionStarted, 2)
			//if g.battlefield.currentTick%(config.Engine.UnitActionPeriod*30) == 1 {
			//	debugWritef("Tick %d, orders logic: %dms\n", g.battlefield.currentTick, time.Since(timeCurrentActionStarted)/time.Millisecond)
			//}
		}
		if g.battlefield.currentTick%config.Engine.BuildingListenOrderPeriod == 0 {
			timeCurrentActionStarted = time.Now()
			for i := g.battlefield.buildings.Front(); i != nil; i = i.Next() {
				g.battlefield.executeOrderForBuilding(i.Value.(*building))
			}
			timeReportString += g.createTimeReportString("blds orders", timeCurrentActionStarted, 2)
			//if g.battlefield.currentTick%(config.Engine.UnitActionPeriod*30) == 1 {
			//	debugWritef("Tick %d, bld orders logic: %dms\n", g.battlefield.currentTick, time.Since(timeCurrentActionStarted)/time.Millisecond)
			//}
		}

		g.battlefield.currentTick++

		timeReportString += g.createTimeReportString("cleanup+orders", timeLogicStarted, 2)

		if g.battlefield.currentTick > config.TargetTPS {
			timeReportString += g.createTimeReportString("whole tick", timeLoopStarted, 15)
		}

		timeCurrentActionStarted = time.Now()
		if g.shouldTickBeRendered(g.battlefield.currentTick, config.TargetFPS, config.TargetTPS) {
			g.render.renderBattlefield(&g.battlefield, pc)
		}
		timeReportString += fmt.Sprintf("Render/sleep %dms", time.Since(timeCurrentActionStarted)/time.Millisecond)

		// 887 is just a bug enough prime
		if (g.battlefield.currentTick-1)%887 == 0 {
			debugWrite(timeReportString)
			// debugWrite(g.battlefield.collectStatisticsForDebug())
		}

		if (g.battlefield.currentTick)%300 == 0 {
			g.checkGameEnd()
		}
	}
}

func (g *game) performAiActions() {
	for i := range g.battlefield.ais {
		if g.battlefield.currentTick%config.AiSettings.AiAnalyzePeriod == i {
			g.battlefield.ais[i].aiAnalyze(&g.battlefield)
		}
	}
	for i := range g.battlefield.ais {
		if g.battlefield.currentTick%config.AiSettings.AiActPeriod == i {
			g.battlefield.ais[i].aiControl(&g.battlefield)
		}
	}
}

func (g *game) traverseAllActors() {
	for _, f := range g.battlefield.factions {
		f.resetCurrents()
		f.cleanExpiredFactionDispatchRequests(g.battlefield.currentTick)
		f.resetVisibilityMaps(len(g.battlefield.tiles), len(g.battlefield.tiles[0]))
	}
	var next *list.Element // for deletion while iterating the list
	for i := g.battlefield.units.Front(); i != nil; i = next {
		next = i.Next()
		unt := i.Value.(*unit)
		tx, ty := geometry.TrueCoordsToTileCoords(unt.getPhysicalCenterCoords())
		if !unt.isAlive() {
			// for deletion while iterating
			setI := i
			if i.Prev() != nil {
				i = i.Prev()
			}
			g.battlefield.RandomlyAddEffectInTileRect(EFFECT_SMALL_EXPLOSION, 25,
				tx, ty, 1, 1, 5,
			)
			g.battlefield.removeActor(setI.Value.(*unit))
		} else {
			unt.faction.exploreAround(tx, ty, 1, 1,
				modifyVisionRangeByUnitExpLevel(unt.getVisionRange(), unt.getExperienceLevel()))
			if g.battlefield.currentTick%config.Engine.RegenHpPeriod == 0 {
				unt.receiveHealing(unt.getRegenAmount())
			}
		}
	}

	for i := g.battlefield.buildings.Front(); i != nil; i = next {
		next = i.Next()
		bld := i.Value.(*building)
		if !bld.isAlive() {
			// for deletion while iterating
			if i.Prev() != nil {
				i = i.Prev()
			}
			g.battlefield.RandomlyAddEffectInTileRect(EFFECT_REGULAR_EXPLOSION, 50,
				bld.topLeftX, bld.topLeftY, bld.getStaticData().W, bld.getStaticData().H, 20,
			)
			g.battlefield.RandomlyAddEffectInTileRect(EFFECT_BIGGER_EXPLOSION, 50,
				bld.topLeftX, bld.topLeftY, bld.getStaticData().W, bld.getStaticData().H, 20,
			)
			g.battlefield.changeTilesCodesInRectTo(
				bld.topLeftX, bld.topLeftY, bld.getStaticData().W, bld.getStaticData().H, TILE_BUILDABLE_DAMAGED,
			)
			if bld.unitPlacedInside != nil {
				g.battlefield.addActor(bld.unitPlacedInside)
				bld.unitPlacedInside = nil
			}
			if bld.currentAction.code == ACTION_BUILD {
				bld.currentAction.targetActor.setHitpoints(0)
			}
			g.battlefield.removeActor(bld)
		} else {
			bld.faction.hasBuildings[bld.code] = true
			if bld.getStaticData().GivesTechLevel > bld.faction.currTechLevel {
				bld.faction.currTechLevel = bld.getStaticData().GivesTechLevel
			}
			if !bld.isUnderConstruction() {
				bld.faction.energyProduction += bld.getStaticData().GivesEnergy
				bld.faction.energyConsumption += bld.getStaticData().ConsumesEnergy
				bld.faction.increaseResourcesStorage(bld.getStaticData().StorageAmount)
			}
			bld.faction.exploreAround(bld.topLeftX, bld.topLeftY, bld.getStaticData().W, bld.getStaticData().H,
				modifyVisionRangeByUnitExpLevel(bld.getVisionRange(), bld.getExperienceLevel()))
			if g.battlefield.currentTick%config.Engine.RegenHpPeriod == 0 {
				bld.receiveHealing(bld.getRegenAmount())
			}
		}
	}
}

func (g *game) checkGameEnd() {
	var remainingFaction *faction
	moreThanOneFactionRemains := false
	playerIsDefeated := true
	for b := g.battlefield.buildings.Front(); b != nil; b = b.Next() {
		if remainingFaction == nil {
			remainingFaction = b.Value.(actor).getFaction()
		} else {
			if b.Value.(actor).getFaction() != remainingFaction {
				moreThanOneFactionRemains = true
			}
		}
		if b.Value.(actor).getFaction() == g.battlefield.factions[0] {
			playerIsDefeated = false
			if moreThanOneFactionRemains {
				break
			}
		}
	}
	if playerIsDefeated {
		// check units then
		for u := g.battlefield.units.Front(); u != nil; u = u.Next() {
			if u.Value.(actor).getFaction() == g.battlefield.factions[0] {
				playerIsDefeated = false
				moreThanOneFactionRemains = true
				break
			}
		}
	}
	if playerIsDefeated {
		if moreThanOneFactionRemains {
			g.render.drawEndGame(g.battlefield.factions, nil, true)
		} else {
			g.render.drawEndGame(g.battlefield.factions, remainingFaction, true)
		}
	} else if !moreThanOneFactionRemains {
		g.render.drawEndGame(g.battlefield.factions, remainingFaction, false)
	} else {
		return
	}
	g.gameIsRunning = false
}

func (g *game) centerPlayerCameraAtStartPosition(pc *playerController) {
	cameraSet := false
	for i := g.battlefield.units.Front(); i != nil; i = i.Next() {
		if i.Value.(actor).getFaction() == pc.controlledFaction {
			tx, ty := i.Value.(*unit).getTileCoords()
			pc.centerCameraAtTile(&g.battlefield, tx, ty)
			cameraSet = true
		}
	}
	if !cameraSet {
		for i := g.battlefield.buildings.Front(); i != nil; i = i.Next() {
			if i.Value.(actor).getFaction() == pc.controlledFaction {
				tx, ty := geometry.TrueCoordsToTileCoords(i.Value.(actor).getPhysicalCenterCoords())
				pc.centerCameraAtTile(&g.battlefield, tx, ty)
				cameraSet = true
			}
		}
	}
}

func (g *game) shouldTickBeRendered(tick, fps, tps int) bool {
	return fps*tick/tps != fps*(tick+1)/tps
}

// returns string, also writes the string to renderer debug lines
func (g *game) createTimeReportString(actionName string, timeSince time.Time, criticalValueMs int) string {
	if len(g.render.timeDebugInfosToRender) == 0 {
		g.render.timeDebugInfosToRender = make([]debugTimeInfo, 0)
	}
	mcs := time.Since(timeSince) / time.Microsecond
	criticalMcs := time.Duration(criticalValueMs) * time.Microsecond
	//if mcs > criticalMcs {
	//	// time.Sleep(1000 * time.Millisecond)
	//	debugWritef("WARNING: %s took %d mcs!\n", actionName, mcs)
	//}

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
