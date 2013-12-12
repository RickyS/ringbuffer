package ringbuffer_test

import (
	//"fmt"
	"ringbuffer"
	"testing"
)

type benchBuffer ringbuffer.RingBuffer

// Write only a byte at a time.
func BenchmarkInsertByteEr(b *testing.B) {
	rb := ringbuffer.New(b.N)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		e := rb.Write(byte(7))
		if nil != e {
			b.Fatalf("Write: %v\n", e)
		}
	}
}

// Write and Read, using a large-ish struct.
// it says 49 B/op.  Makes sense.
func BenchmarkKitchenMedium(b *testing.B) {
	rb := ringbuffer.New(b.N)
	b.ResetTimer()

	nnn := 0
	for i := 0; i < b.N; i++ {
		_ = rb.Write(&kitchenSink{words: "Benchmarking",
			nums: [4]int{nnn, 1 + nnn, 2 + nnn, 3 + nnn}})
		nnn += 4
	}

	//Should we check the values as they come out?
	for j := 0; j < b.N; j++ {
		ksTmp, ok := rb.Read().(*kitchenSink)
		if !ok {
			b.Fatalf("Read Type: %T\n", ksTmp)
		}
	}
}

// Same as above, copying the struct in and out.
// Copying is a little faster!!
func BenchmarkKitchenCopy(b *testing.B) {
	rb := ringbuffer.New(b.N)
	b.ResetTimer()

	nnn := 0
	for i := 0; i < b.N; i++ {
		_ = rb.Write(kitchenSink{words: "Benchmarking",
			nums: [4]int{nnn, 1 + nnn, 2 + nnn, 3 + nnn}})
		nnn += 4
	}

	//Should we check the values as they come out?
	for j := 0; j < b.N; j++ {
		ksTmp, ok := rb.Read().(kitchenSink)
		if !ok {
			b.Fatalf("(Copy) Read Type: %T\n", ksTmp)
		}
	}
}
