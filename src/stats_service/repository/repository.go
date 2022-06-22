package repository

import (
	"context"
	"fmt"
	"stats_service/connections"
	"stats_service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateStatistics(field string, statistics models.VoterStats, groupType, groupName string, minAge int, maxAge int) error {
	client := connections.GetInstanceMongoClient()
	db := client.Database("statistics")
	collection := db.Collection("votes_stats")
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

	updateDocument := bson.M{"$inc": bson.M{field: 1}}
	updateOptions := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.TODO(), query, updateDocument, updateOptions)
	if err != nil {
		return fmt.Errorf("error updating statistics: %v", err)
	}
	return nil

}
