Go-Sequence
===========

**Implements an AutoIncrement counter class similar to PostgreSQL's sequence.**

Honestly, this package is really more of a first attempt at creating a Go
library for myself - and I've found that I've needed to create monotonically
increasing counters for a variety of reasons in my packages. I created this
system to mirror PostgreSQL's sequence data structure to easily add
incrementing data structures to your projects.

Usage
-----

Using sequences is pretty straight forward, here is my common use case:

    import (
        "github.com/bbengfort/benfs/dss/sequence"
    )

    // Initialize new counter starting at 1 and incrementing by 1.
    counter := sequence.New()
    first := counter.Next()
    second := counter.Next()

Sequences maintain their state and you can check if they've been started:

    counter.IsStarted()

You can also reseet sequences to return them to a starting value:

    counter.Restart()

Sequences are bounded by a minimum and maximum value and are incremented
by a step that you can choose (the default is 1). If a counter goes above
a maximum value you get an error. Note that if you simply use a `uint64`
this is not the behavior you will see!

To initialize a counter with a different minimum, maximum and step - say
to count by even numbers starting at 10 and going to 100:

    counter := new(sequence.Sequence)
    counter.Init(9, 10, 100, 2)

Note that the current value should be initialized to one less than the
minimum to start with; this will be taken care of in the future!

Testing
-------

To execute tests for this library:

    go test github.com/bbengfort/sequence

Hopefully everything passes!

### TODO:

1. Fix minimum and current initialization to allow for (min, max] API
2. Eliminate the need for a specification, make it more like `xrange`
3. Implement an iterable for the `range` function
