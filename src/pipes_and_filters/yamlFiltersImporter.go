package pipes_and_filters

import (
	"io/ioutil"
	"strconv"

	l "own_logger"

	yaml "gopkg.in/yaml.v2"
)

type FilterWithName struct {
	Name     string
	Function FilterWithParams
}
type FilterWithParams func(data any, params map[string]any) error

type SelectedFilterFromYaml struct {
	Name   string         `yaml:"name"`
	Params map[string]any `yaml:"params"`
}

func (p *Pipeline) LoadFiltersFromYaml(yamlPath string, availableFilters map[string]FilterWithParams) error {

	// Read yaml file
	yamlFile, errReadingFile := ioutil.ReadFile(yamlPath)
	if errReadingFile != nil {
		l.LogError("Error reading yaml file: " + errReadingFile.Error())
		return errReadingFile
	}

	// Parse yaml file
	var selectedFilters []SelectedFilterFromYaml
	errParsingYaml := yaml.Unmarshal(yamlFile, &selectedFilters)

	if errParsingYaml != nil {
		l.LogError("Error parsing yaml file: " + errParsingYaml.Error())
		return errParsingYaml
	}

	// Insert filters in Pipe
	for _, selectedFilter := range selectedFilters {
		filterName, filterExists := availableFilters[selectedFilter.Name]
		if !filterExists {
			l.LogWarning("Filter " + selectedFilter.Name + " not found")
			continue
		}

		maxRetries, specify := selectedFilter.Params["maxRetries"]
		if !specify {
			maxRetries = 1
		} else {
			maxRetries = maxRetries.(int)
		}
		p.Use(insertParameters(filterName, selectedFilter.Params, maxRetries.(int), selectedFilter.Name))
	}
	return nil
}

func insertParameters(missingParameterFilter FilterWithParams, params map[string]any, maxRetries int, filterName string) Filter {
	return func(data any) error {
		var err error
		for i := 0; i < maxRetries; i++ {
			err = missingParameterFilter(data, params)
			if err == nil {
				return nil
			}
			if maxRetries > 1 {
				l.LogWarning("Filter " + filterName + " failed. Retrying for " + strconv.Itoa(i+1) + "° time...")
			}
		}
		return err
	}
}
