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
		return electionConfig, err
	}

	err = json.Unmarshal(val, &electionConfig)
	return electionConfig, err
}

func getElectionConfigKey(electionId string) string {
	return "election_config_" + electionId
}

func SaveElectionConfigToCache(electionConfig m.ElectionConfig) error {

	l.LogInfo("Saving election config with id" + electionConfig.ElectionId + "to cache")
	cache := connections.GetInstanceRedisClient()
	ctx := context.Background()
	key := getElectionConfigKey(electionConfig.ElectionId)
	val, err := json.Marshal(electionConfig)
	if err != nil {
		return err
	}
	timeOut := 30 * time.Second
	_, err = cache.Set(ctx, key, val, timeOut).Result()
	return err
}
