package data

import (
	"testing"
	"time"
)

func TestRingLength(t *testing.T) {
	r := NewRing(10, func() interface{} { return nil })

	if r.Length() != 10 {
		t.Errorf("r.Length= %d, want %d", r.Length(), 10)
	}

}

func TestRingWriteWaits(t *testing.T) {
	i := -1
	// Very fast writer
	r := NewRing(10, func() interface{} {
		i++
		return i
	})

	// slow reader, should see writes in sequence
	for j := 0; j < 100; j++ {
		time.Sleep(time.Millisecond)
		read := r.Next()
		if read != j {
			t.Errorf("got %d, want %d", read, j)
		}
	}

}

func TestRingReadOverlaps(t *testing.T) {
	i := -1
	// Very slow writer
	r := NewRing(10, func() interface{} {
		i++
		time.Sleep(time.Millisecond)
		return i
	})

	// fast reader, should see multiple times the same write
	for j := 0; j < 100; j++ {
		read := r.Next()
		if read == j {
			t.Errorf("got %d, want %d", read, j)
		}
	}

}

/*

 */
func BenchmarkWithRing(b *testing.B) {
	b.StopTimer()
	v := NewRing(100, func() interface{} { return 10 })
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v.Next()
	}
}
