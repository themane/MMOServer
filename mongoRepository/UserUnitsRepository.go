package mongoRepository

import (
	"github.com/themane/MMOServer/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

func (u *UserRepositoryImpl) ConstructShips(id string, planetId string, unitName string, quantity float64,
	constructionRequirements models.Requirements) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"occupied_planets.$[planetElement].ships.$[unitElement].under_construction.start_time": primitive.NewDateTimeFromTime(time.Now()),
			"occupied_planets.$[planetElement].ships.$[unitElement].under_construction.quantity":   quantity,
		},
		"$inc": bson.M{
			"occupied_planets.$[planetElement].population.soldiers": -(constructionRequirements.Population.Soldiers * quantity),
			"occupied_planets.$[planetElement].population.workers":  -(constructionRequirements.Population.Workers * quantity),
			"occupied_planets.$[planetElement].water.amount":        -(constructionRequirements.Resources.Water * quantity),
			"occupied_planets.$[planetElement].graphene.amount":     -(constructionRequirements.Resources.Graphene * quantity),
			"occupied_planets.$[planetElement].shelio":              -(constructionRequirements.Resources.Shelio * quantity),
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"planetElement._id": planetId},
				bson.M{"unitElement.name": unitName},
			},
		}),
	)
	u.logger.Printf("Added %s %s ships for construction. id: %s, planetId: %s\n", quantity, unitName, id, planetId)
	return nil
}

func (u *UserRepositoryImpl) AddConstructionShips(id string, planetId string, unitName string, quantity float64,
	constructionRequirements models.Requirements) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{
			"occupied_planets.$[planetElement].population.soldiers":                              -(constructionRequirements.Population.Soldiers * quantity),
			"occupied_planets.$[planetElement].population.workers":                               -(constructionRequirements.Population.Workers * quantity),
			"occupied_planets.$[planetElement].water.amount":                                     -(constructionRequirements.Resources.Water * quantity),
			"occupied_planets.$[planetElement].graphene.amount":                                  -(constructionRequirements.Resources.Graphene * quantity),
			"occupied_planets.$[planetElement].shelio":                                           -(constructionRequirements.Resources.Shelio * quantity),
			"occupied_planets.$[planetElement].ships.$[unitElement].under_construction.quantity": quantity,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"planetElement._id": planetId},
				bson.M{"unitElement.name": unitName},
			},
		}),
	)
	u.logger.Printf("Added %s %s ships for construction. id: %s, planetId: %s\n", quantity, unitName, id, planetId)
	return nil
}

func (u *UserRepositoryImpl) RemoveConstructionShips(id string, planetId string, unitName string, quantity float64,
	cancelReturns models.Returns) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{
			"occupied_planets.$[planetElement].population.soldiers":                              cancelReturns.Population.Soldiers * quantity,
			"occupied_planets.$[planetElement].population.workers":                               cancelReturns.Population.Workers * quantity,
			"occupied_planets.$[planetElement].water.amount":                                     cancelReturns.Resources.Water * quantity,
			"occupied_planets.$[planetElement].graphene.amount":                                  cancelReturns.Resources.Graphene * quantity,
			"occupied_planets.$[planetElement].shelio":                                           cancelReturns.Resources.Shelio * quantity,
			"occupied_planets.$[planetElement].ships.$[unitElement].under_construction.quantity": -quantity,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"planetElement._id": planetId},
				bson.M{"unitElement.name": unitName},
			},
		}),
	)
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
			"occupied_planets.$[planetElement].ships.$[unitElement].under_construction": 1,
		},
		"$inc": bson.M{
			"occupied_planets.$[planetElement].population.soldiers": cancelReturns.Population.Soldiers,
			"occupied_planets.$[planetElement].population.workers":  cancelReturns.Population.Workers,
			"occupied_planets.$[planetElement].water.amount":        cancelReturns.Resources.Water,
			"occupied_planets.$[planetElement].graphene.amount":     cancelReturns.Resources.Graphene,
			"occupied_planets.$[planetElement].shelio":              cancelReturns.Resources.Shelio,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"planetElement._id": planetId},
				bson.M{"unitElement.name": unitName},
			},
		}),
	)
	u.logger.Printf("Canceled %s ships from construction. id: %s, planetId: %s\n", unitName, id, planetId)
	return nil
}
func (u *UserRepositoryImpl) DestructShips(id string, planetId string, unitName string, quantity float64,
	destructionReturns models.Returns) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets.$[planetElement].ships.$[unitElement].quantity": -quantity,
		"occupied_planets.$[planetElement].population.soldiers":           destructionReturns.Population.Soldiers * quantity,
		"occupied_planets.$[planetElement].population.workers":            destructionReturns.Population.Workers * quantity,
		"occupied_planets.$[planetElement].water.amount":                  destructionReturns.Resources.Water * quantity,
		"occupied_planets.$[planetElement].graphene.amount":               destructionReturns.Resources.Graphene * quantity,
		"occupied_planets.$[planetElement].shelio":                        destructionReturns.Resources.Shelio * quantity,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"planetElement._id": planetId},
				bson.M{"unitElement.name": unitName},
			},
		}),
	)
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
			"occupied_planets.$[planetElement].defences.$[unitElement].under_construction.start_time": primitive.NewDateTimeFromTime(time.Now()),
			"occupied_planets.$[planetElement].defences.$[unitElement].under_construction.quantity":   quantity,
		},
		"$inc": bson.M{
			"occupied_planets.$[planetElement].population.soldiers": -(constructionRequirements.Population.Soldiers * quantity),
			"occupied_planets.$[planetElement].population.workers":  -(constructionRequirements.Population.Workers * quantity),
			"occupied_planets.$[planetElement].water.amount":        -(constructionRequirements.Resources.Water * quantity),
			"occupied_planets.$[planetElement].graphene.amount":     -(constructionRequirements.Resources.Graphene * quantity),
			"occupied_planets.$[planetElement].shelio":              -(constructionRequirements.Resources.Shelio * quantity),
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"planetElement._id": planetId},
				bson.M{"unitElement.name": unitName},
			},
		}),
	)
	u.logger.Printf("Added %s %s ships for construction. id: %s, planetId: %s\n", quantity, unitName, id, planetId)
	return nil
}

func (u *UserRepositoryImpl) AddConstructionDefences(id string, planetId string, unitName string, quantity float64,
	constructionRequirements models.Requirements) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{
			"occupied_planets.$[planetElement].population.soldiers":                                 -(constructionRequirements.Population.Soldiers * quantity),
			"occupied_planets.$[planetElement].population.workers":                                  -(constructionRequirements.Population.Workers * quantity),
			"occupied_planets.$[planetElement].water.amount":                                        -(constructionRequirements.Resources.Water * quantity),
			"occupied_planets.$[planetElement].graphene.amount":                                     -(constructionRequirements.Resources.Graphene * quantity),
			"occupied_planets.$[planetElement].shelio":                                              -(constructionRequirements.Resources.Shelio * quantity),
			"occupied_planets.$[planetElement].defences.$[unitElement].under_construction.quantity": quantity,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"planetElement._id": planetId},
				bson.M{"unitElement.name": unitName},
			},
		}),
	)
	u.logger.Printf("Added %s %s ships for construction. id: %s, planetId: %s\n", quantity, unitName, id, planetId)
	return nil
}

func (u *UserRepositoryImpl) RemoveConstructionDefences(id string, planetId string, unitName string, quantity float64,
	cancelReturns models.Returns) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$inc": bson.M{
			"occupied_planets.$[planetElement].population.soldiers":                                 cancelReturns.Population.Soldiers * quantity,
			"occupied_planets.$[planetElement].population.workers":                                  cancelReturns.Population.Workers * quantity,
			"occupied_planets.$[planetElement].water.amount":                                        cancelReturns.Resources.Water * quantity,
			"occupied_planets.$[planetElement].graphene.amount":                                     cancelReturns.Resources.Graphene * quantity,
			"occupied_planets.$[planetElement].shelio":                                              cancelReturns.Resources.Shelio * quantity,
			"occupied_planets.$[planetElement].defences.$[unitElement].under_construction.quantity": -quantity,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"planetElement._id": planetId},
				bson.M{"unitElement.name": unitName},
			},
		}),
	)
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
			"occupied_planets.$[planetElement].defences.$[unitElement].under_construction": 1,
		},
		"$inc": bson.M{
			"occupied_planets.$[planetElement].population.soldiers": cancelReturns.Population.Soldiers,
			"occupied_planets.$[planetElement].population.workers":  cancelReturns.Population.Workers,
			"occupied_planets.$[planetElement].water.amount":        cancelReturns.Resources.Water,
			"occupied_planets.$[planetElement].graphene.amount":     cancelReturns.Resources.Graphene,
			"occupied_planets.$[planetElement].shelio":              cancelReturns.Resources.Shelio,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"planetElement._id": planetId},
				bson.M{"unitElement.name": unitName},
			},
		}),
	)
	u.logger.Printf("Canceled %s defences from construction. id: %s, planetId: %s\n", unitName, id, planetId)
	return nil
}
func (u *UserRepositoryImpl) DestructDefences(id string, planetId string, unitName string, quantity float64,
	destructionReturns models.Returns) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{
		"occupied_planets.$[planetElement].defences.$[unitElement].quantity": -quantity,
		"occupied_planets.$[planetElement].population.soldiers":              destructionReturns.Population.Soldiers * quantity,
		"occupied_planets.$[planetElement].population.workers":               destructionReturns.Population.Workers * quantity,
		"occupied_planets.$[planetElement].water.amount":                     destructionReturns.Resources.Water * quantity,
		"occupied_planets.$[planetElement].graphene.amount":                  destructionReturns.Resources.Graphene * quantity,
		"occupied_planets.$[planetElement].shelio":                           destructionReturns.Resources.Shelio * quantity,
	}}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"planetElement._id": planetId},
				bson.M{"unitElement.name": unitName},
			},
		}),
	)
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
			"occupied_planets.$[planetElement].defence_ship_carriers.$[unitElement].name":                          unitName,
			"occupied_planets.$[planetElement].defence_ship_carriers.$[unitElement].level":                         1,
			"occupied_planets.$[planetElement].defence_ship_carriers.$[unitElement].under_construction.start_time": primitive.NewDateTimeFromTime(time.Now()),
		},
		"$inc": bson.M{
			"occupied_planets.$[planetElement].population.soldiers": -constructionRequirements.Population.Soldiers,
			"occupied_planets.$[planetElement].population.workers":  -constructionRequirements.Population.Workers,
			"occupied_planets.$[planetElement].water.amount":        -constructionRequirements.Resources.Water,
			"occupied_planets.$[planetElement].graphene.amount":     -constructionRequirements.Resources.Graphene,
			"occupied_planets.$[planetElement].shelio":              -constructionRequirements.Resources.Shelio,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"planetElement._id": planetId},
				bson.M{"unitElement._id": unitId},
			},
		}),
	)
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
			"occupied_planets.$[planetElement].defence_ship_carriers.$[unitElement].under_construction.start_time": primitive.NewDateTimeFromTime(time.Now()),
		},
		"$inc": bson.M{
			"occupied_planets.$[planetElement].defence_ship_carriers.$[unitElement].level": 1,
			"occupied_planets.$[planetElement].population.soldiers":                        -constructionRequirements.Population.Soldiers,
			"occupied_planets.$[planetElement].population.workers":                         -constructionRequirements.Population.Workers,
			"occupied_planets.$[planetElement].water.amount":                               -constructionRequirements.Resources.Water,
			"occupied_planets.$[planetElement].graphene.amount":                            -constructionRequirements.Resources.Graphene,
			"occupied_planets.$[planetElement].shelio":                                     -constructionRequirements.Resources.Shelio,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"planetElement._id": planetId},
				bson.M{"unitElement._id": unitId},
			},
		}),
	)
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
			"occupied_planets.$[planetElement].defence_ship_carriers.$[unitElement]": 1,
		},
		"$inc": bson.M{
			"occupied_planets.$[planetElement].population.soldiers": cancelReturns.Population.Soldiers,
			"occupied_planets.$[planetElement].population.workers":  cancelReturns.Population.Workers,
			"occupied_planets.$[planetElement].water.amount":        cancelReturns.Resources.Water,
			"occupied_planets.$[planetElement].graphene.amount":     cancelReturns.Resources.Graphene,
			"occupied_planets.$[planetElement].shelio":              cancelReturns.Resources.Shelio,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"planetElement._id": planetId},
				bson.M{"unitElement._id": unitId},
			},
		}),
	)
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
			"occupied_planets.$[planetElement].defence_ship_carriers.$[unitElement].under_construction": 1,
		},
		"$inc": bson.M{
			"occupied_planets.$[planetElement].defence_ship_carriers.$[unitElement].level": -1,
			"occupied_planets.$[planetElement].population.soldiers":                        cancelReturns.Population.Soldiers,
			"occupied_planets.$[planetElement].population.workers":                         cancelReturns.Population.Workers,
			"occupied_planets.$[planetElement].water.amount":                               cancelReturns.Resources.Water,
			"occupied_planets.$[planetElement].graphene.amount":                            cancelReturns.Resources.Graphene,
			"occupied_planets.$[planetElement].shelio":                                     cancelReturns.Resources.Shelio,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"planetElement._id": planetId},
				bson.M{"unitElement._id": unitId},
			},
		}),
	)
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
			"occupied_planets.$[planetElement].defence_ship_carriers.$[unitElement]": 1,
		},
		"$inc": bson.M{
			"occupied_planets.$[planetElement].population.soldiers": destructionReturns.Population.Soldiers,
			"occupied_planets.$[planetElement].population.workers":  destructionReturns.Population.Workers,
			"occupied_planets.$[planetElement].water.amount":        destructionReturns.Resources.Water,
			"occupied_planets.$[planetElement].graphene.amount":     destructionReturns.Resources.Graphene,
			"occupied_planets.$[planetElement].shelio":              destructionReturns.Resources.Shelio,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"planetElement._id": planetId},
				bson.M{"unitElement._id": unitId},
			},
		}),
	)
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
		shipsUpdateModel["occupied_planets.$[planetElement].defence_ship_carriers.$[unitElement].hosting_ships."+shipName] = quantity
	}
	update := bson.M{
		"$set": shipsUpdateModel,
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"planetElement._id": planetId},
				bson.M{"unitElement._id": unitId},
			},
		}),
	)
	u.logger.Printf("Updated deployed ships on defence ship carrier. id: %s, planetId: %s, unitId: %s, ships: %s \n", id, planetId, unitId, ships)
	return nil
}

func (u *UserRepositoryImpl) DeployDefencesOnShield(id string, planetId string, shieldId string,
	defences map[string]int) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	defenceUpdateModel := bson.M{}
	arrayFilters := []interface{}{
		bson.M{"planetElement._id": planetId},
	}
	index := 0
	for defenceName, quantity := range defences {
		defenceElementLabel := "defence" + strconv.Itoa(index)
		defenceUpdateModel["occupied_planets.$[planetElement].defences.$["+defenceElementLabel+"].guarding_shield."+shieldId] = quantity
		arrayFilters = append(arrayFilters, bson.M{defenceElementLabel + ".name": defenceName})
	}
	update := bson.M{
		"$set": defenceUpdateModel,
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{Filters: arrayFilters}))
	u.logger.Printf("Updated deployed defences on shield. id: %s, planetId: %s, shieldId: %s, defences: %s \n", id, planetId, shieldId, defences)
	return nil
}

func (u *UserRepositoryImpl) DeployDefenceShipCarrierOnShield(id string, planetId string, unitId string, shieldId string) error {

	client, ctx := u.getMongoClient()
	defer disconnect(client, ctx)
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"occupied_planets.$[planetElement].defence_ship_carriers.$[unitElement].guarding_shield": shieldId,
		},
	}
	u.getCollection(client).FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetArrayFilters(options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"planetElement._id": planetId},
				bson.M{"unitElement._id": unitId},
			},
		}),
	)
	u.logger.Printf("Destructed %s defence ship carrier. id: %s, planetId: %s\n", unitId, id, planetId)
	return nil
}
