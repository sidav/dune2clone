package main

func (b *battlefield) executeOrderForBuilding(bld *building) {
	if bld.isRepairingSelf {
		b.executeBuildingSelfRepair(bld) // TODO: maybe move to actions?..
	}
	switch bld.currentOrder.code {
	case ORDER_BUILD:
		b.executeBuildOrder(bld)
	case ORDER_WAIT_FOR_BUILDING_PLACEMENT:
		b.executePlaceBuildingOrder(bld)
	case ORDER_PRODUCE:
		b.executeProduceOrder(bld)
	}
}

func (b *battlefield) executeBuildOrder(bld *building) {
	if !bld.faction.isTechAvailableForBuildingOfCode(buildingCode(bld.currentOrder.targetActorCode)) {
		debugWritef("Tech requirements are ignored: bld %s at (%d, %d)", bld.getName(), bld.topLeftX, bld.topLeftY)
	}
	switch bld.getStaticData().BuildType {
	case BTYPE_BUILD_FIRST:
		if bld.currentAction.code != ACTION_BUILD {
			bld.currentAction.code = ACTION_BUILD
			tBld := createBuilding(buildingCode(bld.currentOrder.targetActorCode), 0, 0, bld.faction)
			bld.currentAction.targetActor = tBld
			bld.currentOrder.targetActor = tBld
		} else {
			if bld.currentAction.getCompletionPercent() >= 100 {
				tBld := bld.currentOrder.targetActor
				bld.currentOrder.resetOrder()
				bld.currentOrder.targetActor = tBld
				bld.currentOrder.code = ORDER_WAIT_FOR_BUILDING_PLACEMENT
			}
		}
	case BTYPE_PLACE_FIRST:
		if bld.currentOrder.targetActor == nil {
			bld.currentOrder.resetOrder()
			creatingBuilding := createBuilding(buildingCode(bld.currentOrder.targetActorCode), 0, 0, bld.faction)
			creatingBuilding.currentHitpoints = creatingBuilding.getMaxHitpoints()/100 + creatingBuilding.getMaxHitpoints()%100
			bld.currentOrder.targetActor = creatingBuilding
			bld.currentOrder.code = ORDER_WAIT_FOR_BUILDING_PLACEMENT
			return
		} else if bld.currentAction.getCompletionPercent() >= 100 {
			bld.currentOrder.resetOrder()
			bld.currentAction.resetAction()
			return
		}
	}
}

func (b *battlefield) executePlaceBuildingOrder(bld *building) {
	switch bld.getStaticData().BuildType {
	case BTYPE_BUILD_FIRST:
		whatIsBuilt := bld.currentAction.targetActor.(*building)
		if bld.currentOrder.targetTileX != -1 && bld.currentOrder.targetTileY != -1 {
			whatIsBuilt.topLeftX, whatIsBuilt.topLeftY = bld.currentOrder.targetTileX, bld.currentOrder.targetTileY
			whatIsBuilt.currentAction.code = ACTION_BEING_BUILT
			whatIsBuilt.currentAction.builtAs = BTYPE_BUILD_FIRST
			whatIsBuilt.currentAction.maxCompletionAmount = float64(config.Engine.BuildingAnimationTicks)
			b.addActor(whatIsBuilt)
			bld.currentAction.resetAction()
			bld.currentOrder.resetOrder()
		}
	case BTYPE_PLACE_FIRST:
		whatIsBuilt := bld.currentOrder.targetActor.(*building)
		if bld.currentOrder.targetTileX != -1 && bld.currentOrder.targetTileY != -1 {
			whatIsBuilt.topLeftX, whatIsBuilt.topLeftY = bld.currentOrder.targetTileX, bld.currentOrder.targetTileY
			whatIsBuilt.currentAction.code = ACTION_BEING_BUILT
			whatIsBuilt.currentAction.builtAs = BTYPE_PLACE_FIRST
			whatIsBuilt.currentAction.maxCompletionAmount = float64(whatIsBuilt.getStaticData().BuildTime *
				(config.Engine.TicksPerNominalSecond / config.Engine.BuildingsActionPeriod))
			b.addActor(whatIsBuilt)
			bld.currentAction.code = ACTION_BUILD
			bld.currentAction.targetActor = bld.currentOrder.targetActor
			bld.currentOrder.code = ORDER_BUILD
		}
	}
}

func (b *battlefield) executeProduceOrder(bld *building) {
	if !bld.faction.isTechAvailableForUnitOfCode(bld.currentOrder.targetActorCode) {
		debugWritef("Tech requirements are ignored: %s at (%d, %d)", bld.getName(), bld.topLeftX, bld.topLeftY)
	}
	if bld.currentAction.code != ACTION_BUILD {
		bld.currentAction.code = ACTION_BUILD
		bld.currentAction.targetActor = createUnit(bld.currentOrder.targetActorCode, 0, 0, bld.faction)
		bld.currentOrder.resetOrder()
	}
}
