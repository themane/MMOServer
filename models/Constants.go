package models

type ExperienceConstants struct {
	User ExperienceLevelConstants `json:"user"`
	Clan ExperienceLevelConstants `json:"clan"`
}

type ExperienceLevelConstants struct {
	MaxLevel            int                           `json:"max_level"`
	ExperiencesRequired map[string]ExperienceRequired `json:"experiences_required"`
}

type ExperienceRequired struct {
	ExperienceRequired int `json:"experience_required"`
}

type ResourceConstants struct {
	MaxLevel int                      `json:"max_level"`
	Levels   map[string]LevelConstant `json:"levels"`
}

type LevelConstant struct {
	WaterRequired       int `json:"water_required"`
	GrapheneRequired    int `json:"graphene_required"`
	ShelioRequired      int `json:"shelio_required"`
	MiningRatePerWorker int `json:"mining_rate_per_worker"`
	WorkersMaxLimit     int `json:"workers_max_limit"`
}
