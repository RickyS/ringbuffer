// Some of the test code is in ringbuffer/ringDebug.go, that
// should change.
package ringbuffer_test

import (
	"fmt"
	"math/rand"
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
	var rBuf = ringbuffer.New(11) // Create the ring buffer with the specified size.
	for _, va := range ksa {      // Add in the kitchenSink structs.
		e := rBuf.Write(va)
		if nil != e {
			t.Fatalf("ksa Oopsie\n")
		}
	}
	// Now get them all out:
	for rBuf.HasAny() {
		fmt.Printf("→ %v\n", rBuf.Read())
	}
}

func TestMillion(t *testing.T) {
	const big = 1000000 // A million
	var rBuf = ringbuffer.New(big)
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

	//fmt.Println("Done")
}

func TestInterleaved(t *testing.T) {
	fmt.Println("———————→ Interleaved ←———————")
	r := rand.New(rand.NewSource(99))
	b := ringbuffer.New(45)
	b.Dump()
	SkipCnt := 0
	for i := 0; i < 3317; i++ {
		x := r.Intn(512)
		doRead := 0 == (1 & x)              // isOdd ?
		if doRead && (i > (6 + b.Leng())) { // no Reading until we've overflowed the buffer.
			if 0 < b.Leng() {
				_ = (*dbgBuffer)(b).ReadV()
			} else {
				SkipCnt++
			}
		} else {
			(*dbgBuffer)(b).WriteV() // This provides the value to write.
		}
	}
	for b.HasAny() {
		_ = (*dbgBuffer)(b).ReadV()
	}
	b.Dump()
}
