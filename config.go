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
		UnitsActionPeriod         int `yaml:"units_action_period"`
		UnitsListenOrderPeriod    int `yaml:"units_listen_order_period"`
		BuildingsActionPeriod     int `yaml:"buildings_action_period"`
		BuildingListenOrderPeriod int `yaml:"building_listen_order_period"`
		ProjectilesActionPeriod   int `yaml:"projectiles_action_period"`
		CleanupPeriod             int `yaml:"cleanup_period"`
		RegenHpPeriod             int `yaml:"regen_hp_period"`
		BuildingAnimationTicks    int `yaml:"building_animation_ticks"`
	} `yaml:"engine"`

	Economy struct {
		ResourcesGrowthPeriod       int `yaml:"resources_growth_period,omitempty"`
		ResourcesGrowthRadius       int `yaml:"resources_growth_radius,omitempty"`
		ResourcesGrowthMin          int `yaml:"resources_growth_min,omitempty"`
		ResourcesGrowthMax          int `yaml:"resources_growth_max,omitempty"`
		ResourcesInTileMinGenerated int `yaml:"resources_in_tile_min_generated,omitempty"`
		ResourceInTilePoorMax       int `yaml:"resource_in_tile_poor_max,omitempty"`
		ResourceInTileMediumMax     int `yaml:"resource_in_tile_medium_max,omitempty"`
		ResourceInTileRichMax       int `yaml:"resource_in_tile_rich_max,omitempty"`
	} `yaml:"economy"`

	AiSettings struct {
		AiActPeriod     int `yaml:"ai_acts_each"`
		AiAnalyzePeriod int `yaml:"ai_analyzes_each"`
	} `yaml:"aiSettings"`
}

func (c *yamlConfig) setDefaultValues() {
	c.DebugOutput = true
	c.TargetFPS = 60
	c.TargetTPS = 60

	c.Engine.BuildingsActionPeriod = 5
	c.Engine.BuildingListenOrderPeriod = 11
	c.Engine.UnitsActionPeriod = 2
	c.Engine.UnitsListenOrderPeriod = 11
	c.Engine.ProjectilesActionPeriod = 2
	c.Engine.CleanupPeriod = 6
	c.Engine.RegenHpPeriod = 60
	c.Engine.BuildingAnimationTicks = 12

	c.Economy.ResourcesGrowthPeriod = 420
	c.Economy.ResourcesGrowthRadius = 5
	c.Economy.ResourcesGrowthMin = 10
	c.Economy.ResourcesGrowthMax = 100
	c.Economy.ResourcesInTileMinGenerated = 50
	c.Economy.ResourceInTilePoorMax = 100
	c.Economy.ResourceInTileMediumMax = 225
	c.Economy.ResourceInTileRichMax = 350

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
