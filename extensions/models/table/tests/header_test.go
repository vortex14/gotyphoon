package tests

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vortex14/gotyphoon/extensions/models/table"
	"reflect"
	"testing"
)

func TestCreateHeader(t *testing.T) {
	newTable := table.Table{}
	Convey("set headers", t, func() {
		headers := table.H{"h 1", "h 2"}
		newTable.SetHeaders(headers)
		Convey("check header fields", func() {
			tH := newTable.GetHeaders()
			assertH := append(table.H{"№"}, headers...)
			So(reflect.DeepEqual(tH, assertH), ShouldEqual, true)
		})
	})
}
