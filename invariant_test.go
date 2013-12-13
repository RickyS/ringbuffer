package ringbuffer

import (
	//"fmt"
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

var myT *testing.T

// Randomishly write and read, checking the invariants.
func TestInterleaved(t *testing.T) {
	myT = t
	//fmt.Println("———————→ Interleaved. ←———————")
	const bufferSize = 450
	const maxLoops = 6174 // why not use the "Mysterious Number of Keprekar"?
	Convey("Interleaved Randomly", t, func() {
		So(maxLoops, ShouldBeGreaterThan, bufferSize)
		var b *RingBuffer
		So(b, ShouldBeNil)
		r := rand.New(rand.NewSource(99))
		b = New(bufferSize)
		So(b, ShouldNotBeNil)

		So(b.invariants(), ShouldBeTrue)
		b.Dump()
		SkipCnt := 0
		var phaseCnt = 0 // Have we overflowed?
		for i := 0; i < maxLoops; i++ {
			x := r.Intn(512)
			doRead := (0 == (1 & x)) && (i > (6 + b.Leng())) // isOdd && no Reading until
			// we've overflowed the buffer.

			if 1 == phaseCnt {
				intervene := (x & 0x102) == 0x102 // 2 bits match 1/4 of the time.
				bLeng := b.Leng()
				if bLeng > ((bufferSize * 2) / 3) {
					if intervene {
						doRead = true // Usually read when fullish.
					}
				} else if bLeng < (bufferSize / 3) {
					if intervene {
						doRead = false // When buffer low, write more.
					}
				}
			}

			if doRead {
				phaseCnt = 1
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
			x := b.ReadV()
			So(x, ShouldHaveSameTypeAs, exemplar)
		}
	})
}

//////
type DbgRingElement int

var exemplar DbgRingElement = 17 // For type checking.

var ReadCnt, WriteCnt int = 0, 0

var wValue DbgRingElement = 0   // increasing as the test case.
var Expected DbgRingElement = 0 // wValue supposed to turn into Expected at the other end.

// ReadV and WriteV are for putting stuff in numeric sequence to check that
// it comes out in the same numeric sequence.
func (b *RingBuffer) WriteV() error {
	tmp := b.WriteD(wValue)
	if nil == tmp { // error return is NOT a bug.
		wValue++ // wValue turns to Expected at the other end...
	}
	return tmp
}

// More debuggishness
func (b *RingBuffer) ReadV() DbgRingElement {
	var tmp DbgRingElement
	Convey("ReadV", func() {
		tmp := b.ReadD()
		So(tmp, ShouldHaveSameTypeAs, exemplar)
		So(b.invariants(), ShouldBeTrue)
		So(tmp, ShouldEqual, Expected)
		Expected = tmp + 1 // also re-synchronize if error found.
	})
	return tmp
}

//  ReadD and WriteD call ringbuffer and make basic checks on each call.
func (b *RingBuffer) WriteD(datum DbgRingElement) error {
	var e error
	Convey("WriteD", func() {
		f := b.Full()
		e = b.Write(datum)
		So(b.invariants(), ShouldBeTrue)
		errReturned := (nil != e)
		So(f, ShouldEqual, errReturned)
		if !errReturned {
			WriteCnt++
		}
		//b.Dump()
	})
	return e
}

// Testing code
func (b *RingBuffer) ReadD() DbgRingElement { // More debuggishness
	var tmp DbgRingElement
	Convey("ReadD", func() {
		bufLen := b.Leng()
		So(b.invariants(), ShouldBeTrue)

		var ok bool
		if 0 < bufLen {
			tmp, ok = (b.Read()).(DbgRingElement) // Type assertion.
			So(tmp, ShouldHaveSameTypeAs, exemplar)
			So(ok, ShouldBeTrue)
			So(b.invariants(), ShouldBeTrue)
			ReadCnt++
		}
	})
	//b.Dump()
	return tmp
}
