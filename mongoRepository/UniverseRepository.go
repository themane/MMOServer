package mongoRepository

import (
	"context"
	"github.com/themane/MMOServer/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type UniverseRepositoryImpl struct {
	client     *mongo.Client
	ctx        context.Context
	cancelFunc context.CancelFunc
	mongoDB    string
}

func NewUniverseRepository(client *mongo.Client, ctx context.Context, cancelFunc context.CancelFunc, mongoDB string) *UniverseRepositoryImpl {
	return &UniverseRepositoryImpl{
		client:     client,
		ctx:        ctx,
		cancelFunc: cancelFunc,
		mongoDB:    mongoDB,
	}
}

func (u *UniverseRepositoryImpl) getCollection() *mongo.Collection {
	return u.client.Database(u.mongoDB).Collection("universe")
}

func (u *UniverseRepositoryImpl) GetSector(system int, sector int) (map[string]models.PlanetUni, error) {
	defer disconnect(u.client, u.ctx)
	var result map[string]models.PlanetUni
	cursor, err := u.getCollection().Find(u.ctx, bson.D{{"position.system", system}, {"position.sector", sector}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(u.ctx) {
		var planet models.PlanetUni
		err := cursor.Decode(&planet)
		if err != nil {
			return nil, err
		}
		result[planet.Position.PlanetId()] = planet
	}
	return result, nil
}

func (u *UniverseRepositoryImpl) GetPlanet(system int, sector int, planet int) (*models.PlanetUni, error) {
	defer disconnect(u.client, u.ctx)
	var result *models.PlanetUni
	filter := bson.D{{"position.system", system}, {"position.sector", sector}, {"position.planet", planet}}
	err := u.getCollection().FindOne(u.ctx, filter).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *UniverseRepositoryImpl) MarkOccupied(system int, sector int, planet int) error {
	defer disconnect(u.client, u.ctx)
	filter := bson.D{{"position.system", system}, {"position.sector", sector}, {"position.planet", planet}}
	update := bson.D{{"occupied", true}}
	u.getCollection().FindOneAndUpdate(u.ctx, filter, update)
	log.Printf("Marked planet: %s as occupied\n", models.PlanetId(system, sector, planet))
	return nil
}
