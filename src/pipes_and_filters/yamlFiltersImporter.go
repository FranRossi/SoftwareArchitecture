package pipes_and_filters

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type SelectedFilterFromYaml struct {
	Name   string         `yaml:"name"`
	Params map[string]any `yaml:"params"`
}

type FilterFromYaml struct {
	Name     string
	Function FilterWithParams
}
type FilterWithParams func(data any, params map[string]any) error

func (p *Pipeline) LoadFiltersFromYaml(yamlPath string, availableFilters map[string]FilterFromYaml) {

	yamlFile, errReadingFile := ioutil.ReadFile(yamlPath)
	if errReadingFile != nil {
		panic(errReadingFile)
	}

	var selectedFilters []SelectedFilterFromYaml
	errParsingYaml := yaml.Unmarshal(yamlFile, &selectedFilters)

	if errParsingYaml != nil {
		panic(errParsingYaml)
	}

	for _, selectedFilter := range selectedFilters {
		p.Use(insertParamters(availableFilters[selectedFilter.Name].Function, selectedFilter.Params))
	}

}

func insertParamters(missingParamterFilter FilterWithParams, params map[string]any) Filter {
	return func(data any) error {
		return missingParamterFilter(data, params)
	}
}
