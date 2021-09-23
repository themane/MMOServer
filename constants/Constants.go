package constants

// Resources
const (
	Water    string = "WATER"
	Graphene string = "GRAPHENE"
	Shelio   string = "SHELIO"
)

// Experience Constant Types
const (
	UserExperiences string = "USER"
	ClanExperiences string = "CLAN"
)

// Paltan Roles
const (
	PaltanLeader    string = "LEADER"
	PaltanSubLeader string = "SUB_LEADER"
	PaltanMember    string = "MEMBER"
)

//  Building States
const (
	WorkingState   = "WORKING"
	UpgradingState = "UPGRADING"
)

//  Building Types
const (
	WaterMiningPlant    = "WATER_MINING_PLANT"
	GrapheneMiningPlant = "GRAPHENE_MINING_PLANT"
	Shield              = "SHIELD"
)

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
	Levels   map[string]DefenceLevelConstant `json:"levels"`
}

type DefenceLevelConstant struct {
	HitPoints        int `json:"hit_points" bson:"hit_points"`
	Attack           int `json:"attack" bson:"attack"`
	Range            int `json:"range" bson:"range"`
	SingleHitTargets int `json:"single_hit_targets" bson:"single_hit_targets"`
}

type BuildingConstants struct {
	MaxLevel int                             `json:"max_level"`
	Levels   map[string]DefenceLevelConstant `json:"levels"`
}

type BuildingLevelConstant struct {
	WaterRequired    int `json:"water_required"`
	GrapheneRequired int `json:"graphene_required"`
	ShelioRequired   int `json:"shelio_required"`
	MinutesRequired  int `json:"minutes_required"`
}
