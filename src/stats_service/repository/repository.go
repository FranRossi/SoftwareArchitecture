package repository

import (
	"context"
	"fmt"
	"voter_api/connections"
	"voter_api/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateStatistics(statistics domain.Statistics, groupType, groupName string) error {
	client := connections.GetInstanceMongoClient()
	uruguayDataBase := client.Database("uruguay_votes")
	uruguayanVotesCollection := uruguayDataBase.Collection("statistics")
	var query bson.M
	if groupType == "region" {
		query = bson.M{"election_id": statistics.ElectionId, "group_type": groupType, "group_name": groupName, "region": statistics.Region}
	} else {
		query = bson.M{"election_id": statistics.ElectionId, "group_type": groupType, "group_name": groupName, "circuit": statistics.Circuit}
	}

	updateDocument := bson.M{"$inc": bson.M{"votes": 1}}
	updateOptions := options.Update().SetUpsert(true)
	_, err := uruguayanVotesCollection.UpdateOne(context.TODO(), query, updateDocument, updateOptions)
	if err != nil {
		return fmt.Errorf("error updating statistics: %v", err)
	}
	return nil
}
