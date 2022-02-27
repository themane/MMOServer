package models

import (
	"errors"
	"fmt"
	"github.com/themane/MMOServer/models"
	"log"
	"strconv"
	"strings"
)

type SectorPosition struct {
	Id     string `json:"_id" example:"023:049"`
	System int    `json:"system" example:"23"`
	Sector int    `json:"sector" example:"49"`
}

func (s *SectorPosition) Init(system int, sector int) {
	s.Id = SectorId(system, sector)
	s.System = system
	s.Sector = sector
}

func InitSectorPositionByPosition(system int, sector int) SectorPosition {
	position := SectorPosition{}
	position.Init(system, sector)
	return position
}

func InitSectorPositionByPlanetPosition(planetPosition models.PlanetPosition) SectorPosition {
	position := SectorPosition{}
	position.Init(planetPosition.System, planetPosition.Sector)
	return position
}

func InitSectorPositionById(id string) (*SectorPosition, error) {
	split := strings.Split(id, ":")
	if len(split) != 2 {
		return nil, errors.New("sector-id not correct: " + id)
	}
	system, err := strconv.Atoi(split[0])
	if err != nil {
		log.Print(err)
		return nil, errors.New("sector-id not correct: " + id)
	}
	sector, err := strconv.Atoi(split[1])
	if err != nil {
		log.Print(err)
		return nil, errors.New("sector-id not correct: " + id)
	}
	position := InitSectorPositionByPosition(system, sector)
	return &position, nil
}

func (s SectorPosition) SystemId() string {
	return fmt.Sprintf("%03d", s.System)
}

func (s SectorPosition) SectorId() string {
	return fmt.Sprintf("%03d:%03d", s.System, s.Sector)
}

func SectorId(system int, sector int) string {
	return fmt.Sprintf("%03d:%03d", system, sector)
}

func (s SectorPosition) Clone() SectorPosition {
	position := SectorPosition{}
	position.Init(s.System, s.Sector)
	return position
}
