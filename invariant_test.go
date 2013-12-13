package ringbuffer

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

///// Inspect the INTERNAL state of the ring buffer. ////
//var invNum int // invNum is an error code.

// The conditions checked here can best be understood by drawing the obvious diagram of the array.
func (b *RingBuffer) invariants() bool { // You can remove this function and all ref to it.

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

	return ok
}

// Randomishly write and read, checking the invariants.
func TestInterleaved(t *testing.T) {
	fmt.Println("———————→ Interleaved ←———————")
	const bufferSize = 450
	Convey("Interleaved Randomly", t, func() {
		var b *RingBuffer
		So(b, ShouldBeNil)
		r := rand.New(rand.NewSource(99))
		b = New(bufferSize)
		So(b, ShouldNotBeNil)
		So(b.invariants(), ShouldBeTrue)
		b.Dump()
		SkipCnt := 0
		for i := 0; i < 3317; i++ {
			x := r.Intn(512)
			doRead := 0 == (1 & x)              // isOdd ?
			if doRead && (i > (6 + b.Leng())) { // no Reading until we've overflowed the buffer.
				if 0 < b.Leng() {
					_ = b.ReadV()
				} else {
					SkipCnt++
				}
			} else {
				b.WriteV() // This provides the value to write.
			}
		}
		for b.HasAny() {
			_ = b.ReadV()
		}
	})
}

//////
type DbgRingElement int

/// type RingBuffer RingBuffer

var ReadCnt, WriteCnt int = 0, 0

var wValue DbgRingElement = 0   // increasing as the test case.
var Expected DbgRingElement = 0 // wValue supposed to turn into Expected at the other end.

//var opVcnt int = 0

// ReadV and WriteV are for putting stuff in numeric sequence to check that
// it comes out in the same numeric sequence.
func (b *RingBuffer) WriteV() error {
	tmp := b.WriteD(wValue)
	if nil == tmp {
		wValue++
	}
	return tmp
}

// More debuggishness
func (b *RingBuffer) ReadV() DbgRingElement {
	tmp := b.ReadD()
	if tmp != Expected {
		fmt.Printf("\tERROR: exp %4d != act %4d\n", Expected, tmp)
		// b.Dump()
		// os.Exit(2)
	}
	Expected = tmp + 1 // also re-synchronize if error found.
	return tmp
}

//  ReadD and WriteD call ringbuffer and make basic checks on each call.
func (b *RingBuffer) WriteD(datum DbgRingElement) error {
	f := b.Full()
	e := b.Write(datum)
	if f != (e != nil) {
		fmt.Printf("\tERROR: full %4v but e %4v\n\t:", f, e) // Error in package.
	} else if e != nil {
		//fmt.Printf("E✔\t %q (w %d)\t\n", e, datum) // Healthy error return.
	} else {
		//fmt.Printf("W %3d\t\t:", datum)
		WriteCnt++
	}
	//b.Dump()
	return e
}

// Testing code
func (b *RingBuffer) ReadD() DbgRingElement { // More debuggishness
	bufLen := b.Leng()
	var tmp DbgRingElement
	var ok bool
	if 0 < bufLen {
		tmp, ok = (b.Read()).(DbgRingElement) // Type assertion.
		if !ok {
			fmt.Printf("ReadD Type Failure, size %4d\n", bufLen)
			b.Dump()
		} else {
			ReadCnt++
		}
	}

	//b.Dump()
	return tmp
}
