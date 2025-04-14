package main

import (
	"fmt"
)

func main() {
	n, err := fmt.Println("Hello, World!")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Number of bytes written:", n)
	}
}
