package main

import (
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(2)
		s := fmt.Sprintf("value: %d", i)
		go func(s string) {
			defer wg.Done()
			fmt.Printf("%s\n", s)
		}(s)
		go func(s string) {
			defer wg.Done()
			fmt.Printf("%q\n", s)
		}(s)
	}

	wg.Wait()
}
