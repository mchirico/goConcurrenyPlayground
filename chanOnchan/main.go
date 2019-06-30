package main

import (
	"log"
	"sync"
	"time"
)

type process struct {
	dur time.Duration
	ch  chan time.Duration
	d   chan time.Time
}

func doStuff(ch <-chan process) {
	log.Printf("We wait to start")
	proc := <-ch
	time.Sleep(proc.dur)
	proc.ch <- proc.dur
	proc.d <- time.Now()
}

func main() {

	sendCh := make(chan process)
	for i := 0; i < 10; i++ {
		go doStuff(sendCh)
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

	time.Sleep(2 * time.Second)
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
