package sequence

import (
	"fmt"
	"testing"
)

//===========================================================================
// Basic Tests
//===========================================================================

// Ensure that the Sequence object implements the Incrementer interface.
// This test is more of a compiler check since this code will fail on compile.
func TestInterface(t *testing.T) {
	var _ Incrementer = &Sequence{}
}

// Test the creation of a default Sequence object.
func TestNewDefault(t *testing.T) {
	seq, err := New()
	if err != nil {
		t.Error(err.Error())
	}

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

//  Test the auto increment functionality
func TestNext(t *testing.T) {
	seq, err := New()
	if err != nil {
		t.Error(err.Error())
	}

	for i := uint64(1); i < 100000; i++ {
		j, err := seq.Next()
		if err != nil {
			t.Error(err.Error())
		}
		if j != i {
			t.Error("Mismatch counter value during +1 sequence")
		}
	}
}

// Example of the Basic Usage
func ExampleSequence() {

	seq, _ := New() // Create a new monotonically increasing counter

	// Fetch the first 10 sequence ids.
	for {
		idx, _ := seq.Next()
		if idx > 10 {
			break
		}

		fmt.Printf("%d ", idx)
	}

	// Output:
	// 1 2 3 4 5 6 7 8 9 10
}

// Test the restart functionality
func TestRestart(t *testing.T) {
	seq, err := New()
	if err != nil {
		t.Error(err.Error())
	}

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

// Test the non-initialized restart error
func TestRestartInitError(t *testing.T) {
	seq := &Sequence{}
	if seq.initialized {
		t.Error("sequence is initialized for some reason?")
	}

	err := seq.Restart()
	if err == nil {
		t.Error("Restart should have failed on non initialized sequence")
	}
}

// Test the update functionality
func TestUpdate(t *testing.T) {
	seq, err := New()
	if err != nil {
		t.Error(err.Error())
	}

	for i := 0; i < 100; i++ {
		_, e := seq.Next()
		if e != nil {
			t.Error("test function caused an unintended error!")
		}
	}

	if seq.current != 100 {
		t.Error("pre-restart assertion failed")
	}

	err = seq.Update(111)
	if err != nil {
		t.Error("could not update montonically increasing counter")
	}

	if seq.current != 111 {
		t.Error("update was not effective")
	}

	n, e := seq.Next()
	if e != nil {
		t.Error(e.Error())
	}
	if n != 112 {
		t.Error("sequence was not updated correctly!")
	}
}

// Test update violation errors
func TestBadUpdate(t *testing.T) {
	seq, err := New()
	if err != nil {
		t.Error(err.Error())
	}

	for i := 0; i < 100; i++ {
		_, e := seq.Next()
		if e != nil {
			t.Error("test function caused an unintended error!")
		}
	}

	if seq.current != 100 {
		t.Error("pre-restart assertion failed")
	}

	err = seq.Update(73)
	if err == nil {
		t.Error("no monotonically increasing error was returned!")
	}

	n, e := seq.Next()
	if e != nil {
		t.Error(e.Error())
	}
	if n != 101 {
		t.Error("sequence was not updated correctly!")
	}
}

// An example of updating a sequence.
func ExampleSequence_Update() {
	var idx uint64
	seq, _ := New()

	seq.Next()
	idx, _ = seq.Current()
	fmt.Println(idx)

	seq.Update(42)
	idx, _ = seq.Current()
	fmt.Println(idx)

	seq.Next()
	idx, _ = seq.Current()
	fmt.Println(idx)

	err := seq.Update(42)
	fmt.Println(err)

	// Output:
	// 1
	// 42
	// 43
	// cannot decrease monotonically increasing sequence
}

// Test the get current state functionality
func TestCurrent(t *testing.T) {
	seq, err := New()
	if err != nil {
		t.Error(err.Error())
	}

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

// Insure that unintialized sequences return an error on Current
func TestCurrentInitError(t *testing.T) {
	seq := &Sequence{}
	if seq.initialized {
		t.Error("sequence is initialized for some reason?")
	}

	idx, err := seq.Current()
	if err == nil {
		t.Error("Current should have failed on non initialized sequence")
	}

	if idx != 0 {
		t.Error("Index should be zero valued")
	}
}

// Test the is started functionality
func TestIsStarted(t *testing.T) {
	seq, err := New()
	if err != nil {
		t.Error(err.Error())
	}

	if seq.IsStarted() {
		t.Error("Unstarted sequence says it's started?!")
	}
	seq.Next()

	if !seq.IsStarted() {
		t.Error("Started sequence says it's not started?!")
	}
}

// An example of the human readable state of a sequence.
func ExampleSequence_String() {
	seq, _ := New()

	fmt.Println(seq)

	seq.Next()
	fmt.Println(seq)

	// Output:
	// Unstarted Sequence incremented by 1 between 1 and 18446744073709551614
	// Sequence at 1, incremented by 1 between 1 and 18446744073709551614
}

//===========================================================================
// Test Sequence Serialization
//===========================================================================

// Test the sequence state dump and load functionality.
func TestSerialization(t *testing.T) {
	var err error
	var seqa *Sequence
	var seqb *Sequence

	seqa, err = New()
	if err != nil {
		t.Error(err.Error())
	}

	// Update sequence a bit
	for i := 0; i < 93212; i++ {
		seqa.Next()
	}

	// Dump the data
	data, err := seqa.Dump()
	if err != nil {
		t.Error(err.Error())
	}

	// Load the new sequence
	seqb = &Sequence{}
	err = seqb.Load(data)
	if err != nil {
		t.Error(err.Error())
	}

	// Compare the sequences
	if seqa.current != seqb.current || seqa.increment != seqb.increment || seqa.minvalue != seqb.minvalue || seqa.maxvalue != seqb.maxvalue {
		t.Error("loaded sequence does not match dumped sequence")
	}

	// Ensure that seqb is initialzed
	if !seqb.IsStarted() {
		t.Error("loaded sequence isn't started?!")
	}

}

// Write a sequence to disk to be loaded later.
func ExampleSequence_Dump() {

	seq := new(Sequence)
	seq.Init()

	for i := 0; i < 10; i++ {
		seq.Next()
	}

	data, _ := seq.Dump()
	fmt.Println(string(data))

	// Output:
	// {"current":10,"increment":1,"maxvalue":18446744073709551614,"minvalue":1}
}

// An example of sequence serialization.
func ExampleSequence_Load() {

	seq := new(Sequence)
	seq.Init()

	for i := 0; i < 10; i++ {
		seq.Next()
	}

	data, _ := seq.Dump()

	sequel := new(Sequence)
	sequel.Load(data)
	fmt.Println(sequel)

	// Output:
	// Sequence at 10, incremented by 1 between 1 and 18446744073709551614
}

//===========================================================================
// Test Initialization
//===========================================================================

// Test the creation of a default Sequence object using Init instead of New.
func TestDefaultInit(t *testing.T) {
	seq := new(Sequence)
	err := seq.Init()
	if err != nil {
		t.Error(err.Error())
	}

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

// Do not allow an initialized sequence to be reinitialized
func TestNoDupInit(t *testing.T) {
	seq := &Sequence{}
	err := seq.Init()
	if err != nil {
		t.Error(err.Error())
	}

	if !seq.initialized {
		t.Error("sequence is not initialized?!")
	}

	err = seq.Init()
	if err == nil {
		t.Error("init didn't return an error after second init")
	}
}

// Test single argument init
func Test1ArgInit(t *testing.T) {
	seq := new(Sequence)
	err := seq.Init(100)
	if err != nil {
		t.Error(err.Error())
	}

	if seq.current != 0 {
		t.Error("Current (start) value not initialized correctly")
	}

	if seq.increment != 1 {
		t.Error("Increment value not initialized correctly")
	}

	if seq.minvalue != 1 {
		t.Error("Minimum value not initialized correctly")
	}

	if seq.maxvalue != 100 {
		t.Error("Maximum value not initialized correctly")
	}
}

// Ensure single argument is greater than 0
func Test1ArgInitZero(t *testing.T) {
	seq := new(Sequence)
	err := seq.Init(0)
	if err == nil {
		t.Error("should error when zero is passed to single init!")
	}
}

// Test the creation of a positive range
func Test2ArgInit(t *testing.T) {
	seq := new(Sequence)
	err := seq.Init(10, 100)
	if err != nil {
		t.Error(err.Error())
	}

	if seq.current != 9 {
		t.Error("Current (start) value not initialized correctly")
	}

	if seq.increment != 1 {
		t.Error("Increment value not initialized correctly")
	}

	if seq.minvalue != 10 {
		t.Error("Minimum value not initialized correctly")
	}

	if seq.maxvalue != 100 {
		t.Error("Maximum value not initialized correctly")
	}
}

// Test the creation a positive range failures
func Test2ArgBadInit(t *testing.T) {
	seq := new(Sequence)
	err := seq.Init(100, 10)
	if err == nil {
		t.Error("should not allow second param to be less than first")
	}
}

// Test the creation zeroed positive range failures
func Test2ArgZeroInit(t *testing.T) {
	seq := new(Sequence)
	err := seq.Init(0, 10)
	if err == nil {
		t.Error("should not allow zero valued ranges")
	}
}

// Test the creation of a positive range
func Test3ArgInitPos(t *testing.T) {
	seq := new(Sequence)
	err := seq.Init(10, 100, 5)
	if err != nil {
		t.Error(err.Error())
	}

	if seq.current != 5 {
		t.Error("Current (start) value not initialized correctly")
	}

	if seq.increment != 5 {
		t.Error("Increment value not initialized correctly")
	}

	if seq.minvalue != 10 {
		t.Error("Minimum value not initialized correctly")
	}

	if seq.maxvalue != 100 {
		t.Error("Maximum value not initialized correctly")
	}
}

// Ensure there is no zero step
func Test3ArgInitZero(t *testing.T) {
	seq := new(Sequence)
	err := seq.Init(10, 100, 0)
	if err == nil {
		t.Error("allowed zero value step!?")
	}
}

// Test minimum step error
func Test3ArgInitStepBang(t *testing.T) {
	seq := new(Sequence)
	err := seq.Init(1, 100, 2)
	if err == nil {
		t.Error("allowed step greater than minimum value!?")
	}
}

// Ensure there is an error on more than three args
func Test4ArgInit(t *testing.T) {
	seq := new(Sequence)
	err := seq.Init(10, 100, 4, 20)
	if err == nil {
		t.Error("allowed four arguments!?")
	}
}

//===========================================================================
// Test Range Boundaries
//===========================================================================

// Test that sequence goes to the maximum value then errors
func TestCeiling(t *testing.T) {
	// Create a sequence right at the maximum bound.
	seq := &Sequence{MaximumBound - 1, 1, MinimumBound, MaximumBound, true}

	idx, err := seq.Next()
	if err != nil {
		t.Error(err.Error())
	}

	if idx != MaximumBound {
		t.Error("Could not generate maximum bound value")
	}

	jdx, err := seq.Next()
	if err == nil {
		t.Error("Did not raise exception after going over maximum bound?!")
	}

	if jdx != 0 {
		t.Error("returned an overflow value")
	}
}

// Test an increment value of 3 with intermediate range.
func TestIncrement(t *testing.T) {
	seq := new(Sequence)
	seq.Init(3, 10000, 3)

	for i := uint64(3); i < 30; i += 3 {
		j, _ := seq.Next()
		if j != i {
			t.Error("Mismatch counter value during +3 sequence")
		}
	}
}

// Test that sequence goes to the maximum value then errors on increment
func TestCeilingIncrement(t *testing.T) {
	// Create a sequence right at the maximum bound.
	seq := &Sequence{MaximumBound - 1, 2, MinimumBound, MaximumBound, true}

	jdx, err := seq.Next()
	if err == nil {
		t.Error("Did not raise exception after going over maximum bound?!")
	}

	if jdx != 0 {
		t.Error("returned an overflow value")
	}
}

// Test a maximum value.
func TestMaximum(t *testing.T) {
	seq, err := New(32132)
	if err != nil {
		t.Error(err.Error())
	}

	for i := uint64(1); i < 32133; i++ {
		j, err := seq.Next()
		if err != nil {
			t.Error(err.Error())
		}

		if i != j {
			t.Error("mismatch during maximum value testing")
		}
	}

	for i := 0; i < 10; i++ {
		val, err := seq.Next()
		if err == nil {
			t.Error("should have raised error after maximum reached")
		}
		if val != 0 {
			t.Error("returning non-zero valued response!")
		}
	}
}

// Test a range, note that the monotonically increasing counter range will
// timeout before tests can be completed.
func TestRange(t *testing.T) {
	seq, err := New(23, 231342)
	if err != nil {
		t.Error(err.Error())
	}

	for i := uint64(23); i < 231343; i++ {
		j, err := seq.Next()
		if err != nil {
			t.Error(err.Error())
		}

		if i != j {
			t.Error("mismatch during maximum value testing")
		}
	}

	for i := 0; i < 10; i++ {
		val, err := seq.Next()
		if err == nil {
			t.Error("should have raised error after maximum reached")
		}
		if val != 0 {
			t.Error("returning non-zero valued response!")
		}
	}
}

//===========================================================================
// Benchmarks
//===========================================================================

func BenchmarkSequence(b *testing.B) {

	for i := 0; i < b.N; i++ {
		seq, err := New()
		if err != nil {
			b.Error(err.Error())
		}

		for j := 0; j < b.N; j++ {
			seq.Next()
		}
	}
}
