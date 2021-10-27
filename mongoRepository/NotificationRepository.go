package mongoRepository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationRepositoryImpl struct {
	mongoURL string
	mongoDB  string
	logger   *constants.LoggingUtils
}

func NewNotificationRepository(mongoURL string, mongoDB string, logLevel string) *NotificationRepositoryImpl {
	return &NotificationRepositoryImpl{
		mongoURL: mongoURL,
		mongoDB:  mongoDB,
		logger:   constants.NewLoggingUtils("NOTIFICATION_REPOSITORY", logLevel),
	}
}

func (n *NotificationRepositoryImpl) getMongoClient() (*mongo.Client, context.Context) {
	return getConnection(n.mongoURL)
}

func (n *NotificationRepositoryImpl) getCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(n.mongoDB).Collection("notification_data")
}

func (n *NotificationRepositoryImpl) FindByUserId(userId string) ([]repoModels.NotificationData, error) {
	client, ctx := n.getMongoClient()
	defer disconnect(client, ctx)
	var result []repoModels.NotificationData
	filter := bson.M{"user_id": userId}
	cursor, err := n.getCollection(client).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var notification repoModels.NotificationData
		err := cursor.Decode(&notification)
		if err != nil {
			n.logger.Error("Error in decoding notification data received from Mongo", err)
			return nil, err
		}
		result = append(result, notification)
	}
	return result, nil
}

func (n *NotificationRepositoryImpl) AddNotification(notification models.Notification, userId string) (*repoModels.NotificationData, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		n.logger.Error("error in persisting notification data", err)
		return nil, errors.New("error in persisting notification data")
	}
	notificationData := repoModels.NotificationData{
		Id:           id.String(),
		Notification: notification,
		UserId:       userId,
		Sent:         false,
	}
	client, ctx := n.getMongoClient()
	defer disconnect(client, ctx)
	_, err = n.getCollection(client).InsertOne(ctx, notificationData)
	if err != nil {
		n.logger.Error("error in persisting notification data", err)
		return nil, errors.New("error in persisting notification data")
	}
	return &notificationData, nil
}

func (n *NotificationRepositoryImpl) MarkSent(id string) error {
	client, ctx := n.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"sent": true}}
	n.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	n.logger.Printf("Marked notification as sent, id: %s\n", id)
	return nil
}
