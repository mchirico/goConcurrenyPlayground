package repeatFunction

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func TestCreateRepeat(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(2*time.Second))
	defer cancel()

	rand := func() interface{} { return rand.Int() }
	repeatFn := RepeatFn()
	take := Take()

	count := 0

	for num := range take(ctx, repeatFn(ctx, rand), 10) {
		fmt.Println(num)
		count += 1
	}

	if count != 10 {
		t.FailNow()
	}

}

func TestCreateRepeat2(t *testing.T) {

	t.Logf("Number of CPU's: %d\n", runtime.NumCPU())

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(2*time.Second))
	defer cancel()

	count := Counter{}

	v := func(a interface{}) interface{} {
		count.Inc()
		return a.(string) +
			fmt.Sprintf("%d", count.Value())
	}

	repeatFn2 := RepeatFn2("value: ")
	take := Take()

	for num := range take(ctx, repeatFn2(ctx, v), 10) {
		fmt.Println(num)

	}

	// Note, take will call function 11 times
	if count.Value() != 11 {
		t.Logf("count: %d\n", count.Value())
		t.FailNow()
	}

}

func TestFanIn(t *testing.T) {

	//numFinders := runtime.NumCPU()
	numFinders := 10
	fmt.Printf("Spinning up %d prime finders.\n", numFinders)
	finders := make([]<-chan interface{}, numFinders)

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(2*time.Second))
	defer cancel()

	count := Counter{}

	v := func(a interface{}) interface{} {
		count.Inc()

		return a.(string) +
			fmt.Sprintf("%d", count.Value())
	}

	take := Take()
	fanIn := FanIn()

	for i := 0; i < numFinders; i++ {
		repeatFn2 := RepeatFn2(fmt.Sprintf("func(%d) ", i))
		finders[i] = repeatFn2(ctx, v)
	}

	take_num := 5

	for num := range take(ctx, fanIn(ctx, finders...), take_num) {
		fmt.Printf("\t%s\n", num)
	}

	//
	if count.Value() != (numFinders*2 + take_num) {
		t.Logf("count: %d\n", count.Value())
		t.FailNow()
	}
	t.Logf("count: %d\n", count.Value())
	t.Logf("numFinders*2+take_num : %d\n", numFinders*2+take_num)

}

func TestCreateRepeat3(t *testing.T) {

	t.Logf("Number of CPU's: %d\n", runtime.NumCPU())

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(2*time.Second))
	defer cancel()

	count := Counter{}

	v := func(a interface{}) interface{} {
		count.Inc()
		return a.(string) +
			fmt.Sprintf("%d", count.Value())
	}

	repeatFn3 := RepeatFn3("value: ", 10)
	take := Take()

	for num := range take(ctx, repeatFn3(ctx, v), 10) {
		fmt.Println(num)

	}

	// Note, take will call function 10 times
	if count.Value() != 10 {
		t.Logf("count: %d\n", count.Value())
		t.FailNow()
	}

}

func TestFanInLimit(t *testing.T) {

	//numFinders := runtime.NumCPU()
	numFinders := 8
	num := 2
	fmt.Printf("Spinning up %d finders.\n", numFinders)
	finders := make([]<-chan interface{}, numFinders)

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(2*time.Second))
	defer cancel()

	count := Counter{}

	v := func(a interface{}) interface{} {
		count.Inc()

		return a.(string) +
			fmt.Sprintf("%d", count.Value())
	}

	take := Take()
	fanIn := FanIn()

	for i := 0; i < numFinders; i++ {
		repeatFn3 := RepeatFn3(fmt.Sprintf("func(%d) ", i), num)
		finders[i] = repeatFn3(ctx, v)
	}

	take_num := numFinders * num

	for num := range take(ctx, fanIn(ctx, finders...), take_num) {
		fmt.Printf("\t\t%s\n", num)
	}

	//
	if count.Value() != (numFinders * num) {
		t.Logf("count: %d\n", count.Value())
		t.FailNow()
	}
	t.Logf("count: %d\n", count.Value())
	t.Logf("numFinders*2+take_num : %d\n", numFinders*num)

}

func TestTee(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(2*time.Second))
	defer cancel()

	c0 := make(chan interface{})
	c1 := make(chan interface{})

	go func() {
		c0 <- "test 1"
		c1 <- "test 2"
	}()

	a, b := Tee(ctx, c0, c1)

	if "test 1" == <-a && "test 1" == <-b {
		if "test 2" == <-a && "test 2" == <-b {
		} else {
			t.FailNow()
		}
	}

}
