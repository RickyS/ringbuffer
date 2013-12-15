package ringbuffer

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	//"math/rand"
	//"testing"
)

///// Inspect the INTERNAL state of the ring buffer. ////
//var invNum int // invNum is an error code.

type DbgRingElement int

var exemplar DbgRingElement = 17 // For type checking.

var iReadCnt, iWriteCnt int = 0, 0

var wValue DbgRingElement = 0   // increasing as the test case.
var Expected DbgRingElement = 0 // wValue supposed to turn into Expected at the other end.

// The conditions checked here can best be understood by drawing the obvious diagram of the array.
func (b *RingBuffer) invariants() bool {

	if nil == b {
		invNum = 17
		return false
	}

	if (nil == b.data) && (0 == b.in) && (0 == b.out) && (0 == b.size) {
		// The RingBuffer has been nilled out by calling Clear()
		return true // All is good.
	}

	capacity := cap(b.data)
	invNum = 0
	ok := (0 <= b.in) && (b.in < capacity) &&
		(0 <= b.out) && (b.out < capacity) &&
		(0 <= b.size) && (b.size <= capacity) && // size can equal capacity.  Subscripts cannot.
		(capacity == len(b.data) &&
			(capacity > 0))

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
	xR, xW, dR1, dW1, makR, makW, changeCnt, fR, phaseCnt, interveneCnt, doiReadCnt, iSkipCnt, iWriteErr int
)

/*********************************
// Randomishly write and read, checking the invariants.
func TestInterleaved(t *testing.T) {
	//fmt.Println("———————→ Interleaved. ←———————")
	const bufferSize = 450
	const maxLoops = 6174 // why not use the "Mysterious Number of Keprekar"?
	SkipConvey("Interleaved Randomly", t, func() {
		So(maxLoops, ShouldBeGreaterThan, bufferSize)
		var b *RingBuffer
		So(b, ShouldBeNil)
		r := rand.New(rand.NewSource(99))
		b = New(bufferSize)
		So(b, ShouldNotBeNil)
		So(b.data, ShouldNotBeNil)

		So(b.invariants(), ShouldBeTrue)
		So(cap(b.data), ShouldEqual, bufferSize)
		b.Dump()

		So(b.Full(), ShouldBeFalse)
		// for b.Full() == false {
		//	e := b.WriteVer()
		//	So(e, ShouldBeNil)
		// }
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
					_ = b.ReadVer()
					xR++
				} else {
					iSkipCnt++ // Avoid errors.  Makes calculations simpler.
				}
			} else {
				//So(b, ShouldNotBeNil)
				//So(b.data, ShouldNotBeNil)
				So(cap(b.data), ShouldBeGreaterThan, 0)
				which := fmt.Sprintf("Intlv Wv %2d, Leng %d, cap %d, b %08p\n", wValue, b.Leng(), cap(b.data), b)
				//b.Dump()
				Convey(which, func() {
					b.WriteVer() // This provides the value to write: wValue.
					xW++
				})
			}
		}

		for b.HasAny() {
			x := b.ReadVer()
			So(x, ShouldHaveSameTypeAs, exemplar)
			fR++
		}

		fmt.Printf("iReadCnt %4d, iWriteCnt %4d, iSkipCnt %4d, (sum %4d), Expected %4d: Leftover %4d\n",
			iReadCnt, iWriteCnt, iSkipCnt,
			iReadCnt+iWriteCnt+iSkipCnt, Expected, b.Leng())
		fmt.Printf("InterveneCnt %4d, doiReadCnt %4d, xR %4d, xW %4d, dR1 %4d, dW1 %4d, ph %d\n",
			interveneCnt, doiReadCnt, xR, xW, dR1, dW1, phaseCnt)
		fmt.Printf("makR %4d, makW %4d, changeCnt %4d, fR %4d, b %08p\n", makR, makW, changeCnt, fR, b)
		So(b.invariants(), ShouldBeTrue)
		b.Clear()
		So(b.invariants(), ShouldBeTrue)
	})
}
***************************/
//////

// ReadVer and WriteVer are for putting stuff in numeric sequence to check that
// it comes out in the same numeric sequence.
func (b *RingBuffer) WriteVer() error {
	b.Dump()
	SkipSo(cap(b.data), ShouldBeGreaterThan, 0)
	tmp := b.WriteDet(wValue)
	if nil == tmp { // error return is NOT a bug in our package.  May be a bug by the user.
		wValue++ // wValue turns to Expected at the other end...
	}
	return tmp
}

// Put 'wValue' INTO the ring, Should get same value OUT as 'Expected'

// More debuggishness
func (b *RingBuffer) ReadVer() DbgRingElement {
	var tmp DbgRingElement
	panic("ReadVer")
	Convey("ReadVer", func() {
		tmp := b.ReadDet()
		So(tmp, ShouldHaveSameTypeAs, exemplar)
		So(tmp, ShouldEqual, Expected)
		Expected = tmp + 1 // also re-synchronize if error found.
	})
	return tmp
}

//  ReadDet and WriteDet call ringbuffer and make basic checks on each call.
func (b *RingBuffer) WriteDet(datum DbgRingElement) error {
	var err error
	Convey("WriteDet", func() {
		SkipConvey(" Checking ", func() {
			So(b, ShouldNotBeNil)
			So(b.invariants(), ShouldBeTrue)
			So(b.data, ShouldNotBeNil)
			So(cap(b.data), ShouldBeGreaterThan, 0)
		})
		//So(b.invariants(), ShouldBeTrue)
		preFull := b.Full()
		err = b.Write(datum)
		isFull := b.Full()
		if preFull {
			So(isFull, ShouldBeTrue)
		}
		if !isFull {
			So(preFull, ShouldBeFalse)
		}
		errReturned := (nil != err)
		if errReturned {
			fmt.Printf("\nWriteD err '%v', isFull %v, preFull %v, datum %v, iWriteCnt %3d, Leng %3d, b %08p\n",
				err, isFull, preFull, datum, iWriteCnt, b.Leng(), b)
			b.Dump()
		}
		//So(isFull, ShouldEqual, errReturned)
		if !errReturned {
			iWriteCnt++
		} else {
			iWriteErr++
		}
		//b.Dump()
	})
	return err
}

// Testing code
func (b *RingBuffer) ReadDet() DbgRingElement { // More debuggishness
	var tmp DbgRingElement
	Convey("ReadDet", func() {
		panic("ReadDet")
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
