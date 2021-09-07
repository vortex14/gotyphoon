package tests

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vortex14/gotyphoon/extensions/models/table"
	"testing"
)

func TestRowTable(t *testing.T) {

	Convey("add a new row with data", t, func() {
		newTable := table.Table{}
		Convey("set header", func() {
			newTable.SetHeaders(table.H{"id", "name"})
			Convey("add new row", func() {
				errW := newTable.Append("1", table.H{"1"})
				Convey("check exception after append", func() {
					So(errW, ShouldEqual, nil)
				})

				Convey("check count row after append", func() {
					So(newTable.GetCountRow(), ShouldEqual, 1)
				})

			})

		})

	})

}
