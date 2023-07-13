package main

import (
	"fmt"
)

func AddOne[Any ~string | ~int](t, a Any) Any {
	return t+a
}

type Adder interface {
	~string | ~int
}

func AddTwo[T Adder](t, a T) T {
	return t+a
}

func main() {
	// result := AddOne("hello ","world")
	result := AddTwo("hello ","world")
	fmt.Println(result)
}