package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	l "pipes_and_filters"
)

func main() {

	// Insert all the filters you may want to use:
	availableFilters := make(map[string]l.FilterFromYaml)

	FilterEcho := l.FilterFromYaml{}
	FilterEcho.Function = FilterEchoData

	FilterLowerAge := l.FilterFromYaml{}
	FilterLowerAge.Function = FilterCheckAge

	FilterUpperAge := l.FilterFromYaml{}
	FilterUpperAge.Function = FilterCheckAgeUpper

	availableFilters["echo"] = FilterEcho
	availableFilters["age_lower"] = FilterLowerAge
	availableFilters["age_upper"] = FilterUpperAge

	fmt.Println("Test")

	p := l.Pipeline{}

	p.LoadFiltersFromYaml("test.yaml", availableFilters)

	validateNumber(70, &p)
	fmt.Println("")
	validateNumber(3, &p)
	fmt.Println("")
	validateNumber(50, &p)

	println("Number validated - End")

	fmt.Print("Press Enter to exit")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

}

func validateNumber(num int, p *l.Pipeline) {
	errors := p.Run(num)

	if len(errors) == 0 {
		println("All OK")
	} else {
		println("Errors:")
		for _, err := range errors {
			println(err.Error())
		}
	}
}

// FILTERS:

// Params:
// min_age: int
func FilterCheckAge(data any, params map[string]any) error {

	minAge, validParse := params["min_age"].(int)
	if !validParse {
		return errors.New("min_age is not an int")
	}

	println("Checking age Lower")
	if input, ok := data.(int); ok {
		if input < minAge {
			return errors.New("underage")
		}
		return nil
	} else {
		return errors.New("Invalid data type")
	}
}

// Params:
// max_age: int
func FilterCheckAgeUpper(data any, params map[string]any) error {

	maxAge, validParse := params["max_age"].(int)
	if !validParse {
		return errors.New("max_age is not an int")
	}

	println("Checking age Upper")
	if input, ok := data.(int); ok {
		if input > maxAge {
			return errors.New("Too old")
		}
		return nil
	} else {
		return errors.New("Invalid data type")
	}
}

func FilterCheckAgeBetween(data any, params map[string]any) error {

	minAge, validParse := params["min_age"].(int)
	if !validParse {
		return errors.New("min_age is not an int")
	}

	maxAge, validParse := params["max_age"].(int)
	if !validParse {
		return errors.New("max_age is not an int")
	}

	println("Checking age between")
	if input, ok := data.(int); ok {
		if input < minAge && input > maxAge {
			return errors.New("Either too young or too old")
		}
		return nil
	} else {
		return errors.New("Invalid data type")
	}
}

func FilterEchoData(data any, params map[string]any) error {
	input, _ := data.(int)
	fmt.Printf("data Data: %d\n", input)
	return nil
}
