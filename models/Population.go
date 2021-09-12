package models

type Population struct {
	Total          int                `json:"total" example:"45"`
	GenerationRate int                `json:"generation_rate" example:"3"`
	Unemployed     int                `json:"unemployed" example:"3"`
	Workers        EmployedPopulation `json:"workers"`
	Soldiers       EmployedPopulation `json:"soldiers"`
}

type EmployedPopulation struct {
	Total int `json:"total" example:"21"`
	Idle  int `json:"idle" example:"4"`
}

func (p *Population) Init(planetUser PlanetUser) {
	p.Total = planetUser.Population.Unemployed + planetUser.Population.Workers.Total + planetUser.Population.Soldiers.Total
	p.GenerationRate = planetUser.Population.GenerationRate
	p.Unemployed = planetUser.Population.Unemployed
	p.Workers = planetUser.Population.Workers
	p.Soldiers = planetUser.Population.Soldiers
}
