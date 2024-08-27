package main

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

var NAN_ERROR = "This is not a number"

var x int
var y int
var history []string = []string{}

func main() {
	var keep_going bool = true
	var input string

	fmt.Println("1) Press Enter to sum two numbers")
	fmt.Println("2) Type 'h' to show the history calculated")
	fmt.Println("3) Type anything else to stop")

	for keep_going {
		input = ""
		fmt.Print("\n op> ")
		fmt.Scanf("%s\n", &input)
		fmt.Println()
		keep_going = handleOperation(input)
	}

	fmt.Println("> Final history:")
	printHistory()
}

func inputToNumber(input string) (int, *string) {
	input = strings.Trim(input, " ")
	result, err := strconv.Atoi(input)
	if err != nil {
		return math.MaxInt64, &NAN_ERROR
	}

	return result, nil
}

func handleOperation(input string) bool {
	input = strings.Trim(input, " ")
	if input == "" {
		var err *string
		fmt.Print("First number: ")
		fmt.Scanf("%s\n", &input)
		x, err = inputToNumber(input)
		if err != nil {
			fmt.Println("> ", *err)
			return true
		}

		fmt.Print("Second number: ")
		fmt.Scanf("%s\n", &input)
		y, err = inputToNumber(input)
		if err != nil {
			fmt.Println("> ", *err)
			return true
		}

		tot := x + y
		fmt.Println("> ", tot)
		history = slices.Insert(history, 0, fmt.Sprintf("%d + %d = %d", x, y, tot))
		return true
	}

	if input == "h" {
		fmt.Println("> History:")
		printHistory()
		return true
	}

	return false
}

func printHistory() {
	for _, curr := range history {
		fmt.Println("> ", curr)
	}
}
