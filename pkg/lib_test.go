package pkg

import (
	"fmt"
	"testing"
	"time"
)

func TestOr(t *testing.T) {

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-Or(
		sig(2*time.Second),
		sig(13*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
	if time.Since(start).Seconds() > 3 {
		t.Fatalf("Took too long. Expected: 2, got %f\n", time.Since(start).Seconds())
	}

}

func TestMyWait(t *testing.T) {
	MyWait()
}
