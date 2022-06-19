package repositories

import (
	"consulting_api/connections"
	m "consulting_api/models"
	"context"
	"encoding/json"
	l "own_logger"
	"time"
)

func RequestElectionConfigFromCache(electionId string) (m.ElectionConfig, error) {
	cache := connections.GetInstanceRedisClient()
	ctx := context.Background()
	var electionConfig m.ElectionConfig
	key := getElectionConfigKey(electionId)

	val, err := cache.Get(ctx, key).Bytes()
	if err != nil {
		l.LogInfo("Election config for election " + electionId + " was not found in cache")
		return electionConfig, err
	}

	l.LogInfo("Election config for election " + electionId + " was found in cache")
	err = json.Unmarshal(val, &electionConfig)
	return electionConfig, err
}

func getElectionConfigKey(electionId string) string {
	return "election_config_" + electionId
}

func SaveElectionConfigToCache(electionId string, electionConfig m.ElectionConfig) error {

	l.LogInfo("Saving election config for election " + electionId + " to cache")
	cache := connections.GetInstanceRedisClient()
	ctx := context.Background()
	key := getElectionConfigKey(electionId)
	val, err := json.Marshal(electionConfig)
	if err != nil {
		return err
	}
	timeOut := 30 * time.Second
	val2, errSet := cache.Set(ctx, key, val, timeOut).Result()

	l.LogInfo(val2)
	return errSet
}
