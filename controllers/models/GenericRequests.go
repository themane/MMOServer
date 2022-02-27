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

type UpdateBuildingWorkersRequest struct {
	Username   string `json:"username" example:"devashish"`
	PlanetId   string `json:"planet_id" example:"001:002:03"`
	BuildingId string `json:"building_id" example:"GMP0018"`
	Workers    int    `json:"workers" example:"10"`
}

type VisitSectorRequest struct {
	Username string `json:"username" example:"devashish"`
	Sector   string `json:"sector" example:"005:001"`
}

type TeleportRequest struct {
	Username string `json:"username" example:"devashish"`
	Planet   string `json:"planet" example:"005:001:03"`
}
