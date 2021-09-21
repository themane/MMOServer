package mongoRepository

import (
	"context"
	"github.com/themane/MMOServer/mongoRepository/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type ClanRepositoryImpl struct {
	mongoURL   string
	ctx        context.Context
	cancelFunc context.CancelFunc
	mongoDB    string
}

func NewClanRepository(mongoURL string, mongoDB string) *ClanRepositoryImpl {
	ctx, cancelFunc := context.WithTimeout(context.Background(), connectTimeoutSecs*time.Second)
	return &ClanRepositoryImpl{
		mongoURL:   mongoURL,
		ctx:        ctx,
		cancelFunc: cancelFunc,
		mongoDB:    mongoDB,
	}
}

func (c *ClanRepositoryImpl) getMongoClient() *mongo.Client {
	return getConnection(c.mongoURL, c.ctx)
}

func (c *ClanRepositoryImpl) getCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(c.mongoDB).Collection("clan_data")
}

func (c *ClanRepositoryImpl) FindById(id string) (*models.ClanData, error) {
	client := c.getMongoClient()
	defer disconnect(client, c.ctx)
	var result models.ClanData
	filter := bson.M{"_id": id}
	singleResult := c.getCollection(client).FindOne(c.ctx, filter)
	err := singleResult.Decode(&result)
	if err != nil {
		log.Printf("Error in decoding clan data received from Mongo: %#v\n", err)
		return nil, err
	}
	return &result, nil
}
