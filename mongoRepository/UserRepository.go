package mongoRepository

import (
	"context"
	"github.com/themane/MMOServer/constants"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	u.logger.Printf("Added experience id: %s, experience: %s\n", id, experience)
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
	result := u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	bsonResult, err := result.DecodeBytes()
	if err != nil {
		return err
	}
	u.logger.Printf(bsonResult.String())
	u.logger.Printf("Upgraded id: %s, planetId: %s, buildingId: %s\n", id, planetId, buildingId)
	return nil
}

func (u *UserRepositoryImpl) CancelUpgradeBuildingLevel(id string, planetId string, buildingId string,
	waterReturned int, grapheneReturned int, shelioReturned int) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"occupied_planets." + planetId + ".buildings." + buildingId + ".building_minutes_per_worker": 0,
		},
		"$inc": bson.M{
			"occupied_planets." + planetId + ".buildings." + buildingId + ".building_level": -1,
			"occupied_planets." + planetId + ".water.amount":                                waterReturned,
			"occupied_planets." + planetId + ".graphene.amount":                             grapheneReturned,
			"occupied_planets." + planetId + ".shelio":                                      shelioReturned,
		},
	}
	result := u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	bsonResult, err := result.DecodeBytes()
	if err != nil {
		return err
	}
	u.logger.Printf(bsonResult.String())
	u.logger.Printf("Upgraded id: %s, planetId: %s, buildingId: %s\n", id, planetId, buildingId)
	return nil
}

func (u *UserRepositoryImpl) UpdateWorkers(id string, planetId string, buildingId string, workers int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".buildings." + buildingId + ".workers": workers,
		"occupied_planets." + planetId + ".population.workers":                   -workers,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Employed workers updated id: %s, planetId: %s, buildingId: %s, workers: %s\n", id, planetId, buildingId, workers)
	return nil
}

func (u *UserRepositoryImpl) UpdateSoldiers(id string, planetId string, buildingId string, soldiers int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".buildings." + buildingId + ".soldiers": soldiers,
		"occupied_planets." + planetId + ".population.soldiers":                   -soldiers,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Employed soldiers updated id: %s, planetId: %s, buildingId: %s, soldiers: %s\n", id, planetId, buildingId, soldiers)
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
	u.logger.Printf("Updated population generation rate id: %s, planetId: %s, rate: %s\n", id, planetId, generationRate)
	return nil
}

func (u *UserRepositoryImpl) Recruit(id string, planetId string, workers int, soldiers int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".population.unemployed": -(workers + soldiers),
		"occupied_planets." + planetId + ".population.workers":    workers,
		"occupied_planets." + planetId + ".population.soldiers":   soldiers,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Assigned workers and soldiers id: %s, planetId: %s, workers: %s, soldiers: %s\n", id, planetId, workers, soldiers)
	return nil
}

func (u *UserRepositoryImpl) KillPopulation(id string, planetId string, unemployed int, workers int, soldiers int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".population.unemployed": -unemployed,
		"occupied_planets." + planetId + ".population.workers":    -workers,
		"occupied_planets." + planetId + ".population.soldiers":   -soldiers,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Killed population: %s, planetId: %s, unemployed: %s, workers: %s, soldiers: %s\n", id, planetId, unemployed, workers, soldiers)
	return nil
}

func (u *UserRepositoryImpl) ReserveResources(id string, planetId string, water int, graphene int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".water.reserving":    water,
		"occupied_planets." + planetId + ".water.amount":       -water,
		"occupied_planets." + planetId + ".graphene.reserving": graphene,
		"occupied_planets." + planetId + ".graphene.amount":    -graphene,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Marked for reserving resources: %s, planetId: %s, water: %s, graphene: %s\n", id, planetId, water, graphene)
	return nil
}

func (u *UserRepositoryImpl) ExtractReservedResources(id string, planetId string, water int, graphene int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".water.reserved":    -water,
		"occupied_planets." + planetId + ".water.amount":      water,
		"occupied_planets." + planetId + ".graphene.reserved": -graphene,
		"occupied_planets." + planetId + ".graphene.amount":   graphene,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Extracted reserved resources: %s, planetId: %s, water: %s, graphene: %s\n", id, planetId, water, graphene)
	return nil
}

func (u *UserRepositoryImpl) Research(id string, planetId string, researchName string,
	grapheneRequired float64, waterRequired float64, shelioRequired float64, minutesRequired float64) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".researches." + researchName + ".level":                       1,
		"occupied_planets." + planetId + ".researches." + researchName + ".research_minutes_per_worker": minutesRequired,
		"occupied_planets." + planetId + ".water.amount":                                                -waterRequired,
		"occupied_planets." + planetId + ".graphene.amount":                                             -grapheneRequired,
		"occupied_planets." + planetId + ".shelio":                                                      -shelioRequired,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Initiated research: %s, planetId: %s, researchName: %s\n", id, planetId, researchName)
	return nil
}

func (u *UserRepositoryImpl) ResearchUpgrade(id string, planetId string, researchName string,
	grapheneRequired float64, waterRequired float64, shelioRequired float64, minutesRequired float64) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".researches." + researchName + ".level":                       1,
		"occupied_planets." + planetId + ".researches." + researchName + ".research_minutes_per_worker": minutesRequired,
		"occupied_planets." + planetId + ".water.amount":                                                -waterRequired,
		"occupied_planets." + planetId + ".graphene.amount":                                             -grapheneRequired,
		"occupied_planets." + planetId + ".shelio":                                                      -shelioRequired,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Initiated research: %s, planetId: %s, researchName: %s\n", id, planetId, researchName)
	return nil
}

func (u *UserRepositoryImpl) CancelResearch(id string, planetId string, researchName string,
	grapheneReturned int, waterReturned int, shelioReturned int) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"occupied_planets." + planetId + ".researches." + researchName + ".research_minutes_per_worker": 0,
		},
		"$inc": bson.M{
			"occupied_planets." + planetId + ".researches." + researchName + ".level": -1,
			"occupied_planets." + planetId + ".water.amount":                          waterReturned,
			"occupied_planets." + planetId + ".graphene.amount":                       grapheneReturned,
			"occupied_planets." + planetId + ".shelio":                                shelioReturned,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Canceled research: %s, planetId: %s, researchName: %s\n", id, planetId, researchName)
	return nil
}
