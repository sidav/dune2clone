package main

import (
	"encoding/json"
	"errors"
	"os"
)

func importUnitsDataOrCreateFile() {
	const filePath = "units_data.json"
	fiBytes, err := os.ReadFile(filePath)
	if errors.Is(err, os.ErrNotExist) {

		// create new file with units data
		res, _ := json.MarshalIndent(sTableUnits, "", "\t")
		fo, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := fo.Close(); err != nil {
				panic(err)
			}
		}()
		fo.Write(res)
		return

	} else if err != nil {
		panic(err)
	}

	//nonPointerTable := make(map[string]unitStatic, 0)
	//err = json.Unmarshal(fiBytes, &nonPointerTable)
	//if err != nil {
	//	panic(err)
	//}

	sTableUnits = make(map[int]*unitStatic, 0)
	err = json.Unmarshal(fiBytes, &sTableUnits)
}

func importBuildingsDataOrCreateFile() {
	const filePath = "buildings_data.json"
	fiBytes, err := os.ReadFile(filePath)
	if errors.Is(err, os.ErrNotExist) {

		// create new file with units data
		res, _ := json.MarshalIndent(sTableBuildings, "", "\t")
		fo, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := fo.Close(); err != nil {
				panic(err)
			}
		}()
		fo.Write(res)
		return

	} else if err != nil {
		panic(err)
	}

	//nonPointerTable := make(map[string]unitStatic, 0)
	//err = json.Unmarshal(fiBytes, &nonPointerTable)
	//if err != nil {
	//	panic(err)
	//}

	sTableBuildings = make(map[buildingCode]*buildingStatic, 0)
	err = json.Unmarshal(fiBytes, &sTableBuildings)
}
