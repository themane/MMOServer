package mongoRepository

import (
	"context"
	"github.com/themane/MMOServer/mongoRepository/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type ClanRepositoryImpl struct {
	mongoURL string
	mongoDB  string
}

func NewClanRepository(mongoURL string, mongoDB string) *ClanRepositoryImpl {
	return &ClanRepositoryImpl{
		mongoURL: mongoURL,
		mongoDB:  mongoDB,
	}
}

func (c *ClanRepositoryImpl) getMongoClient() (*mongo.Client, context.Context) {
	return getConnection(c.mongoURL)
}

func (c *ClanRepositoryImpl) getCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(c.mongoDB).Collection("clan_data")
}

func (c *ClanRepositoryImpl) FindById(id string) (*models.ClanData, error) {
	client, ctx := c.getMongoClient()
	defer disconnect(client, ctx)
	result := models.ClanData{}
	filter := bson.M{"_id": id}
	singleResult := c.getCollection(client).FindOne(ctx, filter)
	err := singleResult.Decode(&result)
	if err != nil {
		log.Printf("Error in decoding clan data received from Mongo: %#v\n", err)
		return nil, err
	}
	return &result, nil
}
