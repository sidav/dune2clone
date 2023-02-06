package main

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
)

type yamlConfig struct {
	DebugOutput bool `yaml:"debug_output"`
	TargetFPS   int  `yaml:"target_fps"`
	TargetTPS   int  `yaml:"target_tps"`

	Engine struct {
		UnitActionPeriod        int `yaml:"unit_action_period"`
		UnitListenOrderPeriod   int `yaml:"unit_listen_order_period"`
		BuildingActionPeriod    int `yaml:"building_action_period"`
		BuildingOrderPeriod     int `yaml:"building_order_period"`
		ProjectilesActionPeriod int `yaml:"projectiles_action_period"`
		CleanupPeriod           int `yaml:"cleanup_period"`
	} `yaml:"engine"`

	AiSettings struct {
		AiActPeriod     int `yaml:"ai_acts_each"`
		AiAnalyzePeriod int `yaml:"ai_analyzes_each"`
	} `yaml:"aiSettings"`
}

func (c *yamlConfig) setDefaultValues() {
	c.DebugOutput = true
	c.TargetFPS = 60
	c.TargetTPS = 60

	c.Engine.BuildingActionPeriod = 5
	c.Engine.BuildingOrderPeriod = 11
	c.Engine.UnitActionPeriod = 2
	c.Engine.UnitListenOrderPeriod = 11
	c.Engine.CleanupPeriod = 6

	c.AiSettings.AiActPeriod = 60
	c.AiSettings.AiAnalyzePeriod = 70
}

func (c *yamlConfig) initFromFileOrCreate() {
	const filePath = "config.yaml"

	fiBytes, err := os.ReadFile(filePath)
	if errors.Is(err, os.ErrNotExist) {
		// set default values for current config
		c.setDefaultValues()

		res, _ := yaml.Marshal(c)
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

	err = yaml.Unmarshal(fiBytes, c)
}
