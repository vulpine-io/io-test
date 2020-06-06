package iotest_test

import (
	"errors"
	"fmt"
	"github.com/vulpine-io/io-test/v1/pkg/iotest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWriteCloser_Write(t *testing.T) {
	Convey("WriteCloser.Write", t, func() {

		Convey("with no config", func() {
			test := new(iotest.WriteCloser)
			input := []byte{1, 2, 3}

			a, b := test.Write(input)

			So(a, ShouldEqual, len(input))
			So(b, ShouldBeNil)
			So(test.WriteCalls, ShouldEqual, 1)
			So(test.WrittenBytes, ShouldResemble, input)
		})

		Convey("with defined errors", func() {
			test := &iotest.WriteCloser{
				WriteErrors: []error{
					errors.New("hello"),
					errors.New("world"),
				},
			}
			input := []byte{1, 2, 3}

			a, b := test.Write(input)

			So(a, ShouldEqual, len(input))
			So(b, ShouldResemble, errors.New("hello"))
			So(test.WriteCalls, ShouldEqual, 1)
			So(test.WrittenBytes, ShouldResemble, input)

			a, b = test.Write(input)

			So(a, ShouldEqual, len(input))
			So(b, ShouldResemble, errors.New("world"))
			So(test.WriteCalls, ShouldEqual, 2)
			So(test.WrittenBytes, ShouldResemble, append(input, input...))

			a, b = test.Write(input)

			So(a, ShouldEqual, len(input))
			So(b, ShouldBeNil)
			So(test.WriteCalls, ShouldEqual, 3)
			So(test.WrittenBytes, ShouldResemble, append(append(input, input...), input...))
		})

		Convey("with defined write counts", func() {
			test := &iotest.WriteCloser{
				WriteCounts: []int{1, -1, 2, -1, 3},
			}
			var input []byte

			a, b := test.Write(input)

			So(a, ShouldEqual, 1)
			So(b, ShouldBeNil)
			So(test.WriteCalls, ShouldEqual, 1)
			So(test.WrittenBytes, ShouldResemble, input)

			a, b = test.Write(input)

			So(a, ShouldEqual, len(input))
			So(b, ShouldBeNil)
			So(test.WriteCalls, ShouldEqual, 2)
			So(test.WrittenBytes, ShouldResemble, input)

			a, b = test.Write(input)

			So(a, ShouldEqual, 2)
			So(test.WriteCalls, ShouldEqual, 3)

			a, b = test.Write(input)

			So(a, ShouldEqual, len(input))
			So(test.WriteCalls, ShouldEqual, 4)

			a, b = test.Write(input)

			So(a, ShouldEqual, 3)
			So(test.WriteCalls, ShouldEqual, 5)

			a, b = test.Write(input)

			So(a, ShouldEqual, len(input))
			So(test.WriteCalls, ShouldEqual, 6)
		})

	})
}

func TestWriteCloser_Close(t *testing.T) {
	Convey("WriteCloser.Close", t, func() {

		Convey("with no config", func() {
			test := new(iotest.WriteCloser)

			So(test.Close(), ShouldBeNil)
			So(test.CloseCalls, ShouldEqual, 1)
			So(test.Close(), ShouldBeNil)
			So(test.CloseCalls, ShouldEqual, 2)
			So(test.Close(), ShouldBeNil)
			So(test.CloseCalls, ShouldEqual, 3)
		})

		Convey("with defined errors", func() {
			test := &iotest.WriteCloser{
				CloseErrors: []error{
					errors.New("goodbye"),
					errors.New("cruel"),
					errors.New("world"),
				},
			}

			So(test.Close(), ShouldResemble, errors.New("goodbye"))
			So(test.CloseCalls, ShouldEqual, 1)
			So(test.Close(), ShouldResemble, errors.New("cruel"))
			So(test.CloseCalls, ShouldEqual, 2)
			So(test.Close(), ShouldResemble, errors.New("world"))
			So(test.CloseCalls, ShouldEqual, 3)
			So(test.Close(), ShouldBeNil)
			So(test.CloseCalls, ShouldEqual, 4)
		})

	})
}

// Zero config WriteCloser usage.
func ExampleWriteCloser_Write_noConfig() {
	writer := new(iotest.WriteCloser)
	data := []byte("hey there!")

	fmt.Println(writer.Write(data))
	fmt.Println(string(writer.WrittenBytes))
	fmt.Println(writer.WriteCalls)

	// Output: 10 <nil>
	// hey there!
	// 1
}

// WriteCloser with configured errors.
func ExampleWriteCloser_Write_errConfig() {
	writer := &iotest.WriteCloser{
		WriteErrors: []error{nil, nil, errors.New("sup, bro")},
	}

	fmt.Println(writer.Write([]byte("nah")))
	fmt.Println(writer.Write([]byte("cya")))
	fmt.Println(writer.Write([]byte("l8r")))
	fmt.Println(writer.Write([]byte("brah")))

	// Output: 3 <nil>
	// 3 <nil>
	// 3 sup, bro
	// 4 <nil>
}

// WriteCloser with configured write counts.
func ExampleWriteCloser_Write_countConfig() {
	writer := &iotest.WriteCloser{
		WriteCounts: []int{0, -1, 22, -1},
	}

	fmt.Println(writer.Write([]byte("hey")))
	fmt.Println(writer.Write([]byte("u")))
	fmt.Println(writer.Write([]byte("up")))
	fmt.Println(writer.Write([]byte("rn?")))

	// Output: 0 <nil>
	// 1 <nil>
	// 22 <nil>
	// 3 <nil>
}
