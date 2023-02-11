//go:build dev
// +build dev

package main

func (c *yamlConfig) initFromFileOrCreate() {
	c.setDefaultValues()
}

func importUnitsDataOrCreateFile() {
	return
}

func importBuildingsDataOrCreateFile() {
	return
}
