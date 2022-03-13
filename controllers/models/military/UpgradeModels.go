package military

import (
	"github.com/themane/MMOServer/models"
	repoModels "github.com/themane/MMOServer/mongoRepository/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UnderConstruction struct {
	StartTime     primitive.DateTime `json:"start_time"`
	EndTime       primitive.DateTime `json:"end_time"`
	Quantity      int                `json:"quantity,omitempty"`
	CancelReturns models.Returns     `json:"cancel_returns"`
}

func InitUnderConstruction(underConstruction repoModels.UnderConstruction, levelConstants map[string]interface{}) *UnderConstruction {
	if underConstruction.Quantity > 0 {
		updatedTime := underConstruction.StartTime.Time().Add(time.Minute * time.Duration(int(levelConstants["minutes_required"].(float64))*underConstruction.Quantity))
		u := UnderConstruction{
			StartTime: underConstruction.StartTime,
			EndTime:   primitive.NewDateTimeFromTime(updatedTime),
			Quantity:  underConstruction.Quantity,
		}
		u.CancelReturns.InitCancelReturns(levelConstants, float64(underConstruction.Quantity))
		return &u
	}
	return nil
}

func InitUnderUpGradation(underConstruction repoModels.UnderConstruction, levelConstants map[string]interface{}) *UnderConstruction {
	if underConstruction.StartTime > 0 {
		updatedTime := underConstruction.StartTime.Time().Add(time.Minute * time.Duration(int(levelConstants["minutes_required"].(float64))))
		u := UnderConstruction{
			StartTime: underConstruction.StartTime,
			EndTime:   primitive.NewDateTimeFromTime(updatedTime),
			Quantity:  underConstruction.Quantity,
		}
		u.CancelReturns.InitCancelReturns(levelConstants, 1)
		return &u
	}
	return nil
}
