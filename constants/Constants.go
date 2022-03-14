package constants

type ExperienceConstants struct {
	MaxLevel            int                                `json:"max_level"`
	ExperiencesRequired map[string]ExperienceLevelConstant `json:"experiences_required"`
}

type ExperienceLevelConstant struct {
	ExperienceRequired int `json:"experience_required"`
}

type MiningConstants struct {
	MaxLevel int                            `json:"max_level"`
	Levels   map[string]MiningLevelConstant `json:"levels"`
}

type MiningLevelConstant struct {
	MiningRatePerWorker int `json:"mining_rate_per_worker"`
	WorkersMaxLimit     int `json:"workers_max_limit"`
}

type MilitaryConstants struct {
	MaxLevel int                               `json:"max_level"`
	Type     string                            `json:"type"`
	Levels   map[string]map[string]interface{} `json:"levels"`
}

type UpgradeConstants struct {
	MaxLevel int                             `json:"max_level"`
	Levels   map[string]UpgradeLevelConstant `json:"levels"`
}

type ResearchConstants struct {
	MaxLevel     int                               `json:"max_level"`
	Description  string                            `json:"description"`
	Bonus        map[string]map[string]interface{} `json:"bonus"`
	Requirements map[string]map[string]interface{} `json:"requirements"`
}

type UpgradeLevelConstant struct {
	WaterRequired    int `json:"water_required"`
	GrapheneRequired int `json:"graphene_required"`
	ShelioRequired   int `json:"shelio_required"`
	MinutesRequired  int `json:"minutes_required"`
}

type SpeciesConstants struct {
	AvailableUnits []string `json:"available_units"`
}
