package sequence

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
)

// NewAtomic creates a NewSequence that implements the Incrementer interface.
// It is safe for concurrent use.
func NewAtomic(params ...uint64) (*AtomicSequence, error) {
	seq := new(AtomicSequence)
	err := seq.Init(params...)
	return seq, err
}

// AtomicSequence is a basic Sequence that uses atomic instructions in Sequence methods.
// Although implementation is very close, it is safe for concurrent use.
type AtomicSequence Sequence

// Init a sequence with reasonable defaults based on the number and order of
// the numeric parameters passed into this method. By default, if no arguments
// are passed into Init, then the Sequence will be initialized as a
// monotonically increasing counter in the positive space as follows:
//
//     seq.Init() // count by 1 from 1 to MaximumBound
//
// If only a single argument is passed in, then it is interpreted as the
// maximum bound as follows:
//
//     seq.Init(100) // count by 1 from 1 until 100.
//
// If two arguments are passed in, then it is interpreted as a discrete range.
//
//     seq.Init(10, 100) // count by 1 from 10 until 100.
//
// If three arguments are passed in, then the third is the step.
//
//     seq.Init(2, 100, 2) // even numbers from 2 until 100.
//
// Both endpoints of these ranges are inclusive.
//
// Init can return a variety of errors. The most common error is if Init is
// called twice - that is that an already initialized sequence is attempted to
// be modified in a way that doesn't reset it. This is part of the safety
// features that Sequence provides. Other errors include mismatched or
// non-sensical arguments that won't initialize the Sequence properly.
// It is done in an atomic way and is safe for concurrent use.
func (s *AtomicSequence) Init(params ...uint64) error {
	if s.initialized {
		return errors.New("cannot re-initialize a sequence object")
	}
	// If no parameters, create the default sequence.
	if len(params) == 0 {
		atomic.AddUint64(&s.increment, 1)
		atomic.AddUint64(&s.minvalue, MinimumBound)
		atomic.AddUint64(&s.maxvalue, MaximumBound)
	}

	// If a single parameter create a maximal bounding.
	if len(params) == 1 {

		// Ensure that the parameter is greater than the minimum value.
		if params[0] < MinimumBound {
			return errors.New("must specify a maximal value greater than 0")
		}

		atomic.AddUint64(&s.increment, 1)
		atomic.AddUint64(&s.minvalue, MinimumBound)
		atomic.AddUint64(&s.maxvalue, params[0])
	}

	// If two parameters create a positive range.
	if len(params) == 2 {
		if params[1] < params[0] {
			return errors.New("for a positive increment, the maximum value must be greater than or equal to the minimum value")
		}

		if params[0] < MinimumBound || params[1] > MaximumBound {
			return errors.New("part of the range is out of bounds for positive increment")
		}

		atomic.AddUint64(&s.increment, 1)
		atomic.AddUint64(&s.minvalue, params[0])
		atomic.AddUint64(&s.maxvalue, params[1])
	}

	// If three parameters create a range with a new step.
	if len(params) == 3 {
		// The step cannot be zero
		if params[2] == 0 {
			return errors.New("must have a non-zero step to increment by")
		}

		if params[2] < 0 {
			// If the step is negative
			// TODO: This is not yet implemented since uints have to be positive.
			if params[0] < params[1] {
				return errors.New("for a negative increment, the first value must be greater than or equal to the second value")
			}

			if params[1] < MinimumBound || params[0] > MaximumBound {
				return errors.New("part of the range is out of bounds for negative increment")
			}
		} else {
			// If the step is positive
			if params[1] < params[0] {
				return errors.New("for a positive increment, the second value must be greater than or equal to the first value")
			}

			if params[0] < MinimumBound || params[1] > MaximumBound {
				return errors.New("part of the range is out of bounds for positive increment")
			}
		}

		atomic.AddUint64(&s.increment, params[2])
		atomic.AddUint64(&s.minvalue, params[0])
		atomic.AddUint64(&s.maxvalue, params[1])
	}

	// If more than three parameters then return an error.
	if len(params) > 3 {
		return errors.New("too many arguments specified")
	}

	// Ensure unsigned subtraction won't lead to a problem.
	if int(s.minvalue)-int(s.increment) < 0 {
		return errors.New("the minimum value must be less than or equal to the step")
	}

	atomic.SwapUint64(&s.current, atomic.LoadUint64(&s.minvalue)-atomic.LoadUint64(&s.increment))
	s.initialized = true
	return nil
}

// Next updates the state of the Sequence and return the next item in the
// sequence. It will return an error if either the minimum or the maximal
// value has been reached.
// It is done in an atomic way.
func (s *AtomicSequence) Next() (uint64, error) {
	atomic.AddUint64(&s.current, atomic.LoadUint64(&s.increment))

	// Check for missed minimum condition
	if atomic.LoadUint64(&s.current) < atomic.LoadUint64(&s.minvalue) {
		return 0, errors.New("reached minimum bound of the sequence")
	}

	// Check for reached maximum condition
	if atomic.LoadUint64(&s.current) > atomic.LoadUint64(&s.maxvalue) {
		return 0, errors.New("reached maximum bound of sequence")
	}

	return atomic.LoadUint64(&s.current), nil
}

// Restart the sequence by resetting the current value. This is the only
// method that allows direct manipulation of the sequence state which violates
// the monotonically increasing or decreasing rule. Use with care and as a
// fail safe if required.
// It is done in an atomic way.
func (s *AtomicSequence) Restart() error {
	// Ensure that the sequence has been initialized.
	if !s.initialized {
		return errors.New("sequence has not been initialized")
	}

	// Ensure unsigned subtraction won't lead to a problem.
	if int(atomic.LoadUint64(&s.minvalue))-int(atomic.LoadUint64(&s.increment)) < 0 {
		return errors.New("the minimum value must be less than or equal to the step")
	}

	// Set current based on the minvalue and the increment.
	atomic.SwapUint64(&s.current, atomic.LoadUint64(&s.minvalue)-atomic.LoadUint64(&s.increment))
	return nil
}

// Update the sequence to the current value. If the update value violates the
// monotonically increasing or decreasing rule, an error is returned.
// It is done in an atomic way.
func (s *AtomicSequence) Update(val uint64) error {
	// monotonically increasing error
	if atomic.LoadUint64(&s.increment) > 0 && val < atomic.LoadUint64(&s.current) {
		return errors.New("cannot decrease monotonically increasing sequence")
	}

	// monotonically decreasing error
	if atomic.LoadUint64(&s.increment) < 0 && val > atomic.LoadUint64(&s.current) {
		return errors.New("cannot increase monotonically decreasing sequence")
	}

	// Update the sequence.
	atomic.SwapUint64(&s.current, val)
	return nil
}

// Current gives the current value of this sequence atomically.
func (s *AtomicSequence) Current() (uint64, error) {
	if !s.initialized {
		return 0, errors.New("sequence has not been initialized")
	}

	if !s.IsStarted() {
		return 0, errors.New("sequence has not been started")
	}

	return atomic.LoadUint64(&s.current), nil
}

// IsStarted does atomic checks to see if this sequence has already started.
func (s *AtomicSequence) IsStarted() bool {
	if !s.initialized {
		return false
	}
	return !(atomic.LoadUint64(&s.current) < atomic.LoadUint64(&s.minvalue)) &&
		atomic.LoadUint64(&s.current) < atomic.LoadUint64(&s.maxvalue)
}

// String returns a human readable representation of this sequence.
func (s *AtomicSequence) String() string {
	d := fmt.Sprintf("incremented by %d between %d and %d", atomic.LoadUint64(&s.increment),
		atomic.LoadUint64(&s.minvalue), atomic.LoadUint64(&s.maxvalue))
	if !s.IsStarted() {
		return fmt.Sprintf("Unstarted Sequence %s", d)
	}
	return fmt.Sprintf("Sequence at %d, %s", atomic.LoadUint64(&s.current), d)
}

// Dump uses atomic Loads to Marshal current data from a AtomicSequence into a JSON object
func (s *AtomicSequence) Dump() ([]byte, error) {
	if !s.IsStarted() {
		return nil, errors.New("cannot dump an uninitialized or unstarted sequence")
	}

	data := make(map[string]uint64)
	data["current"] = atomic.LoadUint64(&s.current)
	data["increment"] = atomic.LoadUint64(&s.increment)
	data["minvalue"] = atomic.LoadUint64(&s.minvalue)
	data["maxvalue"] = atomic.LoadUint64(&s.maxvalue)

	return json.Marshal(data)
}

// Load loads data from Dump. If the input is not the same as the output from Dump() then it will return a error.
func (s *AtomicSequence) Load(data []byte) error {
	if s.initialized {
		return errors.New("cannot load into an initialized sequence")
	}

	vals := make(map[string]uint64)
	if err := json.Unmarshal(data, &vals); err != nil {
		return err
	}

	if val, ok := vals["current"]; !ok {
		return errors.New("improperly formatted data or sequence version")
	} else {
		atomic.SwapUint64(&s.current, val)
	}

	if val, ok := vals["increment"]; !ok {
		return errors.New("improperly formatted data or sequence version")
	} else {
		atomic.SwapUint64(&s.increment, val)
	}

	if val, ok := vals["minvalue"]; !ok {
		return errors.New("improperly formatted data or sequence version")
	} else {
		atomic.SwapUint64(&s.minvalue, val)
	}

	if val, ok := vals["maxvalue"]; !ok {
		return errors.New("improperly formatted data or sequence version")
	} else {
		atomic.SwapUint64(&s.maxvalue, val)
	}

	s.initialized = true
	return nil
}
