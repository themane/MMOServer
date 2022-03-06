package buildings

type IntegerBuildingAttributes struct {
	Current int `json:"current" example:"2"`
	Next    int `json:"next" example:"4"`
	Max     int `json:"max,omitempty" example:"16"`
}

type FloatBuildingAttributes struct {
	Current float64 `json:"current" example:"0.2"`
	Next    float64 `json:"next" example:"0.4"`
	Max     float64 `json:"max,omitempty" example:"2.5"`
}
