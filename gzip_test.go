package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var b64enc = "H4sIAAAAAAAE/wAFAPr/aGVsbG8BAAD//4amEDYFAAAA"

func TestCompression(t *testing.T) {
	Convey("it should compress a string correctly", t, func() {
		So(Compress("hello"), ShouldEqual, b64enc)
	})
}

func TestDecompression(t *testing.T) {
	Convey("it should decompress a string", t, func() {
		So(Decompress(b64enc), ShouldEqual, "hello")
	})
}
