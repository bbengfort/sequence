package sequence

import (
	"math"
	"sync"
	"testing"
)

// Ensure that the AtomicSequence object implements the Incrementer interface.
// This test is more of a compiler check since this code will fail on compile.
func TestInterfaceAtomic(t *testing.T) {
	var _ Incrementer = &AtomicSequence{}
}

// Test the creation of a default Sequence object.
func TestNewDefaultAtomic(t *testing.T) {
	seq, err := NewAtomic()
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
func TestNextAtomic(t *testing.T) {
	seq, err := NewAtomic()
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

// Test the restart functionality
func TestRestartAtomic(t *testing.T) {
	seq, err := NewAtomic()
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
func TestRestartInitErrorAtomic(t *testing.T) {
	seq := &AtomicSequence{}
	if seq.initialized {
		t.Error("sequence is initialized for some reason?")
	}

	err := seq.Restart()
	if err == nil {
		t.Error("Restart should have failed on non initialized sequence")
	}
}

// Test the update functionality
func TestUpdateAtomic(t *testing.T) {
	seq, err := NewAtomic()
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
func TestBadUpdateAtomic(t *testing.T) {
	seq, err := NewAtomic()
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

// Test the get current state functionality
func TestCurrentAtomic(t *testing.T) {
	seq, err := NewAtomic()
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
func TestCurrentInitErrorAtomic(t *testing.T) {
	seq := &AtomicSequence{}
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
func TestIsStartedAtomic(t *testing.T) {
	seq, err := NewAtomic()
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

//===========================================================================
// Test Sequence Serialization
//===========================================================================

// Test the sequence state dump and load functionality.
func TestSerializationAtomic(t *testing.T) {
	var err error
	var seqa *AtomicSequence
	var seqb *AtomicSequence

	seqa, err = NewAtomic()
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
	seqb = &AtomicSequence{}
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

//===========================================================================
// Test Initialization
//===========================================================================

// Test the creation of a default Sequence object using Init instead of New.
func TestDefaultInitAtomic(t *testing.T) {
	seq := new(AtomicSequence)
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
func TestNoDupInitAtomic(t *testing.T) {
	seq := &AtomicSequence{}
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
func Test1ArgInitAtomic(t *testing.T) {
	seq := new(AtomicSequence)
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
func Test1ArgInitZeroAtomic(t *testing.T) {
	seq := new(AtomicSequence)
	err := seq.Init(0)
	if err == nil {
		t.Error("should error when zero is passed to single init!")
	}
}

// Test the creation of a positive range
func Test2ArgInitAtomic(t *testing.T) {
	seq := new(AtomicSequence)
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
func Test2ArgBadInitAtomic(t *testing.T) {
	seq := new(AtomicSequence)
	err := seq.Init(100, 10)
	if err == nil {
		t.Error("should not allow second param to be less than first")
	}
}

// Test the creation zeroed positive range failures
func Test2ArgZeroInitAtomic(t *testing.T) {
	seq := new(AtomicSequence)
	err := seq.Init(0, 10)
	if err == nil {
		t.Error("should not allow zero valued ranges")
	}
}

// Test the creation of a positive range
func Test3ArgInitPosAtomic(t *testing.T) {
	seq := new(AtomicSequence)
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
func Test3ArgInitZeroAtomic(t *testing.T) {
	seq := new(AtomicSequence)
	err := seq.Init(10, 100, 0)
	if err == nil {
		t.Error("allowed zero value step!?")
	}
}

// Test minimum step error
func Test3ArgInitStepBangAtomic(t *testing.T) {
	seq := new(AtomicSequence)
	err := seq.Init(1, 100, 2)
	if err == nil {
		t.Error("allowed step greater than minimum value!?")
	}
}

// Ensure there is an error on more than three args
func Test4ArgInitAtomic(t *testing.T) {
	seq := new(AtomicSequence)
	err := seq.Init(10, 100, 4, 20)
	if err == nil {
		t.Error("allowed four arguments!?")
	}
}

//===========================================================================
// Test Range Boundaries
//===========================================================================

// Test that sequence goes to the maximum value then errors
func TestCeilingAtomic(t *testing.T) {
	// Create a sequence right at the maximum bound.
	seq := &AtomicSequence{MaximumBound - 1, 1, MinimumBound, MaximumBound, true}

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
func TestIncrementAtomic(t *testing.T) {
	seq := new(AtomicSequence)
	seq.Init(3, 10000, 3)

	for i := uint64(3); i < 30; i += 3 {
		j, _ := seq.Next()
		if j != i {
			t.Error("Mismatch counter value during +3 sequence")
		}
	}
}

// Test that sequence goes to the maximum value then errors on increment
func TestCeilingIncrementAtomic(t *testing.T) {
	// Create a sequence right at the maximum bound.
	seq := &AtomicSequence{MaximumBound - 1, 2, MinimumBound, MaximumBound, true}

	jdx, err := seq.Next()
	if err == nil {
		t.Error("Did not raise exception after going over maximum bound?!")
	}

	if jdx != 0 {
		t.Error("returned an overflow value")
	}
}

// Test a maximum value.
func TestMaximumAtomic(t *testing.T) {
	seq, err := NewAtomic(32132)
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
func TestRangeAtomic(t *testing.T) {
	seq, err := NewAtomic(23, 231342)
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

func TestIfAtomicIsSafeForConcurrentUse(t *testing.T) {
	seq, err := NewAtomic()
	if err != nil {
		t.Error(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		var i uint64
		for ; i < math.MaxUint16; i++ {
			seq.Next()
		}
	}()
	go func() {
		defer wg.Done()
		// keep reading current
		var i uint64
		for ; i < math.MaxUint16; i++ {
			seq.Current()
		}
	}()

	wg.Wait()
}

//===========================================================================
// Benchmarks
//===========================================================================

func BenchmarkSequenceAtomic(b *testing.B) {
	var s uint64
	f := func(u uint64) {}

	seq, err := NewAtomic()
	if err != nil {
		b.Error(err.Error())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s, _ = seq.Next()
	}
	f(s)
	b.ReportAllocs()
}
