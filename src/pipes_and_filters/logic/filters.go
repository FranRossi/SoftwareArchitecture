package logic

import (
	"errors"
	"fmt"
)

func FilterCheckAge(input Any, params map[string]string) error {

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

func FilterCheckAgeUpper(input Any, params map[string]string) error {

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

func FilterEchoInput(input Any, params map[string]string) error {
	data, _ := input.(int)
	fmt.Printf("Input Data: %d\n", data)
	return nil
}
