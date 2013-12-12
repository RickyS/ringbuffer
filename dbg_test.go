package ringbuffer_test

import (
	"fmt"
	"os"
	"ringbuffer"
	//"testing"
)

// Code to help exercise the ringbuffer package, a FIFO, LILO, and Queue.
// Call the ringbuffer routines with well-defined sequences of data so we know
// what to expect to read.
// All the real code for the ring buffer is in ringbuffer/ringbuffer.go

// Debuggishness: Don't use 'DbgRingElement', it's just for internal test.
// We read and write type DbgRingElement a lot.  And use its integerness as a check of
// integrity of the algorithm.
type DbgRingElement int

type dbgBuffer ringbuffer.RingBuffer

var ReadCnt, WriteCnt int = 0, 0

var wValue DbgRingElement = 0   // increasing as the test case.
var Expected DbgRingElement = 0 // wValue supposed to turn into Expected at the other end.

//var opVcnt int = 0

// ReadV and WriteV are for putting stuff in numeric sequence to check that
// it comes out in the same numeric sequence.
func (b *dbgBuffer) WriteV() error {
	tmp := b.WriteD(wValue)
	if nil == tmp {
		wValue++
	}
	return tmp
}

// More debuggishness
func (b *dbgBuffer) ReadV() DbgRingElement {
	tmp := b.ReadD()
	if tmp != Expected {
		fmt.Printf("\tERROR: exp %4d != act %4d\n", Expected, tmp)
		(*ringbuffer.RingBuffer)(b).Dump()
		os.Exit(2)
	}
	Expected = tmp + 1 // also re-synchronize if error found.
	return tmp
}

//  ReadD and WriteD call ringbuffer and make basic checks on each call.
func (b *dbgBuffer) WriteD(datum DbgRingElement) error {
	f := (*ringbuffer.RingBuffer)(b).Full()
	e := (*ringbuffer.RingBuffer)(b).Write(datum)
	if f != (e != nil) {
		fmt.Printf("\tERROR: full %4v but e %4v\n\t:", f, e) // Error in package.
	} else if e != nil {
		//fmt.Printf("E✔\t %q (w %d)\t\n", e, datum) // Healthy error return.
	} else {
		//fmt.Printf("W %3d\t\t:", datum)
		WriteCnt++
	}
	//(*ringbuffer.RingBuffer)(b).Dump()
	return e
}

// Testing code
func (b *dbgBuffer) ReadD() DbgRingElement { // More debuggishness
	bufLen := (*ringbuffer.RingBuffer)(b).Leng()
	var tmp DbgRingElement
	var ok bool
	if 0 < bufLen {
		tmp, ok = ((*ringbuffer.RingBuffer)(b).Read()).(DbgRingElement) // Type assertion.
		if !ok {
			fmt.Printf("ReadD Type Failure, size %4d\n", bufLen)
			(*ringbuffer.RingBuffer)(b).Dump()
		} else {
			ReadCnt++
		}
	}

	//(*ringbuffer.RingBuffer)(b).Dump()
	return tmp
}
