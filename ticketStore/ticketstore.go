package ticketStore

import (
	"runtime"
	"sync/atomic"
)

type TicketStore struct {
	ticket *uint64
	done   *uint64
	slots  []string
}

func NewTicketStore(n int) TicketStore {

	t := TicketStore{}
	t.ticket = new(uint64)
	t.done = new(uint64)
	t.slots = make([]string, n, n)
	return t

}

func (ts *TicketStore) Put(s string) {
	t := atomic.AddUint64(ts.ticket, 1) - 1
	ts.slots[t] = s
	for !atomic.CompareAndSwapUint64(ts.done, t, t+1) {
		runtime.Gosched()
	}
}

func (ts *TicketStore) GetDone() []string {
	return ts.slots[:atomic.LoadUint64(ts.done)]
}
