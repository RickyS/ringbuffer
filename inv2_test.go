package ringbuffer

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"math/rand"
	"testing"
)

var (
	hasCnt, zeroCnt int
)

// Simplified version of "TestInterleaved", which was too complex to debug.
func TestRand(t *testing.T) {
	//fmt.Println("———————→ TestRand ←———————")
	const bufferSize = 45
	const maxLoops = 111
	// const bufferSize = 6
	// const maxLoops = 10
	phaseCnt = 0
	Convey("TestRand", t, func() {
		So(maxLoops, ShouldBeGreaterThan, bufferSize)
		var b *RingBuffer
		var x int
		So(b, ShouldBeNil)
		r := rand.New(rand.NewSource(99))
		b = New(bufferSize)
		So(b, ShouldNotBeNil)
		So(b.data, ShouldNotBeNil)

		So(b.invariants(), ShouldBeTrue)
		So(cap(b.data), ShouldEqual, bufferSize)
		So(len(b.data), ShouldEqual, bufferSize)
		So(b.Leng(), ShouldEqual, 0)
		So(b.HasAny(), ShouldBeFalse)
		So(b.Full(), ShouldBeFalse)
		//dumpData(b)
		//b.Dump()
		zeroCnt = 0
		for i := 0; i < maxLoops; i++ {
			x = r.Intn(0x601)
			doRead := (1 == (1 & x)) && (i > (6 + bufferSize)) // isOdd && no Reading until
			// we've overflowed the buffer.

			if 0 == b.Leng() {
				doRead = false
				zeroCnt++
			} else {
				hasCnt++
			}
			if doRead {
				//panic("dR1")
				dR1++
			} else {
				dW1++
			}
			if (0 == phaseCnt) && b.Full() {
				phaseCnt = 1
			}
			So(cap(b.data), ShouldEqual, bufferSize) // Now this is redundant.  Used be needed.
			So(b.size, ShouldEqual, b.Leng())

			if doRead {
				doiReadCnt++
				_ = b.ReadVer()
				xR++
			} else {
				erra := b.WriteVer() // This provides the value to write: wValue.
				if nil == erra {
					So(b.Leng(), ShouldBeGreaterThan, 0)
				}
				xW++
			}
		}
		dumpData(b)
		fmt.Printf("Done.\n")
	})
}

func dumpData(b *RingBuffer) {
	fmt.Printf("\niReadCnt %4d, iWriteCnt %4d, iWriteErr %4d, iSkipCnt %4d, (sum %4d), Expected %4d: Leftover %4d\n",
		iReadCnt, iWriteCnt, iWriteErr, iSkipCnt,
		iReadCnt+iWriteCnt+iSkipCnt, Expected, b.Leng())
	fmt.Printf("InterveneCnt %4d, doiReadCnt %4d, xR %4d, xW %4d, dR1 %4d, dW1 %4d, phaseCnt %d\n",
		interveneCnt, doiReadCnt, xR, xW, dR1, dW1, phaseCnt)
	fmt.Printf("makR %4d, makW %4d, changeCnt %4d, fR %4d, b %08p, hasCnt %d, zeroCnt %d\n",
		makR, makW, changeCnt, fR, b, hasCnt, zeroCnt)
	//b.Dump()
	fmt.Printf("——\n")

}
