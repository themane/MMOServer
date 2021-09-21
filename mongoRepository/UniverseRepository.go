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
	mongoURL string
	mongoDB  string
}

func NewUniverseRepository(mongoURL string, mongoDB string) *UniverseRepositoryImpl {
	return &UniverseRepositoryImpl{
		mongoURL: mongoURL,
		mongoDB:  mongoDB,
	}
}

func (u *UniverseRepositoryImpl) getMongoClient() (*mongo.Client, context.Context) {
	return getConnection(u.mongoURL)
}

func (u *UniverseRepositoryImpl) getCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(u.mongoDB).Collection("universe")
}

func (u *UniverseRepositoryImpl) FindById(id string) (*repoModels.PlanetUni, error) {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	result := repoModels.PlanetUni{}
	filter := bson.M{"_id": id}
	singleResult := u.getCollection(client).FindOne(ctx, filter)
	err := singleResult.Decode(&result)
	if err != nil {
		log.Printf("Error in decoding planet data received from Mongo: %#v\n", err)
		return nil, err
	}
	return &result, nil
}

func (u *UniverseRepositoryImpl) FindByPosition(system int, sector int, planet int) (*repoModels.PlanetUni, error) {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	result := repoModels.PlanetUni{}
	filter := bson.M{
		"position.system": system,
		"position.sector": sector,
		"position.planet": planet,
	}
	singleResult := u.getCollection(client).FindOne(ctx, filter)
	err := singleResult.Decode(&result)
	if err != nil {
		log.Printf("Error in decoding planet data received from Mongo: %#v\n", err)
		return nil, err
	}
	return &result, nil
}

func (u *UniverseRepositoryImpl) GetSector(system int, sector int) (map[string]repoModels.PlanetUni, error) {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	result := map[string]repoModels.PlanetUni{}
	filter := bson.M{
		"position.system": system,
		"position.sector": sector,
	}
	cursor, err := u.getCollection(client).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var planet repoModels.PlanetUni
		err := cursor.Decode(&planet)
		if err != nil {
			log.Printf("Error in decoding planet data received from Mongo: %#v\n", err)
			return nil, err
		}
		result[planet.Position.PlanetId()] = planet
	}
	return result, nil
}

func (u *UniverseRepositoryImpl) GetAllOccupiedPlanets(system int) (map[string]repoModels.PlanetUni, error) {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	result := map[string]repoModels.PlanetUni{}
	filter := bson.M{
		"position.system": system,
		"occupied":        bson.M{"$ne": nil},
	}
	cursor, err := u.getCollection(client).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var planet repoModels.PlanetUni
		err := cursor.Decode(&planet)
		if err != nil {
			log.Printf("Error in decoding planet data received from Mongo: %#v\n", err)
			return nil, err
		}
		result[planet.Position.PlanetId()] = planet
	}
	return result, nil
}

func (u *UniverseRepositoryImpl) GetRandomUnoccupiedPlanet(system int) (*repoModels.PlanetUni, error) {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{
		"position.system": system,
		"occupied":        nil,
	}
	randomChoicePipeline := bson.M{
		"$match":  filter,
		"$sample": bson.M{"size": 1},
	}
	cursor, err := u.getCollection(client).Aggregate(ctx, randomChoicePipeline)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		planet := repoModels.PlanetUni{}
		err := cursor.Decode(&planet)
		if err != nil {
			log.Printf("Error in decoding planet data received from Mongo: %#v\n", err)
			return nil, err
		}
		return &planet, nil
	}
	return nil, errors.New("no new planet could be assigned")
}

func (u *UniverseRepositoryImpl) MarkOccupied(system int, sector int, planet int, userId string) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{
		"position.system": system,
		"position.sector": sector,
		"position.planet": planet,
	}
	update := bson.M{"occupied": userId}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	log.Printf("Marked planet: %s as occupied\n", models.PlanetId(system, sector, planet))
	return nil
}
