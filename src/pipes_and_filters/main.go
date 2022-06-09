package main

import (
	"fmt"
	l "pipes_and_filters/logic"
)

func main() {
	p := l.Pipeline{}
	p.Use(l.FilterEchoInput, l.FilterCheckAgeUpper, l.FilterCheckAge)

	validateNumber(70, &p)
	fmt.Println("")
	validateNumber(3, &p)
	fmt.Println("")
	validateNumber(50, &p)

	println("Number validated - End")
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
