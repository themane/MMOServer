package models

type ClanData struct {
	Id         string       `json:"_id"`
	Name       string       `json:"name"`
	Experience int          `json:"experience"`
	Members    []ClanMember `json:"members"`
}

type ClanMember struct {
	Id   string `json:"_id"`
	Role string `json:"role"`
}

type ClanRepository interface {
	FindById(id string) (ClanData, error)
}
