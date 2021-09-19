package models

import "go.mongodb.org/mongo-driver/x/mongo/driver/uuid"

type ClanData struct {
	Id         uuid.UUID    `json:"_id"`
	Name       string       `json:"name"`
	Experience int          `json:"experience"`
	Members    []ClanMember `json:"members"`
}

type ClanMember struct {
	Id   uuid.UUID `json:"_id"`
	Role string    `json:"role"`
}

type ClanRepository interface {
	FindById(id uuid.UUID) (*ClanData, error)
}
