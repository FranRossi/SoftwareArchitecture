package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func GetElection() {
	/*
		election, err := getElectionFromCache

		if err == nil {
			return election
		}

		election, err = getElectionFromDB
		if err == nil {
			storeElectionInCache()
			return election
		}

	*/
}

type VoterModel struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
}

func SaveVoterToCache(voter VoterModel) error {
	cache := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81"})
	ctx := context.Background()

	json, _ := json.Marshal(voter)
	fmt.Println(voter)
	cacheErr := cache.Set(ctx, voter.ID, json, 20*time.Second).Err()
	if cacheErr != nil {
		return cacheErr
	}
	return nil
}

func GetVoter(VoterId string) {
	cache := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81", DB: 0})
	ctx := context.Background()

	val, err := cache.Get(ctx, VoterId).Bytes()
	if err != nil {
		// return GetVoterFromDB(VoterId)
		fmt.Println(err)
	}

	var voter VoterModel
	errParsing := json.Unmarshal(val, &voter)
	if errParsing != nil {
		fmt.Println(errParsing)
	}
	fmt.Println(voter)

}
