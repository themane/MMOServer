package mongoRepository

import (
	"context"
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/mongoRepository/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClanRepositoryImpl struct {
	mongoURL string
	mongoDB  string
	logger   *constants.LoggingUtils
}

func NewClanRepository(mongoURL string, mongoDB string, logLevel string) *ClanRepositoryImpl {
	return &ClanRepositoryImpl{
		mongoURL: mongoURL,
		mongoDB:  mongoDB,
		logger:   constants.NewLoggingUtils("CLAN_REPOSITORY", logLevel),
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
		c.logger.Error("error in decoding retrieved clan data", err)
		return nil, err
	}
	return &result, nil
}
