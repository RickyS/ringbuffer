package ringbuffer

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

///// Inspect the INTERNAL state of the ring buffer. ////
//var invNum int // invNum is an error code.

type DbgRingElement int

var exemplar DbgRingElement = 17 // For type checking.

var iReadCnt, iWriteCnt int = 0, 0

var wValue DbgRingElement = 0   // increasing as the test case.
var Expected DbgRingElement = 0 // wValue supposed to turn into Expected at the other end.

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

var (
	xR, xW, dR1, dW1, makR, makW, changeCnt, fR int
)

// Randomishly write and read, checking the invariants.
func TestInterleaved(t *testing.T) {
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
		iSkipCnt := 0
		var phaseCnt = 0 // Have we overflowed?
		var interveneCnt = 0
		var doiReadCnt = 0

		for i := 0; i < maxLoops; i++ {
			x := r.Intn(0x600)
			doRead := (0 == (1 & x)) && (i > (6 + bufferSize)) // isOdd && no Reading until
			// we've overflowed the buffer.
			if doRead {
				dR1++
			} else {
				dW1++
			}

			if (0 == phaseCnt) && b.Full() {
				phaseCnt = 1
			}
			oldDoRead := doRead
			if 1 == phaseCnt {
				intervene := (x & 0x12) == 0x12 // 2 bits are 1: 1/4 of the time.
				if intervene {
					interveneCnt++
					bLeng := b.Leng()
					if bLeng > ((bufferSize * 2) / 3) {
						doRead = true // Usually read when fullish.
						makR++
					} else if bLeng < (bufferSize / 3) {
						doRead = false // When buffer low, write more.
						makW++
					}
				}
			}
			if oldDoRead != doRead {
				changeCnt++
			}

			if doRead {
				doiReadCnt++
				if b.HasAny() {
					_ = b.ReadV()
					xR++
				} else {
					iSkipCnt++ // Avoid errors.  Makes calculations simpler.
				}
			} else {
				b.WriteV() // This provides the value to write.
				xW++
			}
		}

		for b.HasAny() {
			x := b.ReadV()
			So(x, ShouldHaveSameTypeAs, exemplar)
			fR++
		}

		//So((iReadCnt + iSkipCnt))
		fmt.Printf("iReadCnt %4d, iWriteCnt %4d, iSkipCnt %4d, (sum %4d), Expected %4d: Leftover %4d\n",
			iReadCnt, iWriteCnt, iSkipCnt,
			iReadCnt+iWriteCnt+iSkipCnt, Expected, b.Leng())
		fmt.Printf("InterveneCnt %4d, doiReadCnt %4d, xR %4d, xW %4d, dR1 %4d, dW1 %4d, ph %d\n",
			interveneCnt, doiReadCnt, xR, xW, dR1, dW1, phaseCnt)
		fmt.Printf("makR %4d, makW %4d, changeCnt %4d, fR %4d\n", makR, makW, changeCnt, fR)

	})
}

//////

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
		So(tmp, ShouldEqual, Expected)
		Expected = tmp + 1 // also re-synchronize if error found.
	})
	return tmp
}

//  ReadD and WriteD call ringbuffer and make basic checks on each call.
func (b *RingBuffer) WriteD(datum DbgRingElement) error {
	var e error
	Convey("WriteD", func() {
		isFull := b.Full()
		So(b.invariants(), ShouldBeTrue)
		e = b.Write(datum)
		So(b.invariants(), ShouldBeTrue)
		errReturned := (nil != e)
		So(isFull, ShouldEqual, errReturned)
		if !errReturned {
			iWriteCnt++
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
		So(bufLen, ShouldBeGreaterThan, 0)

		var ok bool
		tmp, ok = (b.Read()).(DbgRingElement) // Type assertion.
		So(tmp, ShouldHaveSameTypeAs, exemplar)
		So(ok, ShouldBeTrue)
		So(b.invariants(), ShouldBeTrue)
		iReadCnt++
	})
	//b.Dump()
	return tmp
}
