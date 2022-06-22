package controllers

import (
	"encoding/json"
	"fmt"
	"stats_service/logic"
	"stats_service/models"

	mq "message_queue"
	l "own_logger"
	p_f "pipes_and_filters"
)

func ListenForNewStats() {

	go listenForStats(logic.AddVoteToCertainGroupActual, "stats-actual")
	go listenForStats(logic.AddVoteToCertainGroupTotal, "stats-total")

}

func listenForStats(filter func(data any, params map[string]any) error, queueName string) {
	pipeLine := p_f.Pipeline{}
	availableFilters := map[string]p_f.FilterWithParams{
		"add_to_vote_group": filter,
	}

	errLoadingYaml := pipeLine.LoadFiltersFromYaml("statisticsGroups.yaml", availableFilters)
	if errLoadingYaml != nil {
		return
	}

	go listen(queueName, pipeLine)

}

func listen(queueName string, pipeLine p_f.Pipeline) {
	mq.GetMQWorker().Listen(50, queueName, func(message []byte) error {
		var stats models.VoterStats
		err := json.Unmarshal(message, &stats)
		if err != nil {
			l.LogError("Couldn't parse message")
			return err
		}
		fmt.Print(".") // TODO optional, use for demo
		go runPipeAndLog(pipeLine, stats)
		return nil
	})
}

func runPipeAndLog(pipeLine p_f.Pipeline, stats models.VoterStats) {
	filtersErrs := pipeLine.Run(stats)
	logFiltersErrors(filtersErrs)
}

func logFiltersErrors(errors []error) {
	for _, er := range errors {
		l.LogError(er.Error())
	}
}
