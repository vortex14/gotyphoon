package tests

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vortex14/gotyphoon/elements/models/table"
	Errors "github.com/vortex14/gotyphoon/errors"
	"testing"
)

func TestCreateTable(t *testing.T) {

	Convey("Create a new table", t, func() {

		newTable := table.Table{}

		Convey("checking state row", func() {
			So(newTable.GetCurrentRow(), ShouldEqual, 0)
		})


		Convey("add a new row with exception TableHeadersNotFound", func() {
			So(newTable.Append("",table.H{}), ShouldEqual, Errors.TableHeadersNotFound)
		})




	})
}
