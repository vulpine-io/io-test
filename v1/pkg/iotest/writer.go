package iotest

// WriteCloser is a test implementation of io.WriteCloser that records function
// calls as well as returning the configured outputs.
type WriteCloser struct {

	// WrittenBytes will contain all bytes passed to the Write method.
	WrittenBytes []byte

	// WriteCalls is a counter of the number of times the Write method was called.
	WriteCalls int

	// WriteErrors defines the errors to return on each call to the Write method.
	//
	// Setting this value to {nil, nil, errors.New("hi")} will cause the Write
	// method to return a nil error on the first two calls, and error("hi") on the
	// third.
	//
	// Calls to the Write method that are outside the length of the WriteErrors
	// slice will return a nil error.
	WriteErrors []error

	// WriteCounts defines the written byte count return value on each call to the
	// Write method.
	//
	// Each entry in this array that is > -1 will be returned in order on calls to
	// the Write method.  A value that is < 0 will cause the Write method to
	// return a written byte count equal to the length of the input byte slice.
	//
	// Calls to the write method that are outside the length of the WriteCounts
	// slice will return the length of the input slice (same as a value < 0).
	WriteCounts []int

	// CloseCalls is a counter of the number of times the Close method was called.
	CloseCalls int

	// CloseErrors defines the errors to return on each call to the close method.
	//
	// Setting this value to {nil, errors.New("hi")} will cause the Close method
	// to return a nil error on the first call, and error("hi") on the next.
	//
	// Calls to the write method that are outside the length of the WriteCounts
	// slice will return a nil error.
	CloseErrors []error
}

func (w *WriteCloser) Write(p []byte) (n int, err error) {
	if len(w.WriteErrors) > w.WriteCalls {
		err = w.WriteErrors[w.WriteCalls]
	}

	if len(w.WriteCounts) > w.WriteCalls && w.WriteCounts[w.WriteCalls] > -1 {
		n = w.WriteCounts[w.WriteCalls]
	} else {
		n = len(p)
	}

	w.WriteCalls++
	w.WrittenBytes = append(w.WrittenBytes, p...)

	return
}

func (w *WriteCloser) Close() (err error) {
	if len(w.CloseErrors) > w.CloseCalls {
		err = w.CloseErrors[w.CloseCalls]
	}

	w.CloseCalls++

	return
}
