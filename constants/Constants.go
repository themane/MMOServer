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

type DefenceConstants struct {
	MaxLevel int                             `json:"max_level"`
	Type     string                          `json:"type"`
	Levels   map[string]DefenceLevelConstant `json:"levels"`
}

type DefenceLevelConstant struct {
	RequiredSoldiers int `json:"required_soldiers"`
	RequiredWorkers  int `json:"required_workers"`
	HitPoints        int `json:"hit_points"`
	Armor            int `json:"armor"`
	MinAttack        int `json:"min_attack"`
	MaxAttack        int `json:"max_attack"`
	Range            int `json:"range"`
	SingleHitTargets int `json:"single_hit_targets"`
}

type ShipConstants struct {
	MaxLevel int                          `json:"max_level"`
	Type     string                       `json:"type"`
	Levels   map[string]ShipLevelConstant `json:"levels"`
}

type ShipLevelConstant struct {
	RequiredSoldiers int `json:"required_soldiers"`
	RequiredWorkers  int `json:"required_workers"`
	HitPoints        int `json:"hit_points"`
	Armor            int `json:"armor"`
	ResourceCapacity int `json:"resource_capacity"`
	WorkerCapacity   int `json:"worker_capacity"`
	MinAttack        int `json:"min_attack"`
	MaxAttack        int `json:"max_attack"`
	Range            int `json:"range"`
	Speed            int `json:"speed"`
}

type UpgradeConstants struct {
	MaxLevel int                             `json:"max_level"`
	Levels   map[string]UpgradeLevelConstant `json:"levels"`
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
