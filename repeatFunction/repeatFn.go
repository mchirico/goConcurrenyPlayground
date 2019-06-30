package repeatFunction

import (
	"context"
	"sync"
)

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

type TakeType func(context.Context, <-chan interface{}, int) <-chan interface{}

//Take calls valueStream num + 1 times, returns num times
func Take() TakeType {

	take := func(
		ctx context.Context,
		valueStream <-chan interface{},
		num int,
	) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-ctx.Done():
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}

	return take
}

type RepeatType func(context.Context, func() interface{}) <-chan interface{}

func RepeatFn() RepeatType {

	repeatFn := func(
		ctx context.Context,
		fn func() interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-ctx.Done():
					return
				case valueStream <- fn():

				}
			}
		}()
		return valueStream
	}

	return repeatFn
}

type RepeatType2 func(context.Context, func(interface{}) interface{}) <-chan interface{}

func RepeatFn2(v interface{}) RepeatType2 {

	repeatFn := func(
		ctx context.Context,
		fn func(interface{}) interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-ctx.Done():
					return
				case valueStream <- fn(v):

				}
			}
		}()
		return valueStream
	}

	return repeatFn
}

type RepeatType3 func(context.Context, func(interface{}) interface{}) <-chan interface{}

func RepeatFn3(v interface{}, num int) RepeatType3 {

	repeatFn := func(
		ctx context.Context,
		fn func(interface{}) interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for i := 0; i < num; i++ {
				select {
				case <-ctx.Done():
					return
				case valueStream <- fn(v):

				}
			}
		}()
		return valueStream
	}

	return repeatFn
}

type FanInType func(context.Context, ...<-chan interface{}) <-chan interface{}

func FanIn() FanInType {
	fanIn := func(
		ctx context.Context,
		channels ...<-chan interface{},
	) <-chan interface{} {
		var wg sync.WaitGroup
		multiplexedStream := make(chan interface{})

		multiplex := func(c <-chan interface{}) { // <3>
			defer wg.Done()
			for i := range c {
				select {
				case <-ctx.Done():
					return
				case multiplexedStream <- i:
				}
			}
		}

		// Select from all the channels
		wg.Add(len(channels)) // <4>
		for _, c := range channels {
			go multiplex(c)
		}

		// Wait for all the reads to complete
		go func() { // <5>
			wg.Wait()
			close(multiplexedStream)
		}()

		return multiplexedStream
	}

	return fanIn
}

func Tee(ctx context.Context, c ...chan interface{}) (<-chan interface{}, <-chan interface{}) {

	a := make(chan interface{})
	b := make(chan interface{})

	go func() {
		defer close(a)
		defer close(b)
		for _, v := range c {
			select {
			case <-ctx.Done():
				return
			case tmp := <-v:
				a <- tmp
				b <- tmp

			}
		}
	}()
	return a, b

}
