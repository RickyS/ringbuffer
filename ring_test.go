// Some of the test code is in ringbuffer/dbg_test.go, also.
package ringbuffer_test

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	//"math/rand"
	"ringbuffer"
	"testing"
)

// Error, Errorf, FailNow, Fatal, FatalIf

type kitchenSink struct { // Arbitrary type.
	words string
	nums  [4]int
}

func (pk kitchenSink) String() string {
	return "\t** " + pk.words + fmt.Sprintf(" {%3d, %3d, %3d, %3d} **",
		pk.nums[0], pk.nums[1], pk.nums[2], pk.nums[3])
}

var ksa = [...]kitchenSink{
	kitchenSink{words: "Ignore this message", nums: [...]int{0, 1, 2, 3}},
	kitchenSink{nums: [4]int{99, 98, 97, 96}, words: "this and that"},
	kitchenSink{nums: [4]int{987654321, 1234567890, 0, -1234567890}, words: "No slices allowed!"},
}

func TestKitchenSmall(t *testing.T) {
	const quantity = 11
	var rBuf = ringbuffer.New(quantity) // Create the ring buffer with the specified size.
	Convey("Blank Buffer 11", t, func() {
		So(rBuf.Leng(), ShouldEqual, 0)
	})
	for _, va := range ksa { // Add in the kitchenSink structs.
		e := rBuf.Write(va)
		if nil != e {
			t.Fatalf("ksa Oopsie\n")
		}
	}
	Convey("Full Buffer Wrote 11", t, func() {
		So(rBuf.Leng(), ShouldEqual, len(ksa))
	})
	// Now get them all out:
	for rBuf.HasAny() {
		fmt.Printf("→ %v\n", rBuf.Read())
	}
	Convey("Empty Buffer 11", t, func() {
		So(rBuf.Leng(), ShouldEqual, 0)
	})
}

func TestMillion(t *testing.T) {
	const big = 1000000 // A million
	Convey("A Million elements", t, func() {
		var rBuf = ringbuffer.New(big)
		So(rBuf, ShouldNotBeNil)
		So(0, ShouldEqual, rBuf.Leng())
		fmt.Println("———————→ Million Buffer March ←———————")

		nnn := 0
		for {
			e := rBuf.Write(&kitchenSink{words: "bogosity", nums: [4]int{nnn, 1 + nnn, 2 + nnn, 3 + nnn}})
			if nil != e {
				//fmt.Printf("Fatal Error Required: %v", e)
				break
			}
			nnn += 4
		}

		rc := 0
		for rBuf.HasAny() {
			_ = rBuf.Read()
			rc++
		}
		fmt.Println("Read ", rc, " times.")
	})
	//fmt.Println("Done")
}
