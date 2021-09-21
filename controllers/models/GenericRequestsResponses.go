package models

type LoginRequest struct {
	Username string `json:"username"`
}

type RefreshRequest struct {
	Username string `json:"username"`
	PlanetId string `json:"planet_id"`
}

type RefreshMineRequest struct {
	Username string `json:"username"`
	PlanetId string `json:"planet_id"`
	MineId   string `json:"mine_id"`
}

type UpgradeBuildingRequest struct {
	Username   string `json:"username"`
	PlanetId   string `json:"planet_id"`
	BuildingId string `json:"building_id"`
}

type UpdateResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
