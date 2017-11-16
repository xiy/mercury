package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var b64enc = "H4sIAAAAAAAE/wAPAPD/PGgxPkhlbGxvITwvaDE+AQAA//9WO9LMDwAAAA=="
var plainString = `<h1>Hello!</h1>`

func TestCompression(t *testing.T) {
	Convey("it should compress a string correctly", t, func() {
		So(Compress(plainString), ShouldEqual, b64enc)
	})
}

func TestDecompression(t *testing.T) {
	Convey("it should decompress a string", t, func() {
		So(Decompress(b64enc), ShouldEqual, plainString)
	})
}
