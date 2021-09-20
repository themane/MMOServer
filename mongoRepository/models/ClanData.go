package models

type ClanData struct {
	Id   string `json:"_id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	//Experience int          `json:"experience" bson:"experience"`
	Members []ClanMember `json:"members" bson:"members"`
}

type ClanMember struct {
	MemberId string `json:"member_id" bson:"member_id"`
	Role     string `json:"role" bson:"role"`
}

type ClanRepository interface {
	FindById(id string) (*ClanData, error)
}
