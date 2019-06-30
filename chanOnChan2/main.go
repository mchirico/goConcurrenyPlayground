package main

import (
	"context"
	"log"
	"sync"
	"time"
)

type process struct {
	dur time.Duration
	ch  chan time.Duration
	d   chan time.Time
}

type FN func(context.Context, func() interface{}) <-chan interface{}

type prep struct {
	f     FN
	start chan time.Time
	end   chan time.Time
	ci    chan interface{}
}

func doStuff(ch <-chan process, chP <-chan prep, i int) {
	log.Printf("We wait to start:%d\n", i)
	proc := <-ch
	log.Printf("processing:%d\n", i)
	time.Sleep(proc.dur)
	proc.ch <- proc.dur
	proc.d <- time.Now()
}

func main() {

	sendCh := make(chan process)
	prepCh := make(chan prep)
	for i := 0; i < 10; i++ {
		go doStuff(sendCh, prepCh, i)
	}

	processes := make([]process, 10)
	for i := 0; i < 10; i++ {
		dur := time.Duration(i+1) * time.Second
		proc := process{dur: dur,
			ch: make(chan time.Duration),
			d:  make(chan time.Time),
		}
		processes[i] = proc
		sendCh <- proc
	}

	//time.Sleep(2 * time.Second)
	// recieve on each struct's ack channel
	var wg sync.WaitGroup // use this to block until all goroutines have received the ack and logged
	for i := 0; i < 10; i++ {
		wg.Add(2)
		go func(ch <-chan time.Duration) {
			defer wg.Done()
			dur := <-ch
			log.Printf("slept for %s", dur)
		}(processes[i].ch)
		go func(ch <-chan time.Time) {
			defer wg.Done()
			dur := <-ch
			log.Printf(" for %s", dur)
		}(processes[i].d)
	}
	wg.Wait()
}
