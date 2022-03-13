package models

import (
	"github.com/themane/MMOServer/constants"
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

func InitPopulation(planetUser repoModels.PlanetUser, militaryConstants map[string]constants.MilitaryConstants) *Population {
	p := new(Population)
	totalEmployedWorkers, totalEmployedSoldiers := repoModels.GetEmployedPopulation(planetUser, militaryConstants)
	p.Workers = models.EmployedPopulation{
		Idle:  planetUser.Population.IdleWorkers,
		Total: totalEmployedWorkers + planetUser.Population.IdleWorkers,
	}
	p.Soldiers = models.EmployedPopulation{
		Idle:  planetUser.Population.IdleSoldiers,
		Total: totalEmployedSoldiers + planetUser.Population.IdleSoldiers,
	}
	p.Total = planetUser.Population.Unemployed + p.Workers.Total + p.Soldiers.Total
	p.GenerationRate = planetUser.Population.GenerationRate
	p.Unemployed = planetUser.Population.Unemployed
	return p
}
