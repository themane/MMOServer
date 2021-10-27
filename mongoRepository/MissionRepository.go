package mongoRepository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MissionRepositoryImpl struct {
	mongoURL string
	mongoDB  string
	logger   *constants.LoggingUtils
}

func NewMissionRepository(mongoURL string, mongoDB string, logLevel string) *MissionRepositoryImpl {
	return &MissionRepositoryImpl{
		mongoURL: mongoURL,
		mongoDB:  mongoDB,
		logger:   constants.NewLoggingUtils("MISSION_REPOSITORY", logLevel),
	}
}

func (c *MissionRepositoryImpl) getMongoClient() (*mongo.Client, context.Context) {
	return getConnection(c.mongoURL)
}

func (c *MissionRepositoryImpl) getCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(c.mongoDB).Collection("mission_data")
}

func (c *MissionRepositoryImpl) FindAttackMissionsFromPlanetId(fromPlanetId string) ([]repoModels.AttackMission, error) {
	client, ctx := c.getMongoClient()
	defer disconnect(client, ctx)
	var result []repoModels.AttackMission
	filter := bson.M{"from_planet_id": fromPlanetId, "mission_type": constants.AttackMission}
	cursor, err := c.getCollection(client).Find(ctx, filter)
	if err != nil {
		c.logger.Error("error in retrieving attack missions from DB, fromPlanetId: "+fromPlanetId, err)
		return nil, err
	}
	err = cursor.All(ctx, &result)
	if err != nil {
		c.logger.Error("error in retrieving attack missions from DB, fromPlanetId: "+fromPlanetId, err)
		return nil, err
	}
	return result, nil
}
func (c *MissionRepositoryImpl) FindSpyMissionsFromPlanetId(fromPlanetId string) ([]repoModels.SpyMission, error) {
	client, ctx := c.getMongoClient()
	defer disconnect(client, ctx)
	var result []repoModels.SpyMission
	filter := bson.M{"from_planet_id": fromPlanetId, "mission_type": constants.SpyMission}
	cursor, err := c.getCollection(client).Find(ctx, filter)
	if err != nil {
		c.logger.Error("error in retrieving spy missions from DB, fromPlanetId: "+fromPlanetId, err)
		return nil, err
	}
	err = cursor.All(ctx, &result)
	if err != nil {
		c.logger.Error("error in retrieving spy missions from DB, fromPlanetId: "+fromPlanetId, err)
		return nil, err
	}
	return result, nil
}

func (c *MissionRepositoryImpl) FindAttackMissionsToPlanetId(toPlanetId string) ([]repoModels.AttackMission, error) {
	client, ctx := c.getMongoClient()
	defer disconnect(client, ctx)
	var result []repoModels.AttackMission
	filter := bson.M{"to_planet_id": toPlanetId, "mission_type": constants.AttackMission}
	cursor, err := c.getCollection(client).Find(ctx, filter)
	if err != nil {
		c.logger.Error("error in retrieving attack missions from DB, toPlanetId: "+toPlanetId, err)
		return nil, err
	}
	err = cursor.All(ctx, &result)
	if err != nil {
		c.logger.Error("error in retrieving attack missions from DB, toPlanetId: "+toPlanetId, err)
		return nil, err
	}
	return result, nil
}
func (c *MissionRepositoryImpl) FindSpyMissionsToPlanetId(toPlanetId string) ([]repoModels.SpyMission, error) {
	client, ctx := c.getMongoClient()
	defer disconnect(client, ctx)
	var result []repoModels.SpyMission
	filter := bson.M{"to_planet_id": toPlanetId, "mission_type": constants.SpyMission}
	cursor, err := c.getCollection(client).Find(ctx, filter)
	if err != nil {
		c.logger.Error("error in retrieving spy missions from DB, toPlanetId: "+toPlanetId, err)
		return nil, err
	}
	err = cursor.All(ctx, &result)
	if err != nil {
		c.logger.Error("error in retrieving spy missions from DB, toPlanetId: "+toPlanetId, err)
		return nil, err
	}
	return result, nil
}

func (c *MissionRepositoryImpl) AddAttackMission(fromPlanetId string, toPlanetId string, formation map[string]map[string][]models.Formation,
	launchTime primitive.DateTime, missionTime primitive.DateTime, returnTime primitive.DateTime,
) (*repoModels.AttackMission, error) {

	id, err := uuid.NewRandom()
	if err != nil {
		c.logger.Error("error in persisting attack mission", err)
		return nil, errors.New("error in persisting attack mission")
	}
	attackMission := repoModels.AttackMission{
		Id:           id.String(),
		FromPlanetId: fromPlanetId,
		ToPlanetId:   toPlanetId,
		Formation:    formation,
		LaunchTime:   launchTime,
		MissionTime:  missionTime,
		ReturnTime:   returnTime,
		State:        constants.DepartureState,
		MissionType:  constants.AttackMission,
	}
	client, ctx := c.getMongoClient()
	defer disconnect(client, ctx)
	_, err = c.getCollection(client).InsertOne(ctx, attackMission)
	if err != nil {
		c.logger.Error("error in persisting attack mission", err)
		return nil, errors.New("error in persisting attack mission")
	}
	return &attackMission, nil
}
func (c *MissionRepositoryImpl) AddSpyMission(fromPlanetId string, toPlanetId string, scouts map[string]int,
	launchTime primitive.DateTime, missionTime primitive.DateTime, returnTime primitive.DateTime,
) (*repoModels.SpyMission, error) {

	id, err := uuid.NewRandom()
	if err != nil {
		c.logger.Error("error in persisting spy mission", err)
		return nil, errors.New("error in persisting spy mission")
	}
	spyMission := repoModels.SpyMission{
		Id:           id.String(),
		FromPlanetId: fromPlanetId,
		ToPlanetId:   toPlanetId,
		Scouts:       scouts,
		LaunchTime:   launchTime,
		MissionTime:  missionTime,
		ReturnTime:   returnTime,
		State:        constants.DepartureState,
		MissionType:  constants.SpyMission,
	}
	client, ctx := c.getMongoClient()
	defer disconnect(client, ctx)
	_, err = c.getCollection(client).InsertOne(ctx, spyMission)
	if err != nil {
		c.logger.Error("error in persisting spy mission: ", err)
		return nil, errors.New("error in persisting spy mission")
	}
	return &spyMission, nil
}

func (c *MissionRepositoryImpl) UpdateAttackResult(id string, result repoModels.AttackResult) error {
	client, ctx := c.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"result": result}}
	_, err := c.getCollection(client).UpdateOne(ctx, filter, update)
	if err != nil {
		c.logger.Error("error in updating attack mission", err)
		return errors.New("error in updating attack mission")
	}
	return nil
}
func (c *MissionRepositoryImpl) UpdateSpyResult(id string, result repoModels.SpyResult) error {
	client, ctx := c.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"result": result}}
	_, err := c.getCollection(client).UpdateOne(ctx, filter, update)
	if err != nil {
		c.logger.Error("error in updating spy mission", err)
		return errors.New("error in updating spy mission")
	}
	return nil
}
func (c *MissionRepositoryImpl) UpdateMissionState(id string, state string) error {
	client, ctx := c.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"state": state}}
	_, err := c.getCollection(client).UpdateOne(ctx, filter, update)
	if err != nil {
		c.logger.Error("error in updating mission state", err)
		return errors.New("error in updating mission state")
	}
	return nil
}
