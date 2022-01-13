package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type calc struct {
}

func (calc) operate(input string, operator string) int {
	cleanInput := strings.Split(input, operator)
	operator1 := parse(cleanInput[0])
	operator2 := parse(cleanInput[1])
	switch operator {
	case "+":
		fmt.Println(operator1 + operator2)
		return operator1 + operator2
	case "-":
		fmt.Println(operator1 - operator2)
		return operator1 - operator2
	case "*":
		fmt.Println(operator1 * operator2)
		return operator1 * operator2
	case "/":
		fmt.Println(operator1 / operator2)
		return operator1 / operator2
	default:
		fmt.Println(operator, "is not valid.")
		return 0
	}
}

func parse(input string) int {
	operator, _ := strconv.Atoi(input)
	return operator
}

func readInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func main() {
	input := readInput()
	operator := readInput()
	c := calc{}
	fmt.Println(c.operate(input, operator))

}
