package main

func (b *battlefield) executeOrderForBuilding(bld *building) {
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
	if bld.currentAction.code != ACTION_BUILD {
		bld.currentAction.code = ACTION_BUILD
		bld.currentAction.targetActor = createBuilding(bld.currentOrder.targetActorCode, 0, 0, bld.faction)
	} else {
		if bld.currentAction.getCompletionPercent() >= 100 {
			bld.currentOrder.resetOrder()
			bld.currentOrder.code = ORDER_WAIT_FOR_BUILDING_PLACEMENT
		}
	}
}

func (b *battlefield) executePlaceBuildingOrder(bld *building) {
	whatIsBuilt := bld.currentAction.targetActor.(*building)
	if bld.currentOrder.targetTileX != -1 && bld.currentOrder.targetTileY != -1 {
		whatIsBuilt.topLeftX, whatIsBuilt.topLeftY = bld.currentOrder.targetTileX, bld.currentOrder.targetTileY
		b.addActor(whatIsBuilt)
		bld.currentAction.reset()
		bld.currentOrder.resetOrder()
	}
}

func (b *battlefield) executeProduceOrder(bld *building) {
	if bld.currentAction.code != ACTION_BUILD {
		bld.currentAction.code = ACTION_BUILD
		bld.currentAction.targetActor = createUnit(bld.currentOrder.targetActorCode, 0, 0, bld.faction)
		bld.currentOrder.resetOrder()
	}
}
