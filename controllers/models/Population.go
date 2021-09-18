package models

import (
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
)

type Population struct {
	Total          int                       `json:"total" example:"45"`
	GenerationRate int                       `json:"generation_rate" example:"3"`
	Unemployed     int                       `json:"unemployed" example:"3"`
	Workers        models.EmployedPopulation `json:"workers"`
	Soldiers       models.EmployedPopulation `json:"soldiers"`
}

func (p *Population) Init(planetUser repoModels.PlanetUser) {
	p.Total = planetUser.Population.Unemployed + planetUser.Population.Workers.Total + planetUser.Population.Soldiers.Total
	p.GenerationRate = planetUser.Population.GenerationRate
	p.Unemployed = planetUser.Population.Unemployed
	p.Workers = planetUser.Population.Workers
	p.Soldiers = planetUser.Population.Soldiers
}
