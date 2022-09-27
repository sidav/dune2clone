package main

import (
	"dune2clone/geometry"
	"fmt"
)

type aiStruct struct {
	name                             string
	personalityName                  string
	controlsFaction                  *faction
	moneyPoorMax                     float64
	moneyRichMin                     float64
	current                          aiAnalytics
	desired                          aiAnalytics
	max                              aiAnalytics
	currBaseCenterX, currBaseCenterY int

	alreadyOrderedBuildThisTick bool
	taskForces                  []*aiTaskForce
}

func (ai *aiStruct) isPoor() bool {
	return ai.controlsFaction.getMoney() <= ai.moneyPoorMax
}

func (ai *aiStruct) isRich() bool {
	return ai.controlsFaction.getMoney() > ai.moneyRichMin
}

func (ai *aiStruct) debugWritef(msg string, args ...interface{}) {
	debugWritef("%s [%s]: %s", ai.name, ai.personalityName, fmt.Sprintf(msg, args...))
}

func (ai *aiStruct) isActorInRangeFromBase(a actor, r int) bool {
	ttx, tty := geometry.TrueCoordsToTileCoords(a.getPhysicalCenterCoords())
	return geometry.GetApproxDistFromTo(ai.currBaseCenterX, ai.currBaseCenterY, ttx, tty) <= r
}

func (ai *aiStruct) aiControl(b *battlefield) {
	ai.debugWritef("acts.\n")
	ai.alreadyOrderedBuildThisTick = false
	for i := b.buildings.Front(); i != nil; i = i.Next() {
		if bld, ok := i.Value.(*building); ok {
			if bld.getFaction() == ai.controlsFaction {
				ai.actForBuilding(b, bld)
			}
		}
	}
	for i := b.units.Front(); i != nil; i = i.Next() {
		if unt, ok := i.Value.(*unit); ok {
			if unt.getFaction() == ai.controlsFaction {
				ai.actForUnit(b, unt)
			}
		}
	}
	ai.giveOrdersToAllTaskForces(b)
}
