package mongoRepository

import (
	"context"
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserRepositoryImpl struct {
	mongoURL string
	mongoDB  string
	logger   *constants.LoggingUtils
}

func NewUserRepository(mongoURL string, mongoDB string, logLevel string) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		mongoURL: mongoURL,
		mongoDB:  mongoDB,
		logger:   constants.NewLoggingUtils("USER_REPOSITORY", logLevel),
	}
}
func (u *UserRepositoryImpl) getMongoClient() (*mongo.Client, context.Context) {
	return getConnection(u.mongoURL)
}

func (u *UserRepositoryImpl) getCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(u.mongoDB).Collection("user_data")
}

func (u *UserRepositoryImpl) FindById(id string) (*repoModels.UserData, error) {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	result := repoModels.UserData{}
	filter := bson.M{"_id": id}
	singleResult := u.getCollection(client).FindOne(ctx, filter)
	err := singleResult.Decode(&result)
	if err != nil {
		u.logger.Error("Error in decoding user data received from Mongo", err)
		return nil, err
	}
	return &result, nil
}

func (u *UserRepositoryImpl) FindByUsername(username string) (*repoModels.UserData, error) {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	result := repoModels.UserData{}
	filter := bson.M{"profile.username": username}
	singleResult := u.getCollection(client).FindOne(ctx, filter)
	err := singleResult.Decode(&result)
	if err != nil {
		u.logger.Error("Error in decoding user data received from Mongo", err)
		return nil, err
	}
	return &result, nil
}

func (u *UserRepositoryImpl) AddExperience(id string, experience int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"profile.experience": experience}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Added experience id: %s, experience: %s\n", id, experience)
	return nil
}

func (u *UserRepositoryImpl) UpdateClanId(id string, clanId string) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"profile.clan_id": clanId}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Updated id: %s, clanId: %s\n", id, clanId)
	return nil
}

func (u *UserRepositoryImpl) UpgradeBuildingLevel(id string, planetId string, buildingId string,
	waterRequired int, grapheneRequired int, shelioRequired int, minutesRequired int) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"occupied_planets." + planetId + ".buildings." + buildingId + ".building_minutes_per_worker": minutesRequired,
		},
		"$inc": bson.M{
			"occupied_planets." + planetId + ".buildings." + buildingId + ".building_level": 1,
			"occupied_planets." + planetId + ".water.amount":                                -waterRequired,
			"occupied_planets." + planetId + ".graphene.amount":                             -grapheneRequired,
			"occupied_planets." + planetId + ".shelio":                                      -shelioRequired,
		},
	}
	result := u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	bsonResult, err := result.DecodeBytes()
	if err != nil {
		return err
	}
	u.logger.Printf(bsonResult.String())
	u.logger.Printf("Upgraded id: %s, planetId: %s, buildingId: %s\n", id, planetId, buildingId)
	return nil
}

func (u *UserRepositoryImpl) UpdateWorkers(id string, planetId string, buildingId string, workers int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".buildings." + buildingId + ".workers": workers,
		"occupied_planets." + planetId + ".population.workers":                   -workers,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Employed workers updated id: %s, planetId: %s, buildingId: %s, workers: %s\n", id, planetId, buildingId, workers)
	return nil
}

func (u *UserRepositoryImpl) UpdatePopulationRate(id string, planetId string, generationRate int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"occupied_planets." + planetId + ".population.generation_rate": generationRate,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Updated population generation rate id: %s, planetId: %s, rate: %s\n", id, planetId, generationRate)
	return nil
}

func (u *UserRepositoryImpl) Recruit(id string, planetId string, workers int, soldiers int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".population.unemployed": -(workers + soldiers),
		"occupied_planets." + planetId + ".population.workers":    workers,
		"occupied_planets." + planetId + ".population.soldiers":   soldiers,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Assigned workers and soldiers id: %s, planetId: %s, workers: %s, soldiers: %s\n", id, planetId, workers, soldiers)
	return nil
}

func (u *UserRepositoryImpl) KillPopulation(id string, planetId string, unemployed int, workers int, soldiers int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".population.unemployed": -unemployed,
		"occupied_planets." + planetId + ".population.workers":    -workers,
		"occupied_planets." + planetId + ".population.soldiers":   -soldiers,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Killed population: %s, planetId: %s, unemployed: %s, workers: %s, soldiers: %s\n", id, planetId, unemployed, workers, soldiers)
	return nil
}

func (u *UserRepositoryImpl) ReserveResources(id string, planetId string, water int, graphene int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".water.reserving":    water,
		"occupied_planets." + planetId + ".water.amount":       -water,
		"occupied_planets." + planetId + ".graphene.reserving": graphene,
		"occupied_planets." + planetId + ".graphene.amount":    -graphene,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Marked for reserving resources: %s, planetId: %s, water: %s, graphene: %s\n", id, planetId, water, graphene)
	return nil
}

func (u *UserRepositoryImpl) ExtractReservedResources(id string, planetId string, water int, graphene int) error {
	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".water.reserved":    -water,
		"occupied_planets." + planetId + ".water.amount":      water,
		"occupied_planets." + planetId + ".graphene.reserved": -graphene,
		"occupied_planets." + planetId + ".graphene.amount":   graphene,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Extracted reserved resources: %s, planetId: %s, water: %s, graphene: %s\n", id, planetId, water, graphene)
	return nil
}

func (u *UserRepositoryImpl) ConstructShips(id string, planetId string, unitName string, quantity float64,
	constructionRequirements models.Requirements) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"occupied_planets." + planetId + ".ships." + unitName + ".under_construction.start_time": primitive.NewDateTimeFromTime(time.Now()),
			"occupied_planets." + planetId + ".ships." + unitName + ".under_construction.quantity":   quantity,
		},
		"$inc": bson.M{
			"occupied_planets." + planetId + ".population.soldiers": -(constructionRequirements.SoldiersRequired * quantity),
			"occupied_planets." + planetId + ".population.workers":  -(constructionRequirements.WorkersRequired * quantity),
			"occupied_planets." + planetId + ".water.amount":        -(constructionRequirements.WaterRequired * quantity),
			"occupied_planets." + planetId + ".graphene.amount":     -(constructionRequirements.GrapheneRequired * quantity),
			"occupied_planets." + planetId + ".shelio":              -(constructionRequirements.ShelioRequired * quantity),
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Added %s %s ships for construction. id: %s, planetId: %s\n", quantity, unitName, id, planetId)
	return nil
}
func (u *UserRepositoryImpl) CancelShipsConstruction(id string, planetId string, unitName string,
	cancelReturns models.Returns) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$unset": bson.M{
			"occupied_planets." + planetId + ".ships." + unitName + ".under_construction": 1,
		},
		"$inc": bson.M{
			"occupied_planets." + planetId + ".population.soldiers": cancelReturns.SoldiersReturned,
			"occupied_planets." + planetId + ".population.workers":  cancelReturns.WorkersReturned,
			"occupied_planets." + planetId + ".water.amount":        cancelReturns.WaterReturned,
			"occupied_planets." + planetId + ".graphene.amount":     cancelReturns.GrapheneReturned,
			"occupied_planets." + planetId + ".shelio":              cancelReturns.ShelioReturned,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Canceled %s ships from construction. id: %s, planetId: %s\n", unitName, id, planetId)
	return nil
}
func (u *UserRepositoryImpl) DestructShips(id string, planetId string, unitName string, quantity float64,
	destructionReturns models.Returns) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".ships." + unitName + ".quantity": -quantity,
		"occupied_planets." + planetId + ".population.soldiers":             destructionReturns.SoldiersReturned * quantity,
		"occupied_planets." + planetId + ".population.workers":              destructionReturns.WorkersReturned * quantity,
		"occupied_planets." + planetId + ".water.amount":                    destructionReturns.WaterReturned * quantity,
		"occupied_planets." + planetId + ".graphene.amount":                 destructionReturns.GrapheneReturned * quantity,
		"occupied_planets." + planetId + ".shelio":                          destructionReturns.ShelioReturned * quantity,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Destructed %s %s ships. id: %s, planetId: %s\n", quantity, unitName, id, planetId)
	return nil
}

func (u *UserRepositoryImpl) ConstructDefences(id string, planetId string, unitName string, quantity float64,
	constructionRequirements models.Requirements) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"occupied_planets." + planetId + ".defences." + unitName + ".under_construction.start_time": primitive.NewDateTimeFromTime(time.Now()),
			"occupied_planets." + planetId + ".defences." + unitName + ".under_construction.quantity":   quantity,
		},
		"$inc": bson.M{
			"occupied_planets." + planetId + ".population.soldiers": -(constructionRequirements.SoldiersRequired * quantity),
			"occupied_planets." + planetId + ".population.workers":  -(constructionRequirements.WorkersRequired * quantity),
			"occupied_planets." + planetId + ".water.amount":        -(constructionRequirements.WaterRequired * quantity),
			"occupied_planets." + planetId + ".graphene.amount":     -(constructionRequirements.GrapheneRequired * quantity),
			"occupied_planets." + planetId + ".shelio":              -(constructionRequirements.ShelioRequired * quantity),
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Added %s %s ships for construction. id: %s, planetId: %s\n", quantity, unitName, id, planetId)
	return nil
}
func (u *UserRepositoryImpl) CancelDefencesConstruction(id string, planetId string, unitName string,
	cancelReturns models.Returns) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$unset": bson.M{
			"occupied_planets." + planetId + ".defences." + unitName + ".under_construction": 1,
		},
		"$inc": bson.M{
			"occupied_planets." + planetId + ".population.soldiers": cancelReturns.SoldiersReturned,
			"occupied_planets." + planetId + ".population.workers":  cancelReturns.WorkersReturned,
			"occupied_planets." + planetId + ".water.amount":        cancelReturns.WaterReturned,
			"occupied_planets." + planetId + ".graphene.amount":     cancelReturns.GrapheneReturned,
			"occupied_planets." + planetId + ".shelio":              cancelReturns.ShelioReturned,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Canceled %s defences from construction. id: %s, planetId: %s\n", unitName, id, planetId)
	return nil
}
func (u *UserRepositoryImpl) DestructDefences(id string, planetId string, unitName string, quantity float64,
	destructionReturns models.Returns) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets." + planetId + ".defences." + unitName + ".quantity": -quantity,
		"occupied_planets." + planetId + ".population.soldiers":                destructionReturns.SoldiersReturned * quantity,
		"occupied_planets." + planetId + ".population.workers":                 destructionReturns.WorkersReturned * quantity,
		"occupied_planets." + planetId + ".water.amount":                       destructionReturns.WaterReturned * quantity,
		"occupied_planets." + planetId + ".graphene.amount":                    destructionReturns.GrapheneReturned * quantity,
		"occupied_planets." + planetId + ".shelio":                             destructionReturns.ShelioReturned * quantity,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Destructed %s %s ships. id: %s, planetId: %s\n", quantity, unitName, id, planetId)
	return nil
}

func (u *UserRepositoryImpl) ConstructDefenceShipCarrier(id string, planetId string, unitName string, unitId string,
	constructionRequirements models.Requirements) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"occupied_planets." + planetId + ".defence_ship_carriers." + unitId + ".name":                          unitName,
			"occupied_planets." + planetId + ".defence_ship_carriers." + unitId + ".level":                         1,
			"occupied_planets." + planetId + ".defence_ship_carriers." + unitId + ".under_construction.start_time": primitive.NewDateTimeFromTime(time.Now()),
		},
		"$inc": bson.M{
			"occupied_planets." + planetId + ".population.soldiers": -constructionRequirements.SoldiersRequired,
			"occupied_planets." + planetId + ".population.workers":  -constructionRequirements.WorkersRequired,
			"occupied_planets." + planetId + ".water.amount":        -constructionRequirements.WaterRequired,
			"occupied_planets." + planetId + ".graphene.amount":     -constructionRequirements.GrapheneRequired,
			"occupied_planets." + planetId + ".shelio":              -constructionRequirements.ShelioRequired,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Added %s defence ship carrier for construction. id: %s, planetId: %s\n", unitName, id, planetId)
	return nil
}

func (u *UserRepositoryImpl) UpgradeDefenceShipCarrier(id string, planetId string, unitId string,
	constructionRequirements models.Requirements) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"occupied_planets." + planetId + ".defence_ship_carriers." + unitId + ".under_construction.start_time": primitive.NewDateTimeFromTime(time.Now()),
		},
		"$inc": bson.M{
			"occupied_planets." + planetId + ".defence_ship_carriers." + unitId + ".level": 1,
			"occupied_planets." + planetId + ".population.soldiers":                        -constructionRequirements.SoldiersRequired,
			"occupied_planets." + planetId + ".population.workers":                         -constructionRequirements.WorkersRequired,
			"occupied_planets." + planetId + ".water.amount":                               -constructionRequirements.WaterRequired,
			"occupied_planets." + planetId + ".graphene.amount":                            -constructionRequirements.GrapheneRequired,
			"occupied_planets." + planetId + ".shelio":                                     -constructionRequirements.ShelioRequired,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Added %s defence ship carrier for up-gradation. id: %s, planetId: %s\n", unitId, id, planetId)
	return nil
}

func (u *UserRepositoryImpl) CancelDefenceShipCarrierConstruction(id string, planetId string, unitId string,
	cancelReturns models.Returns) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$unset": bson.M{
			"occupied_planets." + planetId + ".defence_ship_carriers." + unitId: 1,
		},
		"$inc": bson.M{
			"occupied_planets." + planetId + ".population.soldiers": cancelReturns.SoldiersReturned,
			"occupied_planets." + planetId + ".population.workers":  cancelReturns.WorkersReturned,
			"occupied_planets." + planetId + ".water.amount":        cancelReturns.WaterReturned,
			"occupied_planets." + planetId + ".graphene.amount":     cancelReturns.GrapheneReturned,
			"occupied_planets." + planetId + ".shelio":              cancelReturns.ShelioReturned,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Canceled %s defence ship carrier up-gradation/construction. id: %s, planetId: %s\n", unitId, id, planetId)
	return nil
}

func (u *UserRepositoryImpl) CancelDefenceShipCarrierUpGradation(id string, planetId string, unitId string,
	cancelReturns models.Returns) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$unset": bson.M{
			"occupied_planets." + planetId + ".defence_ship_carriers." + unitId + ".under_construction": 1,
		},
		"$inc": bson.M{
			"occupied_planets." + planetId + ".defence_ship_carriers." + unitId + ".level": -1,
			"occupied_planets." + planetId + ".population.soldiers":                        cancelReturns.SoldiersReturned,
			"occupied_planets." + planetId + ".population.workers":                         cancelReturns.WorkersReturned,
			"occupied_planets." + planetId + ".water.amount":                               cancelReturns.WaterReturned,
			"occupied_planets." + planetId + ".graphene.amount":                            cancelReturns.GrapheneReturned,
			"occupied_planets." + planetId + ".shelio":                                     cancelReturns.ShelioReturned,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Canceled %s defence ship carrier up-gradation/construction. id: %s, planetId: %s\n", unitId, id, planetId)
	return nil
}

func (u *UserRepositoryImpl) DestructDefenceShipCarrier(id string, planetId string, unitId string,
	destructionReturns models.Returns) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$unset": bson.M{
			"occupied_planets." + planetId + ".defence_ship_carriers." + unitId: 1,
		},
		"$inc": bson.M{
			"occupied_planets." + planetId + ".population.soldiers": destructionReturns.SoldiersReturned,
			"occupied_planets." + planetId + ".population.workers":  destructionReturns.WorkersReturned,
			"occupied_planets." + planetId + ".water.amount":        destructionReturns.WaterReturned,
			"occupied_planets." + planetId + ".graphene.amount":     destructionReturns.GrapheneReturned,
			"occupied_planets." + planetId + ".shelio":              destructionReturns.ShelioReturned,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Destructed %s defence ship carrier. id: %s, planetId: %s\n", unitId, id, planetId)
	return nil
}

func (u *UserRepositoryImpl) DeployShipsOnDefenceShipCarrier(id string, planetId string, unitId string,
	ships map[string]int) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	shipsUpdateModel := bson.M{}
	for shipName, quantity := range ships {
		shipsUpdateModel["occupied_planets."+planetId+".defence_ship_carriers."+unitId+".hosting_ships."+shipName] = quantity
	}
	update := bson.M{
		"$set": shipsUpdateModel,
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Updated deployed ships on defence ship carrier. id: %s, planetId: %s, unitId: %s, ships: %s \n", id, planetId, unitId, ships)
	return nil
}

func (u *UserRepositoryImpl) DeployDefencesOnShield(id string, planetId string, shieldId string,
	defences map[string]int) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	defenceUpdateModel := bson.M{}
	for defenceName, quantity := range defences {
		defenceUpdateModel["occupied_planets."+planetId+".defences."+defenceName+".guarding_shield."+shieldId] = quantity
	}
	update := bson.M{
		"$set": defenceUpdateModel,
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Updated deployed defences on shield. id: %s, planetId: %s, shieldId: %s, defences: %s \n", id, planetId, shieldId, defences)
	return nil
}

func (u *UserRepositoryImpl) DeployDefenceShipCarrierOnShield(id string, planetId string, unitId string, shieldId string) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"occupied_planets." + planetId + ".defence_ship_carriers." + unitId + ".guarding_shield": shieldId,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update)
	u.logger.Printf("Destructed %s defence ship carrier. id: %s, planetId: %s\n", unitId, id, planetId)
	return nil
}
