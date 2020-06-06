package iotest_test

import (
	"errors"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vulpine-io/io-test/v1/pkg/iotest"
	"io"
	"testing"
)

func TestReadCloser_Read(t *testing.T) {
	Convey("ReadCloser.Read", t, func() {

		Convey("with no config", func() {
			test := new(iotest.ReadCloser)
			buff := make([]byte, 20)

			a, b := test.Read(buff)

			So(a, ShouldEqual, 20)
			So(b, ShouldBeNil)
			So(test.ReadCalls, ShouldEqual, 1)
			So(string(buff), ShouldEqual, "abcdefghijklmnopqrst")

			a, b = test.Read(buff)

			So(a, ShouldEqual, 20)
			So(b, ShouldBeNil)
			So(test.ReadCalls, ShouldEqual, 2)
			So(string(buff), ShouldEqual, "uvwxyzabcdefghijklmn")

			a, b = test.Read(buff)

			So(a, ShouldEqual, 20)
			So(b, ShouldBeNil)
			So(test.ReadCalls, ShouldEqual, 3)
			So(string(buff), ShouldEqual, "opqrstuvwxyzabcdefgh")
		})

		Convey("with configured data", func() {
			test := &iotest.ReadCloser{
				ReadableData: []byte("i'm a little teapot"),
			}
			buff := make([]byte, 15)

			a, b := test.Read(buff)

			So(a, ShouldEqual, 15)
			So(b, ShouldBeNil)
			So(test.ReadCalls, ShouldEqual, 1)
			So(string(buff), ShouldEqual, "i'm a little te")

			a, b = test.Read(buff)

			So(a, ShouldEqual, 15)
			So(b, ShouldBeNil)
			So(test.ReadCalls, ShouldEqual, 2)
			So(string(buff), ShouldEqual, "apoti'm a littl")

			a, b = test.Read(buff)

			So(a, ShouldEqual, 15)
			So(b, ShouldBeNil)
			So(test.ReadCalls, ShouldEqual, 3)
			So(string(buff), ShouldEqual, "e teapoti'm a l")
		})

		Convey("with configured errors", func() {
			test := &iotest.ReadCloser{
				ReadErrors: []error{nil, io.EOF, errors.New("thot")},
			}
			buff := make([]byte, 8)

			a, b := test.Read(buff)

			So(a, ShouldEqual, 8)
			So(b, ShouldBeNil)
			So(test.ReadCalls, ShouldEqual, 1)
			So(string(buff), ShouldEqual, "abcdefgh")

			a, b = test.Read(buff)

			So(a, ShouldEqual, 8)
			So(b, ShouldEqual, io.EOF)
			So(test.ReadCalls, ShouldEqual, 2)
			So(string(buff), ShouldEqual, "ijklmnop")

			a, b = test.Read(buff)

			So(a, ShouldEqual, 8)
			So(b, ShouldResemble, errors.New("thot"))
			So(test.ReadCalls, ShouldEqual, 3)
			So(string(buff), ShouldEqual, "qrstuvwx")

			a, b = test.Read(buff)

			So(a, ShouldEqual, 8)
			So(b, ShouldBeNil)
			So(test.ReadCalls, ShouldEqual, 4)
			So(string(buff), ShouldEqual, "yzabcdef")
		})

		Convey("with configured counts", func() {
			test := &iotest.ReadCloser{
				ReadCounts: []int{10, -1, 22, 3},
			}
			buff := make([]byte, 8)

			a, b := test.Read(buff)

			So(a, ShouldEqual, 10)
			So(b, ShouldBeNil)
			So(test.ReadCalls, ShouldEqual, 1)
			So(string(buff), ShouldEqual, "abcdefgh")

			a, b = test.Read(buff)

			So(a, ShouldEqual, 8)
			So(b, ShouldBeNil)
			So(test.ReadCalls, ShouldEqual, 2)
			So(string(buff), ShouldEqual, "ijklmnop")

			a, b = test.Read(buff)

			So(a, ShouldEqual, 22)
			So(b, ShouldBeNil)
			So(test.ReadCalls, ShouldEqual, 3)
			So(string(buff), ShouldEqual, "qrstuvwx")

			a, b = test.Read(buff)

			So(a, ShouldEqual, 3)
			So(b, ShouldBeNil)
			So(test.ReadCalls, ShouldEqual, 4)
			So(string(buff), ShouldEqual, "yzatuvwx")

			a, b = test.Read(buff)

			So(a, ShouldEqual, 8)
			So(b, ShouldBeNil)
			So(test.ReadCalls, ShouldEqual, 5)
			So(string(buff), ShouldEqual, "bcdefghi")
		})
	})
}

func TestReadCloser_Close(t *testing.T) {
	Convey("ReadCloser.Close", t, func() {

		Convey("with no config", func() {
			test := new(iotest.ReadCloser)

			So(test.Close(), ShouldBeNil)
			So(test.CloseCalls, ShouldEqual, 1)

			So(test.Close(), ShouldBeNil)
			So(test.CloseCalls, ShouldEqual, 2)
		})

		Convey("with configured errors", func() {
			test := &iotest.ReadCloser{
				CloseErrors: []error{
					nil,
					errors.New("yo"),
				},
			}

			So(test.Close(), ShouldBeNil)
			So(test.CloseCalls, ShouldEqual, 1)

			So(test.Close(), ShouldResemble, errors.New("yo"))
			So(test.CloseCalls, ShouldEqual, 2)

			So(test.Close(), ShouldBeNil)
			So(test.CloseCalls, ShouldEqual, 3)
		})
	})
}

// Zero configuration ReadCloser.
func ExampleReadCloser_Read_noConf() {
	reader := new(iotest.ReadCloser)
	buffer := make([]byte, 10)

	fmt.Println(reader.Read(buffer))
	fmt.Println(string(buffer))
	fmt.Println(reader.ReadCalls)

	// Output: 10 <nil>
	// abcdefghij
	// 1
}

// Short read.
func ExampleReadCloser_Read_shortReads() {
	reader := &iotest.ReadCloser{
		ReadableData: []byte("pastrami"),
		ReadCounts:   []int{2, -1, 5},
	}
	buffer := make([]byte, 10)

	n, _ := reader.Read(buffer)
	fmt.Println(n, buffer)
	n, _ = reader.Read(buffer)
	fmt.Println(n, string(buffer))
	n, _ = reader.Read(buffer)
	fmt.Println(n, string(buffer))

	// Output: 2 [112 97 0 0 0 0 0 0 0 0]
	// 10 stramipast
	// 5 ramipipast
}
