package mongoRepository

import (
	"context"
	"github.com/themane/MMOServer/mongoRepository/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

type ClanRepositoryImpl struct {
	client     *mongo.Client
	ctx        context.Context
	cancelFunc context.CancelFunc
	mongoDB    string
}

func NewClanRepository(client *mongo.Client, ctx context.Context, cancelFunc context.CancelFunc, mongoDB string) *ClanRepositoryImpl {
	return &ClanRepositoryImpl{
		client:     client,
		ctx:        ctx,
		cancelFunc: cancelFunc,
		mongoDB:    mongoDB,
	}
}

func (u *ClanRepositoryImpl) getCollection() *mongo.Collection {
	return u.client.Database(u.mongoDB).Collection("user_data")
}

func (u *ClanRepositoryImpl) FindById(id uuid.UUID) (*models.ClanData, error) {
	defer disconnect(u.client, u.ctx)
	var result *models.ClanData
	filter := bson.D{{"_id", id}}
	err := u.getCollection().FindOne(u.ctx, filter).Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
