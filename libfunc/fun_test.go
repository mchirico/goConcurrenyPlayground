package libfunc

import (
	"fmt"
	"testing"
)

func TestStages(t *testing.T) {

	multiply := func(values []int, multiplier int) []int {
		multipliedValues := make([]int, len(values))
		for i, v := range values {
			multipliedValues[i] = v * multiplier
		}
		return multipliedValues
	}

	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, v := range values {
			addedValues[i] = v + additive
		}
		return addedValues
	}

	ints := []int{1, 2, 3, 4}

	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}

}

func TestSelect(t *testing.T) {

	messages := make(chan string)
	s := make(chan string)

	go func() { messages <- "ping" }()
	go func() { s <- "test" }()

	for i := 0; i < 2; i++ {
		select {
		case msg := <-messages:
			fmt.Println(msg)
		case s0 := <-s:
			fmt.Println(s0)

		}
	}

}

func TestSelect2(t *testing.T) {

	c1 := make(chan interface{})

	go func() { c1 <- struct{}{} }()

	select {
	case <-c1:
		fmt.Println("here")
	}

}
