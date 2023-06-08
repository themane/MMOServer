package models

type Formation struct {
	ShipName string `json:"ship_name" bson:"ship_name" example:"ANUJ"`
	Quantity int    `json:"quantity" bson:"quantity" example:"15"`
}

type Notification struct {
	Tutorial string `json:"tutorial"`
	Warning  string `json:"warn"`
	Error    string `json:"error"`
}

type UserSocialDetails struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	PictureUrl    string `json:"picture_url"`
	Authenticator string `json:"authenticator"`
}

type FbUserDetails struct {
	Id      string        `json:"id"`
	Name    string        `json:"name"`
	Email   string        `json:"email"`
	Picture FbUserPicture `json:"picture"`
}

type FbUserPicture struct {
	Data FbUserPictureData `json:"data"`
}

type FbUserPictureData struct {
	Url string `json:"url"`
}
