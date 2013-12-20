package ringbuffer

import (
	//"fmt"
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

//////

// ReadVer and WriteVer are for putting stuff in numeric sequence to check that
// it comes out in the same numeric sequence.
func (b *RingBuffer) WriteVer() error {
	if 0 >= cap(b.data) {
		So(cap(b.data), ShouldBeGreaterThan, 0) // check it twice to reduce output.
	}
	//fmt.Printf("WriteVer %v, b %p, cap(b.data) %d, b.size %d\n", wValue, b, cap(b.data), b.size)
	eRR := b.WriteDet(wValue)
	if nil == eRR { // error return is NOT a bug in our package.  May be a bug by the user.
		wValue++ // wValue turns to Expected at the other end...
		if 0 >= b.size {
			So(b.size, ShouldBeGreaterThan, 0)
		}
	}
	return eRR
}

//  ReadDet and WriteDet call ringbuffer and make basic checks on each call.
func (b *RingBuffer) WriteDet(datum DbgRingElement) error {
	var err error
	//fmt.Printf("WriteDet %v, b %p, cap(b.data) %d, b.size %d\n", datum, b, cap(b.data), b.size)
	//Convey("WriteDet", func() { // THIS CAUSES BUG. TODO:  TRACK DOWN.
	// SkipConvey(" Checking ", func() {
	//  So(b, ShouldNotBeNil)
	//  So(b.invariants(), ShouldBeTrue)
	//  So(b.data, ShouldNotBeNil)
	//  So(cap(b.data), ShouldBeGreaterThan, 0)
	// })
	//So(b.invariants(), ShouldBeTrue)
	preFull := b.Full()
	err = b.Write(datum)
	//fmt.Printf("WriteDeT %v, b %p, cap(b.data) %d, b.size %d\n", datum, b, cap(b.data), b.size)
	isFull := b.Full()
	if preFull {
		So(isFull, ShouldBeTrue)
	}
	if !isFull {
		So(preFull, ShouldBeFalse)
	}
	errReturned := (nil != err)
	if errReturned {
		// fmt.Printf("\nWriteD err '%v', isFull %v, preFull %v, datum %v, iWriteCnt %3d, Leng %3d, b %08p\n",
		//	err, isFull, preFull, datum, iWriteCnt, b.Leng(), b)
		// b.Dump()
	} else {
		So(b.size, ShouldBeGreaterThan, 0)
	}
	//So(isFull, ShouldEqual, errReturned)
	if !errReturned {
		iWriteCnt++
	} else {
		iWriteErr++
	}
	//b.Dump()
	//})
	return err
}

// Put 'wValue' INTO the ring, Should get same value OUT as 'Expected'

// More debuggishness
func (b *RingBuffer) ReadVer() DbgRingElement {
	var tmp DbgRingElement
	//panic("ReadVer")
	Convey("ReadVer", func() {
		tmp := b.ReadDet()
		So(tmp, ShouldHaveSameTypeAs, exemplar)
		So(tmp, ShouldEqual, Expected)
		Expected = tmp + 1 // also re-synchronize if error found.
	})
	return tmp
}

// Testing code
func (b *RingBuffer) ReadDet() DbgRingElement { // More debuggishness
	var tmp DbgRingElement
	Convey("ReadDet", func() {
		//panic("ReadDet")
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
