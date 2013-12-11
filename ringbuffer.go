// package ringbuffer implements a sequential compact FIFO and LILO. Also called a Queue.
// To Use:
//          type myThing ringbuffer.RingElement
//          var whatever == myThing("whatever") // Assuming a conversion from string.
//          rb := RingBuffer.New(40)
//          rb.Write(myThing) // Et cetera
//          aThing := rb.Read()
//          for 0 < rb.Size() {
//              doSomethingWith(rb.Read())
//          }
//
//  THIS IS NOT CONCURRENT —— ONE GOROUTINE ONLY.
package ringbuffer

// A ring buffer is stored in an array of ringbuffer.RingElement, of the size requested.
type RingElement interface{}

type RingBuffer struct {
	data    []RingElement
	in, out int // Place of next in (Write). Place of next out (Read).  These are subscripts.
	size    int // Number of items currenly in the ring buffer.
}

type RingBufferError struct {
	What string
}

// "Convert" ringbuffer.RingBufferError into a string.
func (e *RingBufferError) Error() string {
	return e.What
}

///// Inspect the internal state of the ring buffer and complain if not ok. ////
var invNum int // invNum is an error code.

// The conditions checked here can best be understood by drawing the obvious diagram of the array.
func (b *RingBuffer) invariant() bool { // You can remove this function and all ref to it.

	if (nil == b.data) && (0 == b.in) && (0 == b.out) && (0 == b.size) {
		// The RingBuffer has been nilled out by calling Clear()
		return true // All is good.
	}

	capacity := cap(b.data)
	invNum = 0
	ok := (0 <= b.in) && (b.in < capacity) &&
		(0 <= b.out) && (b.out < capacity) &&
		(0 <= b.size) && (b.size <= capacity) && // size can equal capacity.  Subscripts cannot.
		(capacity == len(b.data))

	if !ok {
		invNum = 1 // invariant violation number.  Lame but effective.
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

/////////////////////////////////////////////////////////////////////////////////////////////////

// New() allocates and initializes a new ring buffer of specified size
func New(n int) *RingBuffer {
	b := &RingBuffer{data: make([]RingElement, n), // Contents
		in: 0, out: 0, size: 0}
	b.invariant()
	return b
}

// next() does a 'wrapping increment' of a subscript to point to the next element.
func (b *RingBuffer) next(subscript int) int {
	subscript++
	if subscript >= cap(b.data) { // I suspect this is quicker than a modulus calculation.
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

// Read fetches the next element from the ring buffer.
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

// Any left to read?
func (b *RingBuffer) HasAny() bool {
	return b.size > 0
}

// Obliterate, Purge, and Remove the contents of the ring buffer.
// Support your local Garbage Collector!
func (b *RingBuffer) Clear() {
	b.in, b.out, b.size = 0, 0, 0
	for i := 0; i < len(b.data); i++ { // Remove dangling references to avoid leaks.
		b.data[i] = nil
	}
	b.invariant()
	b.data = nil // Let GC collect the array.
}
