package ringbuffer

import (
	"fmt"
	//"ringbuffer"
	//"testing"
)

// Dump displays the internal variables and the ENTIRE contents of the ring buffer.
func (b *RingBuffer) Dump() {
	if b.invariant() { // Calls Dump() when would return false (that is, when not ok).
		b.internalDump(``)
	}
}

// Called by Dump() and by invariant()
// 1) Show internal subscript values in parens.
// 2) Display a line of buffer contents (integers?), followed by:
// 3) A line with the array subscripts of those contents.
func (b *RingBuffer) internalDump(msg string) {
	fmt.Printf("\t(In %3d)   (Out %3d)   (Siz %3d)   (len %3d)   (cap %3d) %s [%d]\n",
		b.in, b.out, b.size, len(b.data), cap(b.data), msg, invNum)
	// invNum is an error code from the ringbuffer.invariant internal routine.
	// It must be zero.

	i, o, s := b.in, b.out, b.size // Save internal ringbuffer state
	fmt.Printf(" ")
	for i := 0; 0 < b.Leng(); i++ { // Display the ENTIRE contents of the RingBuffer!
		fmt.Printf(" %5v ", b.Read())
	}
	fmt.Println()

	b.in, b.out, b.size = i, o, s // Restore internal state.
	bOut := b.out
	for j := 0; j < b.Leng(); j++ { // Display the associate subscripts.
		ixThis := fmt.Sprintf("[%d]", bOut)
		fmt.Printf(" %6s", ixThis)
		bOut = b.next(bOut)
	}
	fmt.Println()
}
