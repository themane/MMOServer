package models

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type EmployedPopulation struct {
	Total int `json:"total" bson:"total" example:"21"`
	Idle  int `json:"idle" bson:"idle" example:"4"`
}

type PlanetPosition struct {
	Id     string `json:"_id" example:"023:049:07"`
	System int    `json:"system" bson:"system" example:"23"`
	Sector int    `json:"sector" bson:"sector" example:"49"`
	Planet int    `json:"planet" bson:"planet" example:"7"`
}

func (p *PlanetPosition) Init(system int, sector int, planet int) {
	p.System = system
	p.Sector = sector
	p.Planet = planet
	p.Id = PlanetId(system, sector, planet)
}

func InitPlanetPositionByPosition(system int, sector int, planet int) PlanetPosition {
	position := PlanetPosition{}
	position.Init(system, sector, planet)
	return position
}

func InitPlanetPositionById(id string) (*PlanetPosition, error) {
	split := strings.Split(id, ":")
	if len(split) != 3 {
		return nil, errors.New("planet-id not correct: " + id)
	}
	system, err := strconv.Atoi(split[0])
	if err != nil {
		log.Print(err)
		return nil, errors.New("planet-id not correct: " + id)
	}
	sector, err := strconv.Atoi(split[1])
	if err != nil {
		log.Print(err)
		return nil, errors.New("planet-id not correct: " + id)
	}
	planet, err := strconv.Atoi(split[2])
	if err != nil {
		log.Print(err)
		return nil, errors.New("planet-id not correct: " + id)
	}
	position := InitPlanetPositionByPosition(system, sector, planet)
	return &position, nil
}

func (p PlanetPosition) SystemId() string {
	return fmt.Sprintf("%03d", p.System)
}

func (p PlanetPosition) SectorId() string {
	return fmt.Sprintf("%03d:%03d", p.System, p.Sector)
}

func (p PlanetPosition) PlanetId() string {
	return fmt.Sprintf("%03d:%03d:%02d", p.System, p.Sector, p.Planet)
}

func PlanetId(system int, sector int, planet int) string {
	return fmt.Sprintf("%03d:%03d:%02d", system, sector, planet)
}

func (p PlanetPosition) Clone() PlanetPosition {
	position := PlanetPosition{}
	position.Init(p.System, p.Sector, p.Planet)
	return position
}
