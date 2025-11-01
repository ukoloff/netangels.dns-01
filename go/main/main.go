package main

import (
	"fmt"
	"na01"
)

func main() {
	fmt.Println("Здравствуй, World!")
	s, e := X()
	if e != nil {
		fmt.Println("Oops!")
	}
	fmt.Println(s)
}

func X() (string, error) {
	return na01.Stringify(6 * 7)
}
