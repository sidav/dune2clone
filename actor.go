package main

type actor interface {
	markSelected(bool)
	getCurrentAction() *action
	getName() string
	getFaction() *faction
}
