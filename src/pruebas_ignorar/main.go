package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {

	castedElement := true
	// print num
	fmt.Println(fmt.Sprint(castedElement))
	// conver string to int
	fmt.Print(strconv.ParseBool(fmt.Sprint(castedElement)))
	// var varPrueba any
	// varPrueba = 'd'
	// fmt.Println(varPrueba)
	// bt := []byte(string(varPrueba))

	// switch v := varPrueba.(type) {
	// case int:
	// 	fmt.Printf("Twice %v is %v\n", v, v*2)
	// case string:
	// 	fmt.Printf("%q is %v bytes long\n", v, len(v))
	// default:
	// 	fmt.Printf("I don't know about type %T!\n", v)
	// }
}

//////////////////////////////////

func pruebaFunciones() {
	multiplicarPor2 := fijarParametro(multiplicar, 2)
	result := multiplicarPor2(34)
	fmt.Println(result)
}

func fijarParametro(function func(a int, b int) int, x int) func(b int) int {
	return func(b int) int {
		return function(x, b)
	}
}

func multiplicar(a int, b int) int {
	return a * b
}

func pruebaTime() {
	startDate := time.Now()
	// add 10 seconds to current time
	startDate = startDate.Add(6 * time.Second)

	timer := time.NewTimer(startDate.Sub(time.Now()))
	//timer := time.NewTimer(10 * time.Second)

	go func() {
		defer timeTrack(time.Now(), "comienzo de elección")
		<-timer.C
		fmt.Println("Comenzó la elección")
	}()

	// read character from stdin

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("-> ")
	text, _ := reader.ReadString('\n')
	print(text)

	// var input string
	// done := make(chan bool)
	// print(input)
	// if input == "exit" {
	// 	done <- true
	// }
	// <-done
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
