package utils

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)


func TestCheckIsNil(t *testing.T) {

	Convey("check by nil in sequence", t, func() {
		Convey("check by nil where all nil", func() {
			statusNil := IsNill(nil, nil, nil, nil)
			So(statusNil, ShouldBeTrue)
		})

		Convey("check by nil where not all nil", func() {
			statusNil := IsNill(nil, 1, nil, nil)
			So(statusNil, ShouldBeFalse)
		})

	})


	Convey("check sequence by not nil", t, func() {
		Convey("check sequence where all is nil", func() {
			statusNotNil := NotNill(nil, nil, nil, nil)
			So(statusNotNil, ShouldBeFalse)
		})

		Convey("check sequence where has not nil", func() {
			statusNotNil := NotNill(nil, 1, nil, nil)
			So(statusNotNil, ShouldBeFalse)
		})

		Convey("check sequence where all exists", func() {
			statusNotNilTrue := NotNill(1, 2, 3, 4)
			So(statusNotNilTrue, ShouldBeTrue)
		})


	})
}
