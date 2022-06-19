package models

type ClanData struct {
	Id         string       `json:"_id" bson:"_id"`
	Name       string       `json:"name" bson:"name"`
	Experience int          `json:"experience" bson:"experience"`
	Members    []ClanMember `json:"members" bson:"members"`
}

type ClanMember struct {
	Id   string `json:"_id" bson:"_id"`
	Role string `json:"role" bson:"role"`
}

type ClanRepository interface {
	FindById(id string) (*ClanData, error)
}
