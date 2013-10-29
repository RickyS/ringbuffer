// package ringbuffer implements a sequential compact FIFO + LILO
// To Use:  type RingElement thing
//          var myThing RingElement
//          rb := RingBuffer.New(40)
//          rb.Write(myThing) // Et cetera
//          aThing := rb.Read()
//          for 0 < rb.Size() {
//              doSomethingWith(rb.Read())
//          }
package ringbuffer

import (
// "fmt"
// "os"
)

type RingElement interface{}

type RingBuffer struct {
	data    []RingElement
	in, out int // place of next in (Write) and next out (Read).  These are subscripts.
	size    int // Number of items currenly in the ring buffer.
}

type RingBufferError struct {
	What string
}

func (e *RingBufferError) Error() string {
	return e.What
}

// Inspect the internal state of the ring buffer and complain if not ok.
var invNum int

func (b *RingBuffer) invariant() bool {
	capacity := cap(b.data)
	invNum = 0
	ok := (0 <= b.in) && (b.in < capacity) &&
		(0 <= b.out) && (b.out < capacity) &&
		(0 <= b.size) && (b.size <= capacity) && // size can equal capacity.  Subscripts cannot.
		(capacity == len(b.data))

	if !ok {
		invNum = 1
	} else {
		if b.out < b.in {
			ok = b.size == b.in-b.out
			if !ok {
				invNum = 2
			}
		} else if b.in < b.out {
			ok = b.size == (capacity-b.out)+b.in
			if !ok {
				invNum = 3
			}
		} else { //  in == out
			ok = (0 == b.size) || (capacity == b.size)
			if !ok {
				invNum = 4
			}
		}
	}

	if !ok {
		b.internalDump("invariant")
	}
	return ok
}

// New allocates and initializes a new ring buffer of specified size
func New(n int) *RingBuffer {
	b := &RingBuffer{data: make([]RingElement, n), // Contents
		in: 0, out: 0, size: 0}
	b.invariant()
	return b
}

func (b *RingBuffer) next(subscript int) int {
	subscript++
	if subscript >= cap(b.data) {
		subscript = 0
	}
	return subscript
}

// Write inserts an element into the ring buffer.
func (b *RingBuffer) Write(datum RingElement) error {
	if b.size >= cap(b.data) {
		return &RingBufferError{"RingBuffer is full\n"}
	}

	b.data[b.in] = datum
	b.in = b.next(b.in)
	b.size++
	b.invariant()

	return nil
}

// Read fetches an element from the ring buffer.
func (b *RingBuffer) Read() RingElement {
	if 0 >= b.size {
		b.invariant()
		return 0
		//return &RingBufferError{"RingBuffer is empty\n"}
	}
	b.size--
	tmp := b.data[b.out]
	b.out = b.next(b.out)
	return tmp
}

// Number of slots currently in use.  Total writes - Total reads.
func (b *RingBuffer) Leng() int {
	b.invariant()
	return b.size
}

// Is the buffer currently full?
func (b *RingBuffer) Full() bool {
	b.invariant()
	return b.size >= cap(b.data)
}

// Obliterate, Purge, and Remove the contents of the ring buffer.
func (b *RingBuffer) Clear() {
	b.in, b.out, b.size = 0, 0, 0
	for i := 0; i < len(b.data); i++ {
		b.data[i] = 0
	}
	b.invariant()
	// free(b.data)  // unimplemented.
}
