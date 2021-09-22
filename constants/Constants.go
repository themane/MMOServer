package constants

const (
	WATER    string = "WATER"
	GRAPHENE string = "GRAPHENE"
	SHELIO   string = "SHELIO"

	LEADER     string = "LEADER"
	SUB_LEADER string = "SUB_LEADER"
	MEMBER     string = "MEMBER"

	PRODUCING_STATE = "PRODUCING"
	UPGRADING_STATE = "UPGRADING"
)

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
	MinutesRequired     int `json:"minutes_required"`
}
