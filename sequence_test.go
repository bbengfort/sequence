package sequence

import (
	"fmt"
	"testing"
)

func TestInterface(t *testing.T) {
	var _ Incrementer = &Sequence{}
}

func TestNew(t *testing.T) {
	seq := New()
	if seq.current != 0 {
		t.Error("Current (start) value not initialized correctly")
	}

	if seq.increment != 1 {
		t.Error("Increment value not initialized correctly")
	}

	if seq.minvalue != 1 {
		t.Error("Minimum value not initialized correctly")
	}

	if seq.maxvalue != ^uint64(0)-1 {
		t.Error("Maximum value not initialized correctly")
	}
}

func TestNext(t *testing.T) {
	seq := New()
	for i := uint64(1); i < 10; i++ {
		j, _ := seq.Next()
		if j != i {
			t.Error("Mismatch counter value during +1 sequence")
		}
	}
}

func TestIncrement(t *testing.T) {
	seq := new(Sequence)
	seq.Init(0, 3)
	for i := uint64(3); i < 30; i += 3 {
		j, _ := seq.Next()
		if j != i {
			t.Error("Mismatch counter value during +2 sequence")
		}
	}
}

func TestMinError(t *testing.T) {
	seq := new(Sequence)
	seq.Init(0, 1, 4)
	j, err := seq.Next()
	if err == nil {
		t.Error("On incorrect minimum floor, the function did not return an error")
	}

	if j != 0 {
		t.Error("On error, sequence returned some integer other than zero")
	}
}

func TestMaxError(t *testing.T) {
	seq := new(Sequence)
	seq.Init(0, 1, 1, 10)
	for i := 0; i < 10; i++ {
		_, e := seq.Next()
		if e != nil {
			t.Error("Test function reached error prematurely")
		}
	}

	j, err := seq.Next()
	if err == nil {
		t.Error("On reaching maximum the function did not return an error")
	}

	if j != 0 {
		t.Error("On error, sequence returned some integer other than zero")
	}
}

func TestRestart(t *testing.T) {
	seq := New()
	for i := 0; i < 100; i++ {
		_, e := seq.Next()
		if e != nil {
			t.Error("test function caused an unintended error!")
		}
	}

	if seq.current != 100 {
		t.Error("pre-restart assertion failed")
	}

	seq.Restart()

	if seq.IsStarted() {
		t.Error("restart was not successful")
	}

	for i := 0; i < 100; i++ {
		_, e := seq.Next()
		if e != nil {
			t.Error("test function caused an unintended error!")
		}
	}

	if seq.current != 100 {
		t.Error("post-restart assertion failed")
	}

}

func TestCurrent(t *testing.T) {
	seq := New()

	j, e := seq.Current()
	if e == nil {
		t.Error("Unstarted sequence did not return an error for current")
	}
	if j != 0 {
		t.Error("On error, sequence returned some integer other than zero")
	}

	for i := uint64(1); i < 100; i++ {
		_, e := seq.Next()
		if e != nil {
			t.Error("Sequence raised an unintended error!")
		}

		if k, _ := seq.Current(); i != k {
			t.Error("Current() does not match correct value for sequence")
		}
	}
}

func TestIsStarted(t *testing.T) {
	seq := New()
	if seq.IsStarted() {
		t.Error("Unstarted sequence says it's started?!")
	}
	seq.Next()

	if !seq.IsStarted() {
		t.Error("Started sequence says it's not started?!")
	}
}

func ExampleSequenceString() {
	seq := New()
	fmt.Println(seq)

	seq.Next()
	fmt.Println(seq)

	// Output:
	// Unstarted Sequence incremented by 1 between 1 and 18446744073709551614
	// Sequence at 1, incremented by 1 between 1 and 18446744073709551614
}

func TestInitNoParams(t *testing.T) {
	seq := new(Sequence)
	seq.Init()
	if seq.current != 0 {
		t.Error("Current (start) value not initialized correctly")
	}

	if seq.increment != 1 {
		t.Error("Increment value not initialized correctly")
	}

	if seq.minvalue != 1 {
		t.Error("Minimum value not initialized correctly")
	}

	if seq.maxvalue != ^uint64(0)-1 {
		t.Error("Maximum value not initialized correctly")
	}
}

func TestInit1Params(t *testing.T) {
	seq := new(Sequence)
	seq.Init(10)
	if seq.current != 10 {
		t.Error("Current (start) value not initialized correctly")
	}

	if seq.increment != 1 {
		t.Error("Increment value not initialized correctly")
	}

	if seq.minvalue != 1 {
		t.Error("Minimum value not initialized correctly")
	}

	if seq.maxvalue != ^uint64(0)-1 {
		t.Error("Maximum value not initialized correctly")
	}
}

func TestInit2Params(t *testing.T) {
	seq := new(Sequence)
	seq.Init(10, 2)
	if seq.current != 10 {
		t.Error("Current (start) value not initialized correctly")
	}

	if seq.increment != 2 {
		t.Error("Increment value not initialized correctly")
	}

	if seq.minvalue != 1 {
		t.Error("Minimum value not initialized correctly")
	}

	if seq.maxvalue != ^uint64(0)-1 {
		t.Error("Maximum value not initialized correctly")
	}
}

func TestInit3Params(t *testing.T) {
	seq := new(Sequence)
	seq.Init(10, 2, 12)
	if seq.current != 10 {
		t.Error("Current (start) value not initialized correctly")
	}

	if seq.increment != 2 {
		t.Error("Increment value not initialized correctly")
	}

	if seq.minvalue != 12 {
		t.Error("Minimum value not initialized correctly")
	}

	if seq.maxvalue != ^uint64(0)-1 {
		t.Error("Maximum value not initialized correctly")
	}
}

func TestInit4Params(t *testing.T) {
	seq := new(Sequence)
	seq.Init(10, 2, 12, 64)
	if seq.current != 10 {
		t.Error("Current (start) value not initialized correctly")
	}

	if seq.increment != 2 {
		t.Error("Increment value not initialized correctly")
	}

	if seq.minvalue != 12 {
		t.Error("Minimum value not initialized correctly")
	}

	if seq.maxvalue != 64 {
		t.Error("Maximum value not initialized correctly")
	}
}

func TestInit5Params(t *testing.T) {
	seq := new(Sequence)
	seq.Init(10, 2, 12, 64, 128)
	if seq.current != 10 {
		t.Error("Current (start) value not initialized correctly")
	}

	if seq.increment != 2 {
		t.Error("Increment value not initialized correctly")
	}

	if seq.minvalue != 12 {
		t.Error("Minimum value not initialized correctly")
	}

	if seq.maxvalue != 64 {
		t.Error("Maximum value not initialized correctly")
	}
}

func BenchmarkSequence(b *testing.B) {

	for i := 0; i < b.N; i++ {
		seq := New()
		for j := 0; j < b.N; j++ {
			seq.Next()
		}
	}
}
