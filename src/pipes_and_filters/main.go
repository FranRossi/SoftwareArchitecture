package main

import (
	"errors"
	"fmt"
)

// TODO revisar si esta linea esta bien
type any = interface{}

type Filter func(any) error

type Pipeline struct {
	filters []Filter
}

func (p *Pipeline) Use(f ...Filter) {
	p.filters = append(p.filters, f...)
}

func (p Pipeline) Run(input any) []error {

	//	var wg sync.WaitGroup
	out := make(chan error, len(p.filters))

	// Runs each filters and saves the error to the out channel
	for _, f := range p.filters {
		out <- f(input)
		//wg.Done()
		// go func() {
		// 	out <- f(input)
		// }()
	}

	// Waits for all the filters to finish
	//	wg.Add(len(p.filters))
	//	wg.Wait()
	close(out)
	var errors []error

	// Saves errors of all the filters
	for err := range out {
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func FilterCheckAge(input any) error {

	println("Checking age Lower")
	if data, ok := input.(int); ok {
		if data < 18 {
			return errors.New("underage")
		}
		return nil
	} else {
		return errors.New("Invalid data type")
	}
}

func FilterCheckAgeUpper(input any) error {

	println("Checking age Upper")
	if data, ok := input.(int); ok {
		if data > 60 {
			return errors.New("Too old")
		}
		return nil
	} else {
		return errors.New("Invalid data type")
	}
}

func FilterEchoInput(input any) error {
	data, _ := input.(int)
	fmt.Println("")
	fmt.Printf("Input Data: %d\n", data)
	return nil
}

func main() {
	p := Pipeline{}
	p.Use(FilterEchoInput, FilterCheckAgeUpper, FilterCheckAge)

	validateNumber(70, &p)
	validateNumber(3, &p)
	validateNumber(50, &p)

	println("Number validated - End")
}

func validateNumber(num int, p *Pipeline) {
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
