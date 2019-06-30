package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type process struct {
	dur time.Duration
	ch  chan time.Duration
	d   chan time.Time
	i   chan int
	s   chan string
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
	if i == 5 {
		proc.s <- "************** 5 ************"
	} else {
		proc.s <- "ok"
	}
	proc.ch <- proc.dur
	proc.d <- time.Now()
	proc.i <- i
}

type accounting struct {
	s   string
	dur time.Duration
}

func Accounting(loops int, pause int) {

	acc := make([]accounting, int(loops))

	sendCh := make(chan process)
	prepCh := make(chan prep)
	for i := 0; i < loops; i++ {
		go doStuff(sendCh, prepCh, i)
	}

	processes := make([]process, loops)
	messages := make(chan string)
	var wg sync.WaitGroup
	for i := 0; i < loops; i++ {
		dur := time.Duration(i+1) * time.Second
		proc := process{dur: dur,
			ch: make(chan time.Duration),
			d:  make(chan time.Time),
			i:  make(chan int),
			s:  make(chan string),
		}
		processes[i] = proc
		sendCh <- proc

		wg.Add(4)
		go func(ch <-chan time.Duration, i int, acc accounting) {
			defer wg.Done()
			dur := <-ch
			acc.dur = dur
			messages <- fmt.Sprintf("%d: slept for %s", i, dur)
		}(processes[i].ch, i, acc[i])
		go func(ch <-chan time.Time, i int, acc accounting) {
			defer wg.Done()
			t := <-ch
			messages <- fmt.Sprintf(" for %s", t)
		}(processes[i].d, i, acc[i])
		go func(ch <-chan int, i int) {
			defer wg.Done()
			idx := <-ch
			messages <- fmt.Sprintf(" idx: %d", idx)
		}(processes[i].i, i)
		go func(ch <-chan string, i int, acc accounting) {
			defer wg.Done()
			s := <-ch
			acc.s = s
			messages <- fmt.Sprintf(" idx: %d, s: %s\n", i, s)
		}(processes[i].s, i, acc[i])

		if i != 0 && (i%pause) == 0 {
			count := 0
		L:
			for {
				select {
				case msg := <-messages:
					fmt.Println("received message", msg)
					count += 1
					if count >= 4*pause {
						time.Sleep(3 * time.Second)
						break L
					}
				default:

					fmt.Println("no message received")
					time.Sleep(2 * time.Second)
				}
			}

			log.Println("******************   CALCULATING ****************")
			for j := 0; j <= i; j++ {
				if acc[i].s != "ok" {
					log.Printf("Got a NOT OK")
					time.Sleep(15 * time.Second)
					break
				}
			}

			time.Sleep(2 * time.Second)
		}
	}
	wg.Wait()

}

// https://play.golang.org/p/8uA60MDFC9E
func main() {
	Accounting(12, 3)
}
