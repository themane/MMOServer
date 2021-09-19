package mongoRepository

import (
	"context"
	"github.com/themane/MMOServer/mongoRepository/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"math"
)

type UserRepositoryImpl struct {
	client     *mongo.Client
	ctx        context.Context
	cancelFunc context.CancelFunc
	mongoDB    string
}

func NewUserRepository(client *mongo.Client, ctx context.Context, cancelFunc context.CancelFunc, mongoDB string) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		client:     client,
		ctx:        ctx,
		cancelFunc: cancelFunc,
		mongoDB:    mongoDB,
	}
}

func (u *UserRepositoryImpl) getCollection() *mongo.Collection {
	return u.client.Database(u.mongoDB).Collection("user_data")
}

func (u *UserRepositoryImpl) FindById(id string) (*models.UserData, error) {
	defer disconnect(u.client, u.ctx)
	var result *models.UserData
	filter := bson.D{{"_id", id}}
	err := u.getCollection().FindOne(u.ctx, filter).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *UserRepositoryImpl) FindByUsername(id string) (*models.UserData, error) {
	defer disconnect(u.client, u.ctx)
	var result *models.UserData
	filter := bson.D{{"_id", id}}
	err := u.getCollection().FindOne(u.ctx, filter).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *UserRepositoryImpl) AddExperience(id string, experience int) error {
	defer disconnect(u.client, u.ctx)
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$inc", bson.D{{"profile.experience", experience}}},
	}
	u.getCollection().FindOneAndUpdate(u.ctx, filter, update)
	log.Printf("Added experience id: %s, experience: %d\n", id, experience)
	return nil
}

func (u *UserRepositoryImpl) UpdateClanId(id string, clanId string) error {
	defer disconnect(u.client, u.ctx)
	filter := bson.D{{"_id", id}}
	update := bson.D{{"profile.clan_id", clanId}}
	u.getCollection().FindOneAndUpdate(u.ctx, filter, update)
	log.Printf("Updated id: %s, clanId: %s\n", id, clanId)
	return nil
}

func (u *UserRepositoryImpl) UpgradeBuildingLevel(id string, planetId string, buildingId string,
	waterRequired int, grapheneRequired int, shelioRequired int) error {

	defer disconnect(u.client, u.ctx)
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$inc", bson.D{{"occupied_planets." + planetId + ".buildings." + buildingId, 1}}},
		{"$inc", bson.D{{"occupied_planets." + planetId + ".water.amount", -waterRequired}}},
		{"$inc", bson.D{{"occupied_planets." + planetId + ".graphene.", -grapheneRequired}}},
		{"$inc", bson.D{{"occupied_planets." + planetId + ".shelio", -shelioRequired}}},
	}
	u.getCollection().FindOneAndUpdate(u.ctx, filter, update)
	log.Printf("Upgraded id: %s, planetId: %s, buildingId: %s\n", id, planetId, buildingId)
	return nil
}

func (u *UserRepositoryImpl) AddResources(id string, planetId string, water int, graphene int, shelio int) error {
	defer disconnect(u.client, u.ctx)
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$inc", bson.D{{"occupied_planets." + planetId + ".water.amount", water}}},
		{"$inc", bson.D{{"occupied_planets." + planetId + ".graphene.amount", graphene}}},
		{"$inc", bson.D{{"occupied_planets." + planetId + ".shelio", shelio}}},
	}
	u.getCollection().FindOneAndUpdate(u.ctx, filter, update)
	log.Printf("Added id: %s, planetId: %s, water: %d, graphene: %d, shelio: %d\n", id, planetId, water, graphene, shelio)
	return nil
}

func (u *UserRepositoryImpl) UpdateMineResources(id string, planetId string, mineId string, water int, graphene int) error {
	defer disconnect(u.client, u.ctx)
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$inc", bson.D{{"occupied_planets." + planetId + ".water.amount", water}}},
		{"$inc", bson.D{{"occupied_planets." + planetId + ".water.amount", water}}},
		{"$inc", bson.D{{"occupied_planets." + planetId + ".graphene.amount", graphene}}},
		{"$inc", bson.D{{"occupied_planets." + planetId + ".mines." + mineId + "mined", math.Max(float64(water), float64(graphene))}}},
	}
	u.getCollection().FindOneAndUpdate(u.ctx, filter, update)
	log.Printf("Updated mine resources id: %s, planetId: %s, water: %d, graphene: %d, mineId: %s\n",
		id, planetId, water, graphene, mineId)
	return nil
}

func (u *UserRepositoryImpl) UpdateWorkers(id string, planetId string, buildingId string, workers int) error {
	defer disconnect(u.client, u.ctx)
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$set", bson.D{{"occupied_planets." + planetId + ".buildings." + buildingId + "workers", workers}}},
	}
	u.getCollection().FindOneAndUpdate(u.ctx, filter, update)
	log.Printf("Updated workers id: %s, planetId: %s, buildingId: %s, workers: %d\n", id, planetId, buildingId, workers)
	return nil
}

func (u *UserRepositoryImpl) AddPopulation(id string, planetId string, population int) error {
	defer disconnect(u.client, u.ctx)
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$inc", bson.D{{"occupied_planets." + planetId + ".population.unemployed", population}}},
	}
	u.getCollection().FindOneAndUpdate(u.ctx, filter, update)
	log.Printf("Added population id: %s, planetId: %s, population: %d\n", id, planetId, population)
	return nil
}

func (u *UserRepositoryImpl) RecruitWorkers(id string, planetId string, workers int) error {
	defer disconnect(u.client, u.ctx)
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$inc", bson.D{{"occupied_planets." + planetId + ".population.unemployed", -workers}}},
		{"$inc", bson.D{{"occupied_planets." + planetId + ".population.workers.idle", workers}}},
	}
	u.getCollection().FindOneAndUpdate(u.ctx, filter, update)
	log.Printf("Assigned workers id: %s, planetId: %s, workers: %d\n", id, planetId, workers)
	return nil
}

func (u *UserRepositoryImpl) RecruitSoldiers(id string, planetId string, soldiers int) error {
	defer disconnect(u.client, u.ctx)
	filter := bson.D{{"_id", id}}
	update := bson.D{
		{"$inc", bson.D{{"occupied_planets." + planetId + ".population.unemployed", -soldiers}}},
		{"$inc", bson.D{{"occupied_planets." + planetId + ".population.soldiers.idle", soldiers}}},
	}
	u.getCollection().FindOneAndUpdate(u.ctx, filter, update)
	log.Printf("Assigned soldiers id: %s, planetId: %s, soldiers: %d\n", id, planetId, soldiers)
	return nil
}