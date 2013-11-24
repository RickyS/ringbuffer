//  Code to help exercise the ringbuffer package, a FIFO.
//  Call the ringbuffer routines with well-defined sequences of data so we know what to check for.
package ringbuffer

import (
	"fmt"
	"os"
)

// Debuggishness
// We read and write type DbgRingElement a lot.  And use its integerness as a check of 
// integrity of the algorithm.
type DbgRingElement int

var ReadCnt, WriteCnt int = 0, 0

var wValue DbgRingElement = 0   // increasing as the test case.
var Expected DbgRingElement = 0 // wValue supposed to turn into Expected at the other end.

var opVcnt int = 0


// ReadV and WriteV are for putting stuff in numeric sequence to check that
// it comes out in the same numeric sequence.
func (b *RingBuffer) WriteV() error {
	tmp := b.WriteD(wValue)
	if nil == tmp {
		wValue++
	}
	return tmp
}

func (b *RingBuffer) ReadV() DbgRingElement { // More debuggishness
	tmp := b.ReadD()
	if tmp != Expected {
		fmt.Printf("\tERROR: exp %4d != act %4d\n", Expected, tmp)
		b.Dump()
		os.Exit(2)
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

func (b *RingBuffer) ReadD() DbgRingElement { // More debuggishness
	bufLen := b.Leng()
	var tmp DbgRingElement
	var ok bool
	if 0 < bufLen {
		tmp, ok = (b.Read()).(DbgRingElement)  // Type assertion.
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

// Dump displays the internal variables and the contents of the ring buffer.
func (b *RingBuffer) Dump() {
	if b.invariant() { // Calls Dump() when would return false.
		b.internalDump(``)
	}
}

// Called by Dump() and by invariant()
// 1) Display a line of buffer contents (integers), followed by:
// 2) A line with the array subscripts of those contents.
func (b *RingBuffer) internalDump(msg string) {
	fmt.Printf("\t(In %3d)   (Out %3d)   (Siz %3d)   (len %3d)   (cap %3d) %s [%d]\n",
		b.in, b.out, b.size, len(b.data), cap(b.data), msg, invNum)
		// invNum is an error code from the ringbuffer.invariant internal routine.
		// Must be zero.

	i, o, s := b.in, b.out, b.size  // Save internal ringbuffer state
	fmt.Printf(" ")
	for i := 0; 0 < b.Leng(); i++ { // Display the contents.
		fmt.Printf(" %5d ", b.Read())
	}
	fmt.Println()

	b.in, b.out, b.size = i, o, s  // Restore internal state.
	bOut := b.out
	for j := 0; j < b.Leng(); j++ { // Display the associate subscripts.
		ixThis := fmt.Sprintf("[%d]", bOut)
		fmt.Printf(" %6s", ixThis)
		bOut = b.next(bOut)
	}
	fmt.Println()
}
