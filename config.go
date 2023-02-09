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

	DamageOnArmorFactorsTable map[string]map[string]float64 `yaml:"damage_on_armor_factors_table"`

	CostForLevelUpMultiplier float64 `yaml:"cost_for_level_up_multiplier"`
	CostForLevelUpExponent   float64 `yaml:"cost_for_level_up_exponent"`

	AiSettings struct {
		AiActPeriod     int `yaml:"ai_act_period"`
		AiAnalyzePeriod int `yaml:"ai_analyze_period"`
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

	c.DamageOnArmorFactorsTable = map[string]map[string]float64{
		string(DAMAGETYPE_ANTI_INFANTRY): {string(ARMORTYPE_INFANTRY): 1, string(ARMORTYPE_BUILDING): 0.25, string(ARMORTYPE_HEAVY): 0.25},
		string(DAMAGETYPE_HEAVY):         {string(ARMORTYPE_INFANTRY): 0.25, string(ARMORTYPE_BUILDING): 1, string(ARMORTYPE_HEAVY): 1},
		string(DAMAGETYPE_ANTI_BUILDING): {string(ARMORTYPE_INFANTRY): 0.25, string(ARMORTYPE_BUILDING): 1, string(ARMORTYPE_HEAVY): 0.25},
		string(DAMAGETYPE_OMNI):          {string(ARMORTYPE_INFANTRY): 1, string(ARMORTYPE_BUILDING): 1, string(ARMORTYPE_HEAVY): 1},
	}
	c.CostForLevelUpMultiplier = 1.5
	c.CostForLevelUpExponent = 1.35

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
