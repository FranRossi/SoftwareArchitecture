package logic

import (
	l "own_logger"
	"stats_service/models"
	"stats_service/repository"
	"time"
)

func AddVoteToCertainGroupGenerics(data any, params map[string]any, field string,
	updateStatsFunc func(string, models.VoterStats, string, string, int, int) error) error {
	statistics := data.(models.VoterStats)
	minAge := params["min_age"].(int)
	maxAge := params["max_age"].(int)
	sex := params["sex"].(string)
	groupType := params["type"].(string)
	groupName := params["name"].(string)
	statistics.Age = getAge(statistics.BirthDate)

	if statistics.Age >= minAge && statistics.Age <= maxAge && statistics.Sex == sex {
		err := updateStatsFunc(field, statistics, groupType, groupName, minAge, maxAge)
		if err != nil {
			l.LogError("error storing statistics on database")
			return err
		}
	}
	return nil
}

func AddVoteToCertainGroupTotal(data any, params map[string]any) error {
	return AddVoteToCertainGroupGenerics(data, params, "capacity", repository.UpdateStatistics)
}

func AddVoteToCertainGroupActual(data any, params map[string]any) error {
	return AddVoteToCertainGroupGenerics(data, params, "votes", repository.UpdateStatistics)
}

func getAge(birthDate string) int {
	t, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return 0
	}
	return time.Now().Year() - t.Year()
}
