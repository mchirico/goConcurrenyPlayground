package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

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

func main() {

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(2*time.Second))
	defer cancel()

	rand := func() interface{} { return rand.Int() }
	repeatFn := RepeatFn()
	take := Take()

	for num := range take(ctx, repeatFn(ctx, rand), 10) {
		fmt.Println(num)
	}

}
