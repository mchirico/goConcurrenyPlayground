package repeatFunction

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
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

type Counter struct {
	sync.Mutex
	value int
}

func (c *Counter) Inc() {
	c.Lock()
	defer c.Unlock()
	c.value += 1
}

func (c *Counter) Value() int {
	c.Lock()
	defer c.Unlock()
	return c.value
}

func TestCreateRepeat2(t *testing.T) {

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

	// Note, get's call 11 times
	if count.Value() != 11 {
		t.Logf("count: %d\n", count.Value())
		t.FailNow()
	}

}
