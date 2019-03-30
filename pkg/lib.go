package pkg

import (
	"context"
	"fmt"
	"sync"
)

func process(s string) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func(s string) {
		defer wg.Done()
		fmt.Printf("%s\n", s)
	}(s)
	go func(s string) {
		defer wg.Done()
		fmt.Printf("%q\n", s)
	}(s)
	wg.Wait()
}

func MyWait() {

	for i := 0; i < 3; i++ {
		s := fmt.Sprintf("value: %d", i)
		process(s)
	}

}

func Or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			case <-Or(append(channels[3:], orDone)...):
			}
		}
	}()
	return orDone
}

type BB func(context.Context, ...interface{}) <-chan interface{}

func Repeat() BB {

	repeat := func(
		ctx context.Context,
		values ...interface{},
	) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-ctx.Done():
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}

	return repeat
}

type TakeType func(context.Context, <-chan interface{}, int) <-chan interface{}

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

func RepeadFN() RepeatType {

	repeatFN := func(
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

	return repeatFN
}
