package logic

import (
	"fmt"
	"own_logger"
	p_f "pipes_and_filters"
	"stats_service/models"
	domain "stats_service/models"
	"stats_service/repository"
	"time"
)

func ListenForNewStats(){

	pipeLineTotal :=
	pipeLineActual := 

}

func createTotalPipeLine(){
	p := p_f.Pipeline{}

}


func RegisterVoteOnCertainGroupActual(voter *models.VoterStats) {
	p := p_f.Pipeline{}
	statistics := transformDataToStatistics(voter)
	statistics.ElectionId = voter.ElectionId

	availableFilters := map[string]p_f.FilterWithParams{
		"add_to_vote_group": AddVoteToCertainGroup,
	}
	errLoadingYaml := p.LoadFiltersFromYaml("statisticsGroups.yaml", availableFilters)
	if errLoadingYaml != nil {
		return
	}
	errors := p.Run(statistics)
	if len(errors) > 0 {
		// TODO add error code
		LogValidationErrors(errors)
		fmt.("statistics are not valid: %v", errors)
	}
}

func LogValidationErrors(errors []error) {
	for _, er := range errors {
		own_logger.LogError(er.Error())
	}
}

func transformDataToStatistics(voter *models.VoterModel) domain.Statistics {
	age := getAge(voter.BirthDate)
	statistics := domain.Statistics{
		Age:     age,
		Region:  voter.Region,
		Circuit: voter.OtherFields["circuit"].(string),
		Sex:     voter.Sex,
	}
	return statistics
}

func getAge(birthDate string) int {
	t, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		return 0
	}
	return time.Now().Year() - t.Year()
}

func AddVoteToCertainGroupTotal(data any, params map[string]any) error {

}

func AddVoteToCertainGroup(data any, params map[string]any) error {
	statistics := data.(domain.Statistics)
	minAge := params["min_age"].(int)
	maxAge := params["max_age"].(int)
	sex := params["sex"].(string)
	groupType := params["type"].(string)
	groupName := params["name"].(string)
	if statistics.Age >= minAge && statistics.Age <= maxAge && statistics.Sex == sex {
		err := repository.UpdateStatistics(statistics, groupType, groupName)
		if err != nil {
			return fmt.Errorf("error storing statistics on database")
		}
	}
	return nil
}
