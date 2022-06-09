package logic

type SelectedFilter struct {
	Name   string
	Params map[string]string
}

func (p *Pipeline) LoadFiltersFromYaml(yaml string) {

	var availableFilters = make(map[string]Filter)

	availableFilters["echo"] = FilterEchoInput
	availableFilters["age_lower"] = FilterCheckAge
	availableFilters["age_upper"] = FilterCheckAgeUpper

	p.Use(FilterEchoInput)

}
