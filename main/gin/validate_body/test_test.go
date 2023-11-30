package main

import (
	"fmt"
	"testing"
)

type Person struct {
	name string
	age  int
}

func TestName(t *testing.T) {
	p1 := Person{name: "John", age: 30}

	// Copy the person by value
	p2 := p1

	// Change the age of p2
	p2.age = 35

	// Print the values of both persons
	fmt.Println(p1) // Output: {John 30}
	fmt.Println(p2) // Output: {John 35}
}
