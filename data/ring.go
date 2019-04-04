package data

import (
	"sync/atomic"
)

/*
Ring Please meet ring
*/
type Ring struct {
	obj   []interface{}
	ridx  int
	widx  int
	delay int32
	run   bool
	ch    chan string
}

/*
NewRing creates a new ring on length length and with the given producer
*/
func NewRing(length int, prod func() interface{}) *Ring {
	r := Ring{make([]interface{}, length), 0, 0, int32(length), true, make(chan string)}
	go func() {
		for r.run {
			if r.delay > 0 {
				// overwriting the current "cell" with the produced object
				r.obj[r.widx] = prod()
				// pointing at the next "cell"
				r.widx = (r.widx + 1) % length
				// decreasing the delay : The producer is "catching up"
				// Should not be a problem ; _I_ am the only one decreasing this
				// so condition of the if cannot be violated.
				atomic.AddInt32(&r.delay, -1)
			} else {
				// the delay can only be 0 at this point.
				// When the channel is writen to, at least one call to the Next method has increased it
				// so we can start looping again.
				<-r.ch
			}
		}
	}()
	return &r
}

/*
Stop stops producing events in the ring
*/
func (r *Ring) Stop() {
	r.run = false
	// wake up the writer, just in case
	select {
	case r.ch <- "stop":
	default:
	}
}

/*
Length of this buffer.
*/
func (r *Ring) Length() int {
	return len(r.obj)
}

/*
Next gives the next element in the ring
*/
func (r *Ring) Next() interface{} {
	res := r.obj[r.ridx]
	r.ridx = (r.ridx + 1) % r.Length()
	// incrementing the delay
	if 0 == r.delay {
		var read = r.delay
		for !atomic.CompareAndSwapInt32(&r.delay, read, (read+1)%int32(r.Length())) {
			read = r.delay
		}
		select {
		case r.ch <- "GoOn":
		default:
		}
	}

	return res
}
