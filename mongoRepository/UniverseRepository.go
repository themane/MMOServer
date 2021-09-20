package mongoRepository

import (
	"context"
	"errors"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type UniverseRepositoryImpl struct {
	mongoURL   string
	ctx        context.Context
	cancelFunc context.CancelFunc
	mongoDB    string
}

func NewUniverseRepository(mongoURL string, ctx context.Context, cancelFunc context.CancelFunc, mongoDB string) *UniverseRepositoryImpl {
	return &UniverseRepositoryImpl{
		mongoURL:   mongoURL,
		ctx:        ctx,
		cancelFunc: cancelFunc,
		mongoDB:    mongoDB,
	}
}

func (u *UniverseRepositoryImpl) getMongoClient() *mongo.Client {
	return getConnection(u.mongoURL, u.ctx)
}

func (u *UniverseRepositoryImpl) getCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(u.mongoDB).Collection("universe")
}

func (u *UniverseRepositoryImpl) GetSector(system int, sector int) (map[string]repoModels.PlanetUni, error) {
	client := u.getMongoClient()
	defer disconnect(client, u.ctx)
	var result map[string]repoModels.PlanetUni
	filter := bson.D{
		{"position.system", system},
		{"position.sector", sector},
	}
	cursor, err := u.getCollection(client).Find(u.ctx, filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(u.ctx) {
		var planet repoModels.PlanetUni
		err := cursor.Decode(&planet)
		if err != nil {
			return nil, err
		}
		result[planet.Position.PlanetId()] = planet
	}
	return result, nil
}

func (u *UniverseRepositoryImpl) GetPlanet(system int, sector int, planet int) (*repoModels.PlanetUni, error) {
	client := u.getMongoClient()
	defer disconnect(client, u.ctx)
	var result *repoModels.PlanetUni
	filter := bson.D{
		{"position.system", system},
		{"position.sector", sector},
		{"position.planet", planet},
	}
	err := u.getCollection(client).FindOne(u.ctx, filter).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *UniverseRepositoryImpl) GetAllOccupiedPlanets(system int) (map[string]repoModels.PlanetUni, error) {
	client := u.getMongoClient()
	defer disconnect(client, u.ctx)
	var result map[string]repoModels.PlanetUni
	filter := bson.D{
		{"position.system", system},
		{"occupied", bson.D{{"$ne", nil}}},
	}
	cursor, err := u.getCollection(client).Find(u.ctx, filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(u.ctx) {
		var planet repoModels.PlanetUni
		err := cursor.Decode(&planet)
		if err != nil {
			return nil, err
		}
		result[planet.Position.PlanetId()] = planet
	}
	return result, nil
}

func (u *UniverseRepositoryImpl) GetRandomUnoccupiedPlanet(system int) (*repoModels.PlanetUni, error) {
	client := u.getMongoClient()
	defer disconnect(client, u.ctx)
	filter := bson.D{
		{"position.system", system},
		{"occupied", nil},
	}
	randomChoicePipeline := bson.D{
		{"$match", filter},
		{"$sample", bson.E{Key: "size", Value: 1}},
	}
	cursor, err := u.getCollection(client).Aggregate(u.ctx, randomChoicePipeline)
	if err != nil {
		return nil, err
	}
	for cursor.Next(u.ctx) {
		var planet repoModels.PlanetUni
		err := cursor.Decode(&planet)
		if err != nil {
			return nil, err
		}
		return &planet, nil
	}
	return nil, errors.New("no new planet could be assigned")
}

func (u *UniverseRepositoryImpl) MarkOccupied(system int, sector int, planet int, userId string) error {
	client := u.getMongoClient()
	defer disconnect(client, u.ctx)
	filter := bson.D{{"position.system", system}, {"position.sector", sector}, {"position.planet", planet}}
	update := bson.D{{"occupied", userId}}
	u.getCollection(client).FindOneAndUpdate(u.ctx, filter, update)
	log.Printf("Marked planet: %s as occupied\n", models.PlanetId(system, sector, planet))
	return nil
}
