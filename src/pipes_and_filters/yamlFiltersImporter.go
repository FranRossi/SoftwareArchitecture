package pipes_and_filters

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type SelectedFilterFromYaml struct {
	Name   string         `yaml:"name"`
	Params map[string]any `yaml:"params"`
}

type FilterFunctionWithName struct {
	Name     string
	Function FilterWithParams
}
type FilterWithParams func(data any, params map[string]any) error

func (p *Pipeline) LoadFiltersFromYaml(yamlPath string, availableFilters []FilterFunctionWithName) {

	// Array to map
	filtersMap := make(map[string]FilterWithParams)
	for _, filter := range availableFilters {
		filtersMap[filter.Name] = filter.Function
	}

	// Read yaml file
	yamlFile, errReadingFile := ioutil.ReadFile(yamlPath)
	if errReadingFile != nil {
		panic(errReadingFile)
	}

	// Parse yaml file
	var selectedFilters []SelectedFilterFromYaml
	errParsingYaml := yaml.Unmarshal(yamlFile, &selectedFilters)

	if errParsingYaml != nil {
		panic(errParsingYaml)
	}

	// Insert filters in Pipe
	for _, selectedFilter := range selectedFilters {
		filterName, filterExists := filtersMap[selectedFilter.Name]
		if !filterExists {
			panic("Filter " + selectedFilter.Name + " not found")
		}
		p.Use(insertParamters(filterName, selectedFilter.Params))
	}
}

func insertParamters(missingParameterFilter FilterWithParams, params map[string]any) Filter{
	return func(data any) error {
		return missingParameterFilter(data, params)
	}
}
