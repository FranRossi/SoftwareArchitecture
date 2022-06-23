package repository

import (
	"context"
	"fmt"
	"os"
	"stats_service/connections"
	"stats_service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DropDataBases() {
	client := connections.GetInstanceMongoClient()
	statsDataBase := client.Database(os.Getenv("DB"))
	statsDataBase.Drop(context.TODO())
}

func UpdateActualStatistics(statistics models.VoterStats, groupType, groupName string, minAge int, maxAge int) error {
	client := connections.GetInstanceMongoClient()
	db := client.Database(os.Getenv("DB"))
	collection := db.Collection(os.Getenv("COL"))
	var query bson.M
	if groupType == "region" {
		query = bson.M{
			"election_id": statistics.ElectionId,
			"group_type":  groupType,
			"group_name":  groupName,
			"region":      statistics.Region,
			"max_age":     maxAge,
			"min_age":     minAge,
			"sex":         statistics.Sex,
		}
	} else {
		query = bson.M{
			"election_id": statistics.ElectionId,
			"group_type":  groupType,
			"group_name":  groupName,
			"circuit":     statistics.Circuit,
			"max_age":     maxAge,
			"min_age":     minAge,
			"sex":         statistics.Sex,
		}
	}

	updateDocument := bson.M{"$inc": bson.M{"votes": 1}}
	updateOptions := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.TODO(), query, updateDocument, updateOptions)
	if err != nil {
		return fmt.Errorf("error updating statistics: %v", err)
	}
	return nil

}

func UpdateTotalStatistics(statistics models.VoterStats, groupType, groupName string, minAge int, maxAge int) error {
	client := connections.GetInstanceMongoClient()
	db := client.Database(os.Getenv("DB"))
	collection := db.Collection(os.Getenv("COL"))
	var query bson.M
	if groupType == "region" {
		query = bson.M{
			"election_id": statistics.ElectionId,
			"group_type":  groupType,
			"group_name":  groupName,
			"region":      statistics.Region,
			"max_age":     maxAge,
			"min_age":     minAge,
			"sex":         statistics.Sex,
			"votes":       0,
		}
	} else {
		query = bson.M{
			"election_id": statistics.ElectionId,
			"group_type":  groupType,
			"group_name":  groupName,
			"circuit":     statistics.Circuit,
			"max_age":     maxAge,
			"min_age":     minAge,
			"sex":         statistics.Sex,
			"votes":       0,
		}
	}

	updateDocument := bson.M{"$inc": bson.M{"capacity": 1}}
	updateOptions := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.TODO(), query, updateDocument, updateOptions)
	if err != nil {
		return fmt.Errorf("error updating statistics: %v", err)
	}
	return nil

}
