package main

type actor interface {
	markSelected(bool)
	getCurrentAction() *action
	getName() string
	getFaction() *faction
	getPhysicalCenterCoords() (float64, float64)
	isPresentAt(int, int) bool
}
