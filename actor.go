package main

type actor interface {
	markSelected(bool)
	getCurrentAction() *action
	getCurrentOrder() *order
	getName() string
	getHitpoints() int
	getMaxHitpoints() int
	getHitpointsPercentage() int
	getFaction() *faction
	getPhysicalCenterCoords() (float64, float64)
	getVisionRange() int
	isPresentAt(int, int) bool
	isInAir() bool
	isAlive() bool

	getExperience() int
	receiveExperienceAmount(int)
	getExperienceLevel() int
}
