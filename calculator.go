package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	operation := scanner.Text()
	fmt.Println(operation)
	values := strings.Split(operation, "+")
	fmt.Println(values)
	fmt.Println(values[0] + values[1])
	operator1, _ := strconv.Atoi(values[0])
	operator2, _ := strconv.Atoi(values[1])
	fmt.Println(operator1 + operator2)
}
