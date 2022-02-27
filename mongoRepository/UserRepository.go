package mongoRepository

import (
	"context"
	"github.com/themane/MMOServer/constants"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
)

type UserRepositoryImpl struct {
	mongoURL string
	mongoDB  string
	logger   *constants.LoggingUtils
}

func NewUserRepository(mongoURL string, mongoDB string, logLevel string) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		mongoURL: mongoURL,
		mongoDB:  mongoDB,
		logger:   constants.NewLoggingUtils("USER_REPOSITORY", logLevel),
	}
}
func (u *UserRepositoryImpl) getMongoClient() (*mongo.Client, context.Context) {
	return getConnection(u.mongoURL)
}

func (u *UserRepositoryImpl) getCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(u.mongoDB).Collection("user_data")
}

func (u *UserRepositoryImpl) FindById(id string) (*repoModels.UserData, error) {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	result := repoModels.UserData{}
	filter := bson.M{"_id": id}
	singleResult := u.getCollection(client).FindOne(ctx, filter)
	err := singleResult.Decode(&result)
	if err != nil {
		u.logger.Error("Error in decoding user data received from Mongo", err)
		return nil, err
	}
	return &result, nil
}

func (u *UserRepositoryImpl) FindByUsername(username string) (*repoModels.UserData, error) {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	result := repoModels.UserData{}
	filter := bson.M{"profile.username": username}
	singleResult := u.getCollection(client).FindOne(ctx, filter)
	err := singleResult.Decode(&result)
	if err != nil {
		u.logger.Error("Error in decoding user data received from Mongo", err)
		return nil, err
	}
	return &result, nil
}

func (u *UserRepositoryImpl) AddExperience(id string, experience int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"profile.experience": experience}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Added experience id: %s, experience: %d\n", id, experience)
	return nil
}

func (u *UserRepositoryImpl) UpdateClanId(id string, clanId string) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"profile.clan_id": clanId}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Updated id: %s, clanId: %s\n", id, clanId)
	return nil
}

func (u *UserRepositoryImpl) UpgradeBuildingLevel(id string, planetId string, buildingId string,
	waterRequired int, grapheneRequired int, shelioRequired int, minutesRequired int) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"occupied_planets." + planetId + ".buildings." + buildingId + ".building_minutes_per_worker": minutesRequired,
		},
		"$inc": bson.M{
			"occupied_planets." + planetId + ".buildings." + buildingId + ".building_level": 1,
			"occupied_planets." + planetId + ".water.amount":                                -waterRequired,
			"occupied_planets." + planetId + ".graphene.amount":                             -grapheneRequired,
			"occupied_planets." + planetId + ".shelio":                                      -shelioRequired,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Upgraded id: %s, planetId: %s, buildingId: %s\n", id, planetId, buildingId)
	return nil
}

func (u *UserRepositoryImpl) AddResources(id string, planetId string, water int, graphene int, shelio int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".water.amount":    water,
		"occupied_planets." + planetId + ".graphene.amount": graphene,
		"occupied_planets." + planetId + ".shelio":          shelio,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Added id: %s, planetId: %s, water: %d, graphene: %d, shelio: %d\n", id, planetId, water, graphene, shelio)
	return nil
}

func (u *UserRepositoryImpl) UpdateMineResources(id string, planetId string, mineId string, water int, graphene int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".water.amount":              water,
		"occupied_planets." + planetId + ".graphene.amount":           graphene,
		"occupied_planets." + planetId + ".mines." + mineId + "mined": math.Max(float64(water), float64(graphene)),
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Updated mine resources id: %s, planetId: %s, water: %d, graphene: %d, mineId: %s\n",
		id, planetId, water, graphene, mineId)
	return nil
}

func (u *UserRepositoryImpl) UpdateWorkers(id string, planetId string, buildingId string, workers int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".buildings." + buildingId + "workers": workers,
		"occupied_planets." + planetId + ".population.workers.idle":             -workers,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Employed workers updated id: %s, planetId: %s, buildingId: %s, workers: %d\n", id, planetId, buildingId, workers)
	return nil
}

func (u *UserRepositoryImpl) AddPopulation(id string, planetId string, population int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".population.unemployed": population,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Added population id: %s, planetId: %s, population: %d\n", id, planetId, population)
	return nil
}

func (u *UserRepositoryImpl) UpdatePopulationRate(id string, planetId string, generationRate int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"occupied_planets." + planetId + ".population.generation_rate": generationRate,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Updated population generation rate id: %s, planetId: %s, rate: %d\n", id, planetId, generationRate)
	return nil
}

func (u *UserRepositoryImpl) RecruitWorkers(id string, planetId string, workers int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".population.unemployed":   -workers,
		"occupied_planets." + planetId + ".population.workers.idle": workers,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Assigned workers id: %s, planetId: %s, workers: %d\n", id, planetId, workers)
	return nil
}

func (u *UserRepositoryImpl) RecruitSoldiers(id string, planetId string, soldiers int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".population.unemployed":    -soldiers,
		"occupied_planets." + planetId + ".population.soldiers.idle": soldiers,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Assigned soldiers id: %s, planetId: %s, soldiers: %d\n", id, planetId, soldiers)
	return nil
}

func (u *UserRepositoryImpl) ScheduledPopulationIncrease(id string, planetIdGenerationRateMap map[string]int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	var incPopulationUpdate bson.D
	for planetId, generationRate := range planetIdGenerationRateMap {
		incPopulationUpdate = append(incPopulationUpdate,
			bson.E{Key: "occupied_planets." + planetId + ".population.unemployed", Value: generationRate},
		)
	}
	update := bson.D{{"$inc", incPopulationUpdate}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	return nil
}

func (u *UserRepositoryImpl) ScheduledWaterIncrease(id string, planetIdMiningRateMap map[string]map[string]int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	var miningUpdates bson.D
	for planetId, miningRates := range planetIdMiningRateMap {
		for mineId, miningRate := range miningRates {
			miningUpdates = append(miningUpdates,
				bson.E{Key: "occupied_planets." + planetId + ".mines." + mineId + ".mined", Value: miningRate},
				bson.E{Key: "occupied_planets." + planetId + ".water.amount", Value: miningRate},
			)
		}
	}
	update := bson.D{{"$inc", miningUpdates}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	return nil
}

func (u *UserRepositoryImpl) ScheduledGrapheneIncrease(id string, planetIdMiningRateMap map[string]map[string]int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	var miningUpdates bson.D
	for planetId, miningRates := range planetIdMiningRateMap {
		for mineId, miningRate := range miningRates {
			miningUpdates = append(miningUpdates,
				bson.E{Key: "occupied_planets." + planetId + ".mines." + mineId + ".mined", Value: miningRate},
				bson.E{Key: "occupied_planets." + planetId + ".graphene.amount", Value: miningRate},
			)
		}
	}
	update := bson.D{{"$inc", miningUpdates}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	return nil
}

func (u *UserRepositoryImpl) ScheduledPopulationConsumption(id string, planetIdTotalPopulationMap map[string]int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	var miningUpdates bson.D
	for planetId, totalPopulation := range planetIdTotalPopulationMap {
		miningUpdates = append(miningUpdates,
			bson.E{Key: "occupied_planets." + planetId + ".water.amount", Value: -totalPopulation},
		)
	}
	update := bson.D{{"$inc", miningUpdates}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	return nil
}

func (u *UserRepositoryImpl) ScheduledBuildingStateUpdate(id string, planetIdBuildingStateMap map[string]map[string]int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	var miningUpdates bson.D
	for planetId, buildingStates := range planetIdBuildingStateMap {
		for buildingId, buildingMinutes := range buildingStates {
			miningUpdates = append(miningUpdates,
				bson.E{Key: "occupied_planets." + planetId + ".buildings." + buildingId + ".building_minutes_per_worker", Value: -buildingMinutes},
			)
		}
	}
	update := bson.D{{"$inc", miningUpdates}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	return nil
}
