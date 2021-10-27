package models

type Formation struct {
	ShipName string `json:"ship_name" bson:"ship_name" example:"ANUJ"`
	Quantity int    `json:"quantity" bson:"quantity" example:"15"`
}

type Notification struct {
	Info    string `json:"info"`
	Warning string `json:"warn"`
	Error   string `json:"error"`
}
