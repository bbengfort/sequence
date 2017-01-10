package sequence

import (
	"errors"
	"fmt"
)

const maxuint64 = ^uint64(0) - 1

// New returns a new default sequence (infinite increment by 1 from 1).
func New() *Sequence {
	return &Sequence{0, 1, 1, maxuint64}
}

//=============================================================================

// Sequence implements an AutoIncrement counter class similar to the
// PostgreSQL sequence object.
type Sequence struct {
	current   uint64 // The current value of the sequence
	increment uint64 // The value to increment by (usually 1)
	minvalue  uint64 // The minimum value of the counter (usually 1)
	maxvalue  uint64 // The max value of the counter (usually bounded by type)
}

// Incrementer defines the interface for sequence-like objects.
type Incrementer interface {
	Init(params ...uint64)    // Initialize the Incrementer with values
	Next() (uint64, error)    // Get the next value in the sequence and update
	Restart()                 // Restarts the sequence
	Current() (uint64, error) // Returns the current value of the Incrementer
	IsStarted() bool          // Returns the state of the Incrementer
}

//=============================================================================

// Init a sequence with uint64 params, ordered similarly to the struct
func (s *Sequence) Init(params ...uint64) {
	if len(params) > 0 {
		s.current = params[0]
	} else {
		s.current = 0
	}

	if len(params) > 1 {
		s.increment = params[1]
	} else {
		s.increment = 1
	}

	if len(params) > 2 {
		s.minvalue = params[2]
	} else {
		s.minvalue = 1
	}

	if len(params) > 3 {
		s.maxvalue = params[3]
	} else {
		s.maxvalue = maxuint64
	}
}

// Next updates the sequence and return the next value
func (s *Sequence) Next() (uint64, error) {
	s.current += s.increment

	// Check for missed minimum condition
	if s.current < s.minvalue {
		return 0, errors.New("Could not reach minimum from current with increment.")
	}

	// Check for reached maximum condition
	if s.current > s.maxvalue {
		return 0, errors.New("Reached maximum bound of sequence.")
	}

	return s.current, nil
}

// Restart the sequence
func (s *Sequence) Restart() {
	s.current = s.minvalue - s.increment
}

// Current returns the current value of the sequence
func (s *Sequence) Current() (uint64, error) {
	if !s.IsStarted() {
		return 0, errors.New("Sequence is unstarted")
	}

	return s.current, nil
}

// IsStarted returns the state of the sequence (started or unstarted)
func (s *Sequence) IsStarted() bool {
	return !(s.current < s.minvalue) && s.current < s.maxvalue
}

// String Representation of the Sequence
func (s *Sequence) String() string {
	d := fmt.Sprintf("incremented by %d between %d and %d", s.increment, s.minvalue, s.maxvalue)
	if !s.IsStarted() {
		return fmt.Sprintf("Unstarted Sequence %s", d)
	}
	return fmt.Sprintf("Sequence at %d, %s", s.current, d)
}
