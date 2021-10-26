package models

type LoginRequest struct {
	Username string `json:"username" example:"devashish"`
}

type RefreshRequest struct {
	Username string `json:"username" example:"devashish"`
	PlanetId string `json:"planet_id" example:"001:002:03"`
}

type RefreshMineRequest struct {
	Username string `json:"username" example:"devashish"`
	PlanetId string `json:"planet_id" example:"001:002:03"`
	MineId   string `json:"mine_id" example:"G018"`
}

type UpgradeBuildingRequest struct {
	Username   string `json:"username" example:"devashish"`
	PlanetId   string `json:"planet_id" example:"001:002:03"`
	BuildingId string `json:"building_id" example:"GMP0018"`
}
