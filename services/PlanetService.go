package services

import (
	"errors"
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/controllers/models/researches"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"strconv"
)

type PlanetService struct {
	userRepository    repoModels.UserRepository
	upgradeConstants  map[string]constants.UpgradeConstants
	buildingConstants map[string]map[string]map[string]interface{}
	researchConstants map[string]constants.ResearchConstants
	logger            *constants.LoggingUtils
}

func NewPlanetService(
	userRepository repoModels.UserRepository,
	buildingConstants map[string]map[string]map[string]interface{},
	researchConstants map[string]constants.ResearchConstants,
	logLevel string,
) *PlanetService {
	return &PlanetService{
		userRepository:    userRepository,
		buildingConstants: buildingConstants,
		researchConstants: researchConstants,
		logger:            constants.NewLoggingUtils("PLANET_SERVICE", logLevel),
	}
}

func (p *PlanetService) UpdatePopulationRate(username string, planetId string, generationRate int) error {
	userData, err := p.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	planetUser := userData.GetOccupiedPlanet(planetId)
	if planetUser == nil {
		return errors.New("planet not occupied")
	}
	currentDeployedWorkers := planetUser.GetBuilding(constants.PopulationControlCenter).Workers
	populationControlCenterLevel := planetUser.GetBuilding(constants.PopulationControlCenter).BuildingLevel

	populationControlCenterConstants := p.buildingConstants[constants.PopulationControlCenter][strconv.Itoa(populationControlCenterLevel)]
	maxPopulationGenerationRate := int(models.MaxSelectablePopulationGenerationRate(populationControlCenterConstants, currentDeployedWorkers))
	if maxPopulationGenerationRate < generationRate {
		return errors.New("rate above maximum")
	}

	err = p.userRepository.UpdatePopulationRate(userData.Id, planetId, generationRate)
	if err != nil {
		return err
	}
	return nil
}

func (p *PlanetService) EmployPopulation(username string, planetId string, workers int, soldiers int) error {
	userData, err := p.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	planetUser := userData.GetOccupiedPlanet(planetId)
	if planetUser == nil {
		return errors.New("planet not occupied")
	}
	unemployedPopulation := planetUser.Population.Unemployed
	if workers+soldiers > unemployedPopulation {
		return errors.New("not enough population to employ")
	}

	err = p.userRepository.Recruit(userData.Id, planetId, workers, soldiers)
	if err != nil {
		return err
	}
	return nil
}

func (p *PlanetService) KillPopulation(username string, planetId string, unemployed int, workers int, soldiers int) error {
	userData, err := p.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	planetUser := userData.GetOccupiedPlanet(planetId)
	if planetUser == nil {
		return errors.New("planet not occupied")
	}
	if unemployed > planetUser.Population.Unemployed ||
		workers > planetUser.Population.IdleWorkers ||
		soldiers > planetUser.Population.IdleSoldiers {
		return errors.New("not enough population to kill")
	}

	err = p.userRepository.KillPopulation(userData.Id, planetId, unemployed, workers, soldiers)
	if err != nil {
		return err
	}
	return nil
}

func (p *PlanetService) ReserveResources(username string, planetId string, water int, graphene int) error {
	userData, err := p.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	planetUser := userData.GetOccupiedPlanet(planetId)
	if planetUser == nil {
		return errors.New("planet not occupied")
	}
	if water > planetUser.Water.Amount ||
		graphene > planetUser.Graphene.Amount {
		return errors.New("not enough resources to reserve")
	}

	err = p.userRepository.ReserveResources(userData.Id, planetId, water, graphene)
	if err != nil {
		return err
	}
	return nil
}

func (p *PlanetService) CancelReserveResources(username string, planetId string) error {
	userData, err := p.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	planetUser := userData.GetOccupiedPlanet(planetId)
	if planetUser == nil {
		return errors.New("planet not occupied")
	}
	if planetUser.Water.Reserving == 0 &&
		planetUser.Graphene.Reserving == 0 {
		return errors.New("nothing to cancel")
	}

	err = p.userRepository.ReserveResources(userData.Id, planetId, -planetUser.Water.Reserving, -planetUser.Graphene.Reserving)
	if err != nil {
		return err
	}
	return nil
}

func (p *PlanetService) ExtractReservedResources(username string, planetId string, water int, graphene int) error {
	userData, err := p.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	planetUser := userData.GetOccupiedPlanet(planetId)
	if planetUser == nil {
		return errors.New("planet not occupied")
	}
	if water > planetUser.Water.Reserved ||
		graphene > planetUser.Graphene.Reserved {
		return errors.New("not enough resources to extract")
	}

	err = p.userRepository.ExtractReservedResources(userData.Id, planetId, water, graphene)
	if err != nil {
		return err
	}
	return nil
}

func (p *PlanetService) Research(username string, planetId string, researchName string) error {
	userData, err := p.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	planetUser := userData.GetOccupiedPlanet(planetId)
	if planetUser == nil {
		return errors.New("planet not occupied")
	}
	requirements, err1 := p.validateAndGetRequirements(*planetUser, researchName)
	if err1 != nil {
		return err1
	}
	if planetUser.GetResearch(researchName).Level == 0 {
		err = p.userRepository.Research(userData.Id, planetId, researchName,
			requirements.GrapheneRequired, requirements.WaterRequired, requirements.ShelioRequired, requirements.MinutesRequiredPerWorker)
		if err != nil {
			return err
		}
	} else {
		err = p.userRepository.ResearchUpgrade(userData.Id, planetId, researchName,
			requirements.GrapheneRequired, requirements.WaterRequired, requirements.ShelioRequired, requirements.MinutesRequiredPerWorker)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PlanetService) validateAndGetRequirements(planetUser repoModels.PlanetUser, researchName string) (*researches.NextLevelRequirements, error) {
	if planetUser.GetResearch(researchName).ResearchMinutesPerWorker != 0 {
		return nil, errors.New("already under progress")
	}
	nextLevelRequirements := researches.NextLevelRequirements{}
	nextLevelRequirements.Init(planetUser.GetResearch(researchName).Level, p.researchConstants[researchName])
	if int(nextLevelRequirements.WaterRequired) > planetUser.Water.Amount ||
		int(nextLevelRequirements.GrapheneRequired) > planetUser.Graphene.Amount ||
		int(nextLevelRequirements.ShelioRequired) > planetUser.Shelio {
		return nil, errors.New("not enough resources to research")
	}
	for _, r := range nextLevelRequirements.SpecialRequirements {
		if !r.Fulfilled {
			return nil, errors.New("requirements not fulfilled")
		}
	}
	return &nextLevelRequirements, nil
}

func (p *PlanetService) CancelResearch(username string, planetId string, researchName string) error {
	userData, err := p.userRepository.FindByUsername(username)
	if err != nil {
		return err
	}
	planetUser := userData.GetOccupiedPlanet(planetId)
	if planetUser == nil {
		return errors.New("planet not occupied")
	}
	researchUser := planetUser.GetResearch(researchName)
	if researchUser.ResearchMinutesPerWorker < 0 {
		return errors.New("nothing to cancel")
	}
	cancelReturns := researches.CancelReturns{}
	cancelReturns.Init(researchUser.ResearchMinutesPerWorker, researchUser.Level, p.researchConstants[researchName])
	err = p.userRepository.CancelResearch(userData.Id, planetId, researchName,
		cancelReturns.GrapheneReturned, cancelReturns.WaterReturned, cancelReturns.ShelioReturned)
	if err != nil {
		return err
	}
	return nil
}
