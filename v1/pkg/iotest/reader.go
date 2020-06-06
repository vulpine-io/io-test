package iotest

const defaultData = "abcdefghijklmnopqrstuvwxyz"

// ReadCloser is a test implementation of io.ReadCloser that records function
// calls as well as returning the configured outputs.
type ReadCloser struct {

	// ReadableData defines the data that will be read into the passed buffer.
	//
	// The data will be repeated if the length of the buffer or repeated calls
	// exceed the length of this slice.
	//
	// If this value is nil or empty, the returned bytes will be repeating
	// segments of the ascii lowercase alphabet.
	ReadableData []byte

	// ReadCounts defines the "bytes read" count values that will be returned from
	// each call to the read method.
	//
	// This also defines the actual number of bytes that will be copied from the
	// ReadableData slice into the passed buffer.
	//
	// If a defined read count is greater than the size of the input buffer,
	// len(buffer) bytes will be copied, but the configured Read length will be
	// returned.
	//
	// Calls to this method that fall outside the defined returned values will
	// copy the full length of the input buffer and return that length.
	//
	// A value in this array that is < 0 will cause the method to use it's default
	// behavior: copy and return the length of the passed buffer.
	ReadCounts []int

	// ReadErrors defines the error values that will be returned from each call to
	// the Read method.
	//
	// Calls to this method that fall outside the defined return errors will have
	// a returned error of nil.
	ReadErrors []error

	// ReadCalls is a counter of the total number of times the Read method has
	// been called.
	ReadCalls int

	// CloseErrors defines the error values that will be returned from each call
	// to the Close method.
	//
	// Calls to this method that fall outside the defined return errors will have
	// a returned error of nil.
	CloseErrors []error

	// CloseCalls is a counter of the total number of times the Close method has
	// been called.
	CloseCalls int

	readPos int
}

func (r *ReadCloser) Read(p []byte) (n int, err error) {
	// Default the data if none was set
	if len(r.ReadableData) == 0 {
		r.ReadableData = []byte(defaultData)
	}

	if len(r.ReadErrors) > r.ReadCalls {
		err = r.ReadErrors[r.ReadCalls]
	}

	if len(r.ReadCounts) > r.ReadCalls && r.ReadCounts[r.ReadCalls] > -1 {
		n = r.ReadCounts[r.ReadCalls]
	} else {
		n = len(p)
	}

	r.ReadCalls++
	r.readPos = cyclingRead(n, r.readPos, r.ReadableData, p)

	return
}

// Close increments the ReadCloser.CloseCalls counter and optionally returns the
// next error in ReadCloser.CloseErrors.
func (r *ReadCloser) Close() (err error) {
	if len(r.CloseErrors) > r.CloseCalls {
		err = r.CloseErrors[r.CloseCalls]
	}

	r.CloseCalls++

	return
}

func cyclingRead(n, c int, src, dest []byte) int {
	nn := minInt(n, len(dest))

	endPos := nn + c
	// Easy copy, no looping
	if endPos <= len(src) {
		copy(dest, src[c:endPos])
		return endPos
	}

	remaining := nn - len(src[c:])
	copy(dest, src[c:])
	cyclingRead(remaining, 0, src, dest[len(src[c:]):])
	return remaining
}

func minInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}
