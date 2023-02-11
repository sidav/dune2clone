package main

type yamlConfig struct {
	DebugOutput bool `yaml:"debug_output"`
	LogToFile   bool `yaml:"log_to_file"`
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
		TicksPerNominalSecond     int `yaml:"ticks_per_nominal_second"`
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
		HarvestingSpeed             int `yaml:"harvesting_speed"`
		HarvesterUnloadSpeed        int `yaml:"harvester_unload_speed"`
	} `yaml:"economy"`

	Gameplay struct {
		DamageOnArmorFactorsTable                     map[string]map[string]float64 `yaml:"damage_on_armor_factors_table"`
		CostForLevelUpMultiplier                      float64                       `yaml:"cost_for_level_up_multiplier"`
		CostForLevelUpExponent                        float64                       `yaml:"cost_for_level_up_exponent"`
		VeterancyDamageBonusForLevelPercent           int                           `yaml:"veterancy_damage_bonus_for_level_percent"`
		VeterancyHpBonusForLevelPercent               int                           `yaml:"veterancy_hp_bonus_for_level_percent"`
		VeterancyFireCooldownReductionForLevelPercent int                           `yaml:"veterancy_fire_cooldown_reduction_for_level_percent"`
		VeterancySpeedBonusForLevel                   float64                       `yaml:"veterancy_speed_bonus_for_level"`
	} `yaml:"gameplay"`

	AiSettings struct {
		AiActPeriod            int     `yaml:"ai_act_period"`
		AiAnalyzePeriod        int     `yaml:"ai_analyze_period"`
		AiBuildSpeedMultiplier float64 `yaml:"ai_build_speed_multiplier"`
		AiExperienceMultiplier float64 `yaml:"ai_experience_multiplier"`
		AiStoragesMultiplier   float64 `yaml:"ai_storages_multiplier"`
		AiVisionCheat          bool    `yaml:"ai_vision_cheat"`
	} `yaml:"aiSettings"`
}

func (c *yamlConfig) setDefaultValues() {
	c.DebugOutput = true
	c.LogToFile = false
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
	c.Engine.TicksPerNominalSecond = 60

	c.Economy.ResourcesGrowthPeriod = 300
	c.Economy.ResourcesGrowthRadius = 4
	c.Economy.ResourcesGrowthMin = 25
	c.Economy.ResourcesGrowthMax = 75
	c.Economy.ResourcesInTileMinGenerated = 50
	c.Economy.ResourceInTilePoorMax = 100
	c.Economy.ResourceInTileMediumMax = 225
	c.Economy.ResourceInTileRichMax = 350
	c.Economy.HarvestingSpeed = 2
	c.Economy.HarvesterUnloadSpeed = 3

	c.Gameplay.DamageOnArmorFactorsTable = map[string]map[string]float64{
		string(DAMAGETYPE_ANTI_INFANTRY): {string(ARMORTYPE_INFANTRY): 1, string(ARMORTYPE_BUILDING): 0.1, string(ARMORTYPE_HEAVY): 0.1},
		string(DAMAGETYPE_HEAVY):         {string(ARMORTYPE_INFANTRY): 0.25, string(ARMORTYPE_BUILDING): 1, string(ARMORTYPE_HEAVY): 1},
		string(DAMAGETYPE_ANTI_BUILDING): {string(ARMORTYPE_INFANTRY): 0.25, string(ARMORTYPE_BUILDING): 1, string(ARMORTYPE_HEAVY): 0.25},
		string(DAMAGETYPE_OMNI):          {string(ARMORTYPE_INFANTRY): 1, string(ARMORTYPE_BUILDING): 1, string(ARMORTYPE_HEAVY): 1},
	}
	c.Gameplay.CostForLevelUpMultiplier = 1.5
	c.Gameplay.CostForLevelUpExponent = 1.35
	c.Gameplay.VeterancyDamageBonusForLevelPercent = 10
	c.Gameplay.VeterancyHpBonusForLevelPercent = 5
	c.Gameplay.VeterancyFireCooldownReductionForLevelPercent = 4
	c.Gameplay.VeterancySpeedBonusForLevel = 0.05

	c.AiSettings.AiActPeriod = 60
	c.AiSettings.AiAnalyzePeriod = 70
	c.AiSettings.AiBuildSpeedMultiplier = 1.11
	c.AiSettings.AiStoragesMultiplier = 2
	c.AiSettings.AiExperienceMultiplier = 2.5
	c.AiSettings.AiVisionCheat = false
}
