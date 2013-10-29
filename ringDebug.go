package ringbuffer

import (
	"fmt"
	"os"
)

// Debuggishness
type DbgRingElement int

var ReadCnt, WriteCnt int = 0, 0

var wValue DbgRingElement = 0   // increasing as the test case.
var Expected DbgRingElement = 0 // wValue supposed to turn into Expected at the other end.

var opVcnt int = 0

func (b *RingBuffer) WriteV() error {
	tmp := b.WriteD(wValue)
	if nil == tmp {
		wValue++
	}
	return tmp
}

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

func (b *RingBuffer) ReadD() DbgRingElement { // More debuggishness
	bufLen := b.Leng()
	var tmp DbgRingElement
	var ok bool
	if 0 < bufLen {
		tmp, ok = (b.Read()).(DbgRingElement)
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
	if b.invariant() { // Calls Dump() if would return false.
		b.internalDump(``)
	}
}

// Called by Dump() and by invariant()
func (b *RingBuffer) internalDump(msg string) {
	fmt.Printf("\t(In %3d)   (Out %3d)   (Siz %3d)   (len %3d)   (cap %3d) %s [%d]\n",
		b.in, b.out, b.size, len(b.data), cap(b.data), msg, invNum)

	i, o, s := b.in, b.out, b.size
	fmt.Printf(" ")
	for i := 0; 0 < b.Leng(); i++ { // Display the contents.
		fmt.Printf(" %5d ", b.Read())
	}
	fmt.Println()

	b.in, b.out, b.size = i, o, s
	bOut := b.out
	for j := 0; j < b.Leng(); j++ { // Display the associate subscripts.
		ixThis := fmt.Sprintf("[%d]", bOut)
		fmt.Printf(" %6s", ixThis)
		bOut = b.next(bOut)
	}
	fmt.Println()
}

/*
func (b *RingBuffer) dumper(s string) {
	fmt.Printf("\t%s: ", s)
	b.Dump()
}
*/
