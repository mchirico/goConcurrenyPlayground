package pingpong

import "fmt"

func ping(pings chan<- string, msg string) {
	pings <- msg
}

func t(pings <-chan string, nxt string) chan string {
	pongs := make(chan string, 1)
	msg := <-pings
	msg += " " + nxt
	pongs <- msg
	return pongs

}

func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}

func RunStuff() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")

	pong(t(t(pings, "stuff"), "more"), pongs)
	fmt.Println(<-pongs)
}
