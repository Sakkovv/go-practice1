package main

import (
	"fmt"
)

func main() {
	var a, b float64
	var operation string

	fmt.Scan(&a, &operation, &b)

	switch operation {
	case "+":
		fmt.Println(a + b)
	case "-":
		fmt.Println(a - b)
	case "*":
		fmt.Println(a * b)
	case "/":
		if b != 0 {
			fmt.Println(a / b)
		} else {
			fmt.Println("error")
		}
	default:
		fmt.Println("error")
	}
}
