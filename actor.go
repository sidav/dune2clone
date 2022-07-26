package main

type actor interface {
	markSelected(bool)
	getName() string
}
